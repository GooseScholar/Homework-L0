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

func NewHttpServer(ctx context.Context, repo *database.DB, cache *cache.Cache) *httpServer {
	ServerMux := http.NewServeMux()

	ServerMux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		id := r.FormValue("id")
		if len(id) < 1 {
			io.WriteString(w, `<!DOCTYPE html>
			<form action="/orders" method="get">
			Ввевидети ваш номер заказа:
			<br>
			<input name="id" type="text">
			<br> 
			<button type="submit">Показать</button>
		</form>`)
			return
		}

		if cached, found := cache.GetOrder(id); found == false {
			order, err := repo.GetOrder(ctx, id)
			if err != nil {
				io.WriteString(w, fmt.Sprintf(`<!DOCTYPE html>
				<form action="/orders" method="get">
				Ввевидети ваш номер заказа:
				<br>
				<input name="id" type="text">
				<br> 
				<button type="submit">Показать</button>
			</form>
			<div> "Order <%s> not found"
			</div>`, id))
				return
			}

			jsn, err := json.Marshal(order)

			if err != nil {
				io.WriteString(w, fmt.Sprintf("Why? %v", err))
				return
			}

			cache.Data[order.Order_uid] = string(jsn)

			pref := `<!DOCTYPE html>
			<form action="/orders" method="get">
			Ввевидети ваш номер заказа:
			<br>
			<input name="id" type="text">
			<br> 
			<button type="submit">Показать</button>
		</form>
		<div>`
			suff := `</div>`
			io.WriteString(w, fmt.Sprintf("%s%v%s", pref, string(jsn), suff))

		} else {
			pref := `<!DOCTYPE html>
			<form action="/orders" method="get">
			Ввевидети ваш номер заказа:
			<br>
			<input name="id" type="text">
			<br> 
			<button type="submit">Показать</button>
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
