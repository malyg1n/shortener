package rest

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	md "github.com/go-chi/chi/v5/middleware"
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
	fmt.Println("server starting")
	cfg := config.GetConfig()
	storage, err := pgsql.NewLinksStoragePG(ctx)
	if err != nil {
		return err
	}
	fmt.Println("storage started")

	linker, err := v1.NewDefaultLinker(storage)
	if err != nil {
		return err
	}
	fmt.Println("linker started")

	handler, err := handlers.NewHandlerManager(linker)
	if err != nil {
		return err
	}
	fmt.Println("handler started")

	router := chi.NewRouter()
	router.Get("/{linkId}", handler.GetLink)
	router.Post("/", handler.SetLink)
	router.Post("/api/shorten", handler.APISetLink)
	router.Get("/user/urls", handler.GetLinksByUser)
	router.Get("/ping", handler.PingDB)
	router.Post("/api/shorten/batch", handler.APISetBatchLinks)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: middleware.Compress(middleware.Decompress(middleware.Cookies(md.Logger(router)))),
	}

	go func() {
		err := srv.ListenAndServe()
		fmt.Println(err.Error())
		fmt.Println("server stopped")
	}()

	<-ctx.Done()
	fmt.Println("ctx done")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		fmt.Println("timeout context cancel")
		cancel()
	}()

	e := storage.Close()
	fmt.Println("storage closed", e)

	return srv.Shutdown(ctxShutDown)
}
