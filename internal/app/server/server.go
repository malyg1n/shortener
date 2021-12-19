package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/internal/app/handlers"
	"github.com/malyg1n/shortener/internal/app/services"
	"github.com/malyg1n/shortener/internal/app/storage"
	"log"
	"net/http"
)

// Run ...
func Run() {
	h := handlers.NewBaseHandler(services.NewDefaultLinksService(storage.NewLinksStorageMap()))
	r := chi.NewRouter()
	r.Get("/{linkID}", h.GetLink)
	r.Post("/", h.SetLink)
	log.Fatalln(http.ListenAndServe(":8080", r))
}
