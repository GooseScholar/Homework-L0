package main

import (
	"context"
	"homework-l0/internal/app"
	"homework-l0/internal/database"
	"homework-l0/internal/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	repo, err := database.NewDB(ctx)
	if err != nil {
		log.Fatalf("connect DB: <%v>", err)
	}
	cache := app.NewCache()

	//service := service.New(db)

	server.NewHttpServer(ctx, repo, cache)

	app.Subscriber(ctx, repo, cache)

	<-ctx.Done()

}
