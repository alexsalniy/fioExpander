CREATE TABLE IF NOT EXISTS fio (
  id VARCHAR NOT NULL UNIQUE PRIMARY KEY,
  name VARCHAR (20) NOT NULL, 
  surname VARCHAR (20) NOT NULL,
  patronymic VARCHAR (20),
  age INT,
  gender VARCHAR(10),
  gender_probability FLOAT,
  nation VARCHAR
);