package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/handlers"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/filesystem"
	"net/http"
	"time"
)

// RunServer init routes adn listen
func RunServer(ctx context.Context) error {
	cfg := config.GetConfig()
	storage, err := filesystem.NewLinksStorageFile()
	if err != nil {
		return err
	}

	linker, err := v1.NewDefaultLinker(storage)
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
		Addr:    cfg.Addr,
		Handler: middleware.Compress(middleware.Decompress(router)),
	}

	go func() {
		_ = srv.ListenAndServe()
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	storage.Close()

	return srv.Shutdown(ctxShutDown)
}
