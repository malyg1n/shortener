package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/internal/app/api/rest/handlers"
	v1 "github.com/malyg1n/shortener/internal/app/services/linker/v1"
	"github.com/malyg1n/shortener/internal/app/storage/inmemory"
	"net/http"
)

// RunServer init routes adn listen
func RunServer() error {
	handler := handlers.NewLinksHandler(v1.NewDefaultLinker(inmemory.NewLinksStorageMap()))

	router := chi.NewRouter()
	router.Get("/{linkId}", handler.GetLink)
	router.Post("/", handler.SetLink)

	return http.ListenAndServe(":8080", router)
}
