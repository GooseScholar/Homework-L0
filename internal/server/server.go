package server

import (
	"context"
	"encoding/json"
	"fmt"
	"homework-l0/internal/app"
	"homework-l0/internal/database"
	"io"
	"log"
	"net/http"
	"os"
)

type httpServer struct {
	repo  Repository
	cache app.Cache
}

//_ = products.NewHttpServer(ctx)

func NewHttpServer(ctx context.Context, repo *database.DB, cache *app.Cache) *httpServer {
	ServerMux := http.NewServeMux()

	ServerMux.HandleFunc("/postgres", func(w http.ResponseWriter, r *http.Request) {
		idStrs, ok := r.URL.Query()["id"]
		if !ok {
			io.WriteString(w, "Missing query parameter <id>")
		}

		id := idStrs[0]
		if len(id) < 1 {
			io.WriteString(w, "Incorrect query parameter <id>")
		}

		if cached, found := cache.GetOrder(id); found == false {
			order, err := repo.GetOrder(ctx, id)
			if err != nil {
				io.WriteString(w, fmt.Sprintf("Order <%s> not found", id))
				return
			}

			jsn, err := json.Marshal(order)

			if err != nil {
				io.WriteString(w, fmt.Sprintf("Why? %v", err))
				return
			}

			io.WriteString(w, fmt.Sprintf("Not found in cache: %v", jsn))

		} else {
			io.WriteString(w, fmt.Sprintf("Found in cache: %v", cached))
		}
		go func() {
			err := http.ListenAndServe(":"+os.Getenv("portServer"), ServerMux)
			if err != nil {
				log.Println("Failed to listen" + os.Getenv("portServer"))
			}
		}()

	})

	return &httpServer{
		repo:  repo,
		cache: *cache,
	}
}
