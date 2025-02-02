package server

import (
	"apollo-image-processor/internal/handler"
	"net/http"

	"github.com/go-chi/chi"
)

func addRoutes(r *chi.Mux, ih *handler.ImageHandler) {

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ping!"))
	})

	r.Route("/images", func(r chi.Router) {

		r.Post("/upload", ih.UploadImages)

	})

}
