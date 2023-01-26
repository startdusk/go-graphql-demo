package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database database
}

func New(filenames ...string) (*Config, error) {
	if err := godotenv.Load(filenames...); err != nil {
		return nil, err
	}

	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
	}, nil
}

type database struct {
	URL string
}
