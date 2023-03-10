package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/startdusk/twitter/config"
	"github.com/startdusk/twitter/data"
	"github.com/startdusk/twitter/data/postgres"
	"github.com/startdusk/twitter/domain"
	"github.com/startdusk/twitter/jwt"
)

var (
	authService data.AuthService

	tweetService data.TweetService
)

// StartDB starts a database instance.
func StartDB(dbname string) (*Container, error) {
	image := "postgres:14-alpine"
	port := "5432"
	args := []string{
		"-e", fmt.Sprintf("POSTGRES_DB=%s", dbname),
		"-e", "POSTGRES_USER=postgres",
		"-e", "POSTGRES_PASSWORD=password",
	}

	return StartContainer(image, port, args...)
}

// StopDB stops a running database instance.
func StopDB(c *Container) {
	StopContainer(c.ID)
}

func TestMain(m *testing.M) {
	conf, err := config.New("../.env.test")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	dbname := uuid.NewString()
	container, err := StartDB(dbname)
	if err != nil {
		panic(err)
	}
	defer StopDB(container)

	databaseURL := fmt.Sprintf(
		"postgres://postgres:password@%s/%s?sslmode=disable",
		container.Host,
		dbname,
	)

	db, err := postgres.New(ctx, databaseURL, 10)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Migrate(); err != nil {
		panic(err)
	}

	userRepo := &postgres.UserRepo{
		DB: db,
	}
	tweetRepo := &postgres.TweetRepo{
		DB: db,
	}
	refreshTokenRepo := &postgres.RefreshTokenRepo{
		DB: db,
	}

	authTokenService := jwt.NewTokenService(&conf.JWT)
	authService = domain.NewAuthService(userRepo, authTokenService, refreshTokenRepo)
	tweetService = domain.NewTweetService(tweetRepo)
	m.Run()
}
