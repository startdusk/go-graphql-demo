package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/startdusk/twitter/cmd/twitter/graph"
	"github.com/startdusk/twitter/config"
	"github.com/startdusk/twitter/data/postgres"
	"github.com/startdusk/twitter/domain"
	"github.com/startdusk/twitter/jwt"
)

func main() {
	ctx := context.Background()
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n", conf)
	db, err := postgres.New(ctx, conf.Database.URL, 10)
	if err != nil {
		panic(err)
	}
	if err := db.Migrate(); err != nil {
		panic(err)
	}

	tokenService := jwt.NewTokenService(&conf.JWT)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(requestid.New())
	router.Use(gin.Recovery())
	router.Use(authMiddleware(tokenService))
	router.GET("/", func(ctx *gin.Context) {
		playground.Handler("Twitter clone", "/query").ServeHTTP(ctx.Writer, ctx.Request)
	})
	router.POST("/query", func(ctx *gin.Context) {
		handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{
				AuthService: domain.NewAuthService(&postgres.UserRepo{
					DB: db,
				}, jwt.NewTokenService(&conf.JWT)),
			},
			Directives: graph.DirectiveRoot{},
			Complexity: graph.ComplexityRoot{},
		})).ServeHTTP(ctx.Writer, ctx.Request)
	})

	log.Println("run it success")
	log.Fatal(http.ListenAndServe(":8888", router))
}
