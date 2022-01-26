package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/handlers"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/pgsql"
	"net/http"
	"time"
)

// RunServer init routes adn listen
func RunServer(ctx context.Context) error {
	cfg := config.GetConfig()
	storage, err := pgsql.NewLinksStoragePG(ctx)
	if err != nil {
		return err
	}

	linker, err := v1.NewDefaultLinker(storage)
	if err != nil {
		return err
	}

	handler, err := handlers.NewHandlerManager(linker)
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Get("/{linkId}", handler.GetLink)
	router.Post("/", handler.SetLink)
	router.Post("/api/shorten", handler.APISetLink)
	router.Get("/user/urls", handler.GetLinksByUser)
	router.Get("/ping", handler.PingDB)
	router.Post("/api/shorten/batch", handler.APISetBatchLinks)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: middleware.Compress(middleware.Decompress(middleware.Cookies(router))),
	}

	go func() {
		_ = srv.ListenAndServe()
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	_ = storage.Close()

	return srv.Shutdown(ctxShutDown)
}
