package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database database
	JWT      JWT
}

func New(filenames ...string) (*Config, error) {
	if err := godotenv.Load(filenames...); err != nil {
		return nil, err
	}

	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
		JWT: JWT{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("DOMAIN"),
		},
	}, nil
}

type database struct {
	URL string
}

type JWT struct {
	Secret string
	Issuer string
}
