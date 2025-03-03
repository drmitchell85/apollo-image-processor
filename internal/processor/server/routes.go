package processorserver

import (
	"net/http"

	"github.com/go-chi/chi"
)

func addRoutes(r *chi.Mux) {

	r.Use(prometheusMiddleware)
	initPrometheus(r)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ping!"))
	})

}
