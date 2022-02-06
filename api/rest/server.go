package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/handlers"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/services/linker"
	"net/http"
)

// APIServer struct.
type APIServer struct {
	handlerManager *handlers.HandlerManager
	server         *http.Server
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

	router := chi.NewRouter()
	router.Get("/{linkId}", srv.handlerManager.GetLink)
	router.Post("/", srv.handlerManager.SetLink)
	router.Post("/api/shorten", srv.handlerManager.APISetLink)
	router.Get("/user/urls", srv.handlerManager.GetLinksByUser)
	router.Get("/ping", srv.handlerManager.PingDB)
	router.Post("/api/shorten/batch", srv.handlerManager.APISetBatchLinks)

	srv.server.Handler = middleware.Compress(middleware.Decompress(middleware.Cookies(router)))

	go func() {
		_ = srv.server.ListenAndServe()
	}()
}
