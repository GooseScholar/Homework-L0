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
	log.Printf("стоим ньюкешем %v", err)
	cache, err := repo.GetInitialCache(ctx)
	if err != nil {
		log.Fatalf("cache: <%v>", err)
	}
	log.Printf("стоим перед %v", err)
	log.Printf("стоим после %v", err)
	if err != nil {
		log.Printf("recovery cache failed: <%v>", err)
	}

	//service := service.New(db)

	server.NewHttpServer(ctx, repo, cache)

	app.Subscriber(ctx, repo, cache)

	<-ctx.Done()
	// Wait for Ctrl+C
	/*	doneCh1 := make(chan bool)
		go func() {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)
			<-sigCh
			doneCh1 <- true
		}()
		<-doneCh1
	*/
}
