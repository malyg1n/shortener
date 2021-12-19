package server

import (
	"github.com/malyg1n/shortener/internal/app/handlers"
	"github.com/malyg1n/shortener/internal/app/services"
	"github.com/malyg1n/shortener/internal/app/storage"
	"log"
	"net/http"
)

// Run ...
func Run() {
	hs := handlers.NewHandlers(services.NewDefaultLinksService(storage.NewLinksStorageMap()))
	http.HandleFunc("/", hs.Resolve)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
