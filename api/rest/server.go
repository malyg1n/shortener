package rest

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/handlers"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/services/linker"
	"github.com/malyg1n/shortener/storage"
	"net/http"
	"time"
)

// APIServer struct.
type APIServer struct {
	handlerManager *handlers.HandlerManager
	server         *http.Server
	storage        storage.LinksStorage
}

// NewAPIServer creates new instance
func NewAPIServer(service linker.Linker, addr string) (*APIServer, error) {
	handler, err := handlers.NewHandlerManager(service)
	if err != nil {
		return nil, err
	}

	server := &APIServer{
		handlerManager: handler,
		server:         &http.Server{Addr: addr},
	}

	return server, nil
}

// Run server.
func (srv *APIServer) Run(ctx context.Context) {

	router := chi.NewRouter().With(middleware.Compress, middleware.Decompress, middleware.Cookies)
	router.Get("/{linkId}", srv.handlerManager.GetLink)
	router.Post("/", srv.handlerManager.SetLink)
	router.Post("/api/shorten", srv.handlerManager.APISetLink)
	router.Get("/user/urls", srv.handlerManager.GetLinksByUser)
	router.Get("/ping", srv.handlerManager.PingDB)
	router.Post("/api/shorten/batch", srv.handlerManager.APISetBatchLinks)
	router.Delete("/api/user/urls", srv.handlerManager.DeleteUserLinks)

	srv.server.Handler = router
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		fmt.Println("shutdown")
		srv.server.Shutdown(ctxShutdown)
	}()

	go func() {
		_ = srv.server.ListenAndServe()
	}()
}
