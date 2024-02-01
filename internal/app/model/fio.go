package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type Nation struct {
	CID  string  `json:"country_id"`
	Prob float32 `json:"probability"`
}

type fioResponse struct {
	Count   int      `json:"count"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Gender  string   `json:"gender"`
	Prob    float32  `json:"probability"`
	Country []Nation `json:"country"`
}

type FIO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type ExtendedFIO struct {
	ID          uuid.UUID
	Name        string `json:"name"`
	Surname     string
	Patronymic  string
	Age         int      `json:"age"`
	Gender      string   `json:"gender"`
	Probability float32  `json:"probability"`
	Country     []Nation `json:"country"`
	Nation      string
}

type FindFIO struct {
	ExtendedFIO
	Page int `json:"page"`
}

func (e *ExtendedFIO) Validator() bool {

	if len(e.Name) == 0 {
		fmt.Println("struct does not have the field", e.Name)
		return false
	}

	if len(e.Name) == 0 {
		fmt.Println("struct does not have the field", e.Surname)
		return false
	}
	return true
}

func (r *ExtendedFIO) GetExtension(extUrl string) (*ExtendedFIO, error) {
	url := extUrl + r.Name

	res, err := http.Get(fmt.Sprintf(url))
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}

	rjson, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}

	res.Body.Close()

	err = json.Unmarshal(rjson, &r)
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}

	return r, nil
}
