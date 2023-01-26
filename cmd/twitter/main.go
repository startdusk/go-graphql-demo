package main

import (
	"context"
	"log"

	"github.com/startdusk/twitter/config"
	"github.com/startdusk/twitter/data/postgres"
)

func main() {
	ctx := context.Background()
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n", conf)
	db, err := postgres.New(ctx, conf.Database.URL)
	if err != nil {
		panic(err)
	}
	if err := db.Migrate(); err != nil {
		panic(err)
	}
	log.Println("run it success")
}
