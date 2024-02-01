package psql

import (
	"database/sql"
	"fio-expander/internal/app/model"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dbInfo string) (*Storage, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dbInfo, err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dbInfo, err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../internal/storage/migrations",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("cannot create new migrate instance: %w", err)
	}

	if err = m.Up(); err != nil {
		fmt.Println("failed to run up migrate:", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Create(e *model.ExtendedFIO) error {
	if e.Validator() {

		fmt.Println("fio before extansion", e)

		e.GetExtension("https://api.agify.io/?name=")
		e.GetExtension("https://api.genderize.io/?name=")
		e.GetExtension("https://api.nationalize.io/?name=")

		fmt.Println("fio after extansion", e)

		e.ID = uuid.New()

		for i := 0; i < len(e.Country); i++ {
			e.Nation = e.Nation + e.Country[i].CID + " " + fmt.Sprintf("%f", e.Country[i].Prob) + " "
		}

		if err := s.db.QueryRow(`
			INSERT INTO fio (id, name, surname, patronymic, age, gender, gender_probability, nation)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`,
			e.ID,
			e.Name,
			e.Surname,
			e.Patronymic,
			e.Age,
			e.Gender,
			e.Probability,
			e.Nation,
		).Scan(&e.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Update(e *model.ExtendedFIO) error {
	if e.Validator() {

		e.GetExtension("https://api.agify.io/?name=")
		e.GetExtension("https://api.genderize.io/?name=")
		e.GetExtension("https://api.nationalize.io/?name=")

		for i := 0; i < len(e.Country); i++ {
			e.Nation = e.Nation + e.Country[i].CID + " " + fmt.Sprintf("%f", e.Country[i].Prob) + " "
		}

		query := fmt.Sprintf(`
		UPDATE fio SET (name, surname, patronymic, age, gender, gender_probability, nation) =
				(%s, %s, %s, %v, %s, %v, %s)
				WHERE id = %v
				RETURNING id`, e.Name, e.Surname, e.Patronymic, e.Age, e.Gender, e.Probability, e.Nation, e.ID)
		fmt.Println(query)

		if err := s.db.QueryRow(`
			UPDATE fio SET (name, surname, patronymic, age, gender, gender_probability, nation) =
			($2, $3, $4, $5, $6, $7, $8)
		 	WHERE id = $1
			RETURNING id`,

			e.ID,
			e.Name,
			e.Surname,
			e.Patronymic,
			e.Age,
			e.Gender,
			e.Probability,
			e.Nation,
		).Scan(&e.ID); err != nil {
			return err
		}

		// rows, err := s.db.Query(
		// 	fmt.Sprintf(`
		// 		UPDATE fio SET (name, surname, patronymic, age, gender, gender_probability, nation) =
		// 		(%s, %s, %s, %v, %s, %v, %s)
		// 		WHERE id = '%v'
		// 		RETURNING id`,
		// 		e.Name,
		// 		e.Surname,
		// 		e.Patronymic,
		// 		e.Age,
		// 		e.Gender,
		// 		e.Probability,
		// 		e.Nation,
		// 		e.ID,
		// 	),
		// )
		// if err != nil {
		// 	return err
		// }
		// if err := rows.Scan(&e.ID); err != nil {
		// 	return err
		// }
	}

	return nil
}

func (s *Storage) FindBy(e *model.FindFIO) (error, *[]model.ExtendedFIO) {

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	selectQuery := queryBuilder.Select("*").From("fio")

	var zeroID uuid.UUID

	if e.ID != zeroID {
		selectQuery = selectQuery.Where(sq.Eq{
			"id": e.ID,
		})
	}

	if e.Name != "" {
		selectQuery = selectQuery.Where(sq.Eq{
			"name": e.Name,
		})
	}

	if e.Surname != "" {
		selectQuery = selectQuery.Where(sq.Eq{
			"surname": e.Surname,
		})
	}

	if e.Patronymic != "" {
		selectQuery = selectQuery.Where(sq.Eq{
			"patronymic": e.Patronymic,
		})
	}

	if e.Age != 0 {
		selectQuery = selectQuery.Where(sq.Eq{
			"age": e.Age,
		})
	}

	if e.Gender != "" {
		selectQuery = selectQuery.Where(sq.Eq{
			"gender": e.Gender,
		})
	}

	var limit uint64 = 5
	var offset uint64 = limit * uint64(e.Page-1)
	selectQuery = selectQuery.Limit(limit).Offset(offset)

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return err, nil
	}

	fmt.Println("query for fin by", sql)
	rows, err := s.db.Query(sql, args...)
	fmt.Println(rows)
	defer rows.Close()
	res := make([]model.ExtendedFIO, 0)
	for rows.Next() {
		var ans model.ExtendedFIO
		err = rows.Scan(&ans.ID, &ans.Name, &ans.Surname, &ans.Patronymic, &ans.Age, &ans.Gender, &ans.Probability, &ans.Nation)
		if err != nil {
			return err, nil
		}
		res = append(res, ans)
	}
	return nil, &res
}

func (s *Storage) Delete(e *model.ExtendedFIO) error {

	if _, err := s.db.Query(`
	DELETE FROM fio 
	WHERE id = $1`,
		e.ID,
	); err != nil {
		return err
	}

	return nil
}
