package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/handlers"
	"github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/inmemory"
	"net/http"
	"time"
)

// RunServer init routes adn listen
func RunServer(ctx context.Context) error {
	linker, err := v1.NewDefaultLinker(inmemory.NewLinksStorageMap())
	if err != nil {
		return err
	}

	handler, err := handlers.NewLinksHandler(linker)
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Get("/{linkId}", handler.GetLink)
	router.Post("/", handler.SetLink)
	router.Post("/api/shorten", handler.APISetLink)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		srv.ListenAndServe()
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	return srv.Shutdown(ctxShutDown)
}
