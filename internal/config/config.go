package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	Server   string `env:"SERVER" envDefault:"localhost:8080"`
	Dbhost   string `env:"DBHOST"`
	Dbport   string `env:"DBPORT"`
	Dbuser   string `env:"DBUSER"`
	Dbpass   string `env:"DBPASS"`
	Dbname   string `env:"DBNAME"`
	Dbinfo   string
	Dbsource string
}

func MustLoad() *Config {
	configPath := "../../.env"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	cfg.Dbinfo = fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		cfg.Dbhost, cfg.Dbport, cfg.Dbuser, cfg.Dbpass, cfg.Dbname)
	cfg.Dbsource = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Dbuser, cfg.Dbpass, cfg.Dbhost, cfg.Dbport, cfg.Dbname)
	return &cfg
}
