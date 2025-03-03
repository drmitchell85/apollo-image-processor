package apiserver

import (
	"apollo-image-processor/internal/api/handler"
	"net/http"

	"github.com/go-chi/chi"
)

func addRoutes(r *chi.Mux, ih *handler.ImageHandler) {

	r.Use(prometheusMiddleware)
	initPrometheus(r)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ping!"))
	})

	r.Route("/images", func(r chi.Router) {

		r.Post("/upload", ih.UploadImages)

	})

}
