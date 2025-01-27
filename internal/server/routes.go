package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

func addRoutes(r *chi.Mux) {

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ping!"))
	})

}
