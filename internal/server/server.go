package server

import (
	"context"
	"encoding/json"
	"fmt"
	"homework-l0/internal/cache"
	"homework-l0/internal/database"
	"io"
	"log"
	"net/http"
	"os"
)

type httpServer struct {
	repo  Repository
	cache cache.Cache
}

//_ = products.NewHttpServer(ctx)

func NewHttpServer(ctx context.Context, repo *database.DB, cache *cache.Cache) *httpServer {
	ServerMux := http.NewServeMux()

	ServerMux.HandleFunc("/postgres", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Зашли")
		r.ParseForm()

		id := r.FormValue("id")
		if len(id) < 1 {
			io.WriteString(w, "Incorrect query parameter <id>")
			return
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

			cache.Data[order.Order_uid] = string(jsn)

			io.WriteString(w, fmt.Sprintf("Not found in cache: %v", jsn))

		} else {
			pref := `<form action="/postgres" method="get">
			<input name="id" type="text">
			<button>Показать</button>
		</form>
		<div>`
			suff := `</div>`
			io.WriteString(w, fmt.Sprintf("%s%v%s", pref, cached, suff))
		}

	})

	go func() {
		err := http.ListenAndServe(":"+os.Getenv("portServer"), ServerMux)
		if err != nil {
			log.Println("Failed to listen" + os.Getenv("portServer"))
		}
	}()

	return &httpServer{
		repo:  repo,
		cache: *cache,
	}
}
