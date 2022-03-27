package rest

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/handlers"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/services/linker"
	"github.com/malyg1n/shortener/storage"
	"log"
	"net/http"
	"time"
)

// APIServer struct.
type APIServer struct {
	handlerManager *handlers.HandlerManager
	server         *http.Server
	storage        storage.LinksStorage
	useSSL         bool
	SSLCert        string
	SSLKey         string
}

// NewAPIServer creates new instance
func NewAPIServer(service linker.Linker, addr string, useSSL bool, sslCert, sslKey string) (*APIServer, error) {
	handler, err := handlers.NewHandlerManager(service)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Addr: addr}

	server := &APIServer{
		handlerManager: handler,
		server:         srv,
		useSSL:         useSSL,
	}

	if useSSL {
		server.server.Addr = ""
		server.SSLCert = sslCert
		server.SSLKey = sslKey
	}

	return server, nil
}

// Run server.
func (srv *APIServer) Run(ctx context.Context) {

	router := chi.NewRouter().With(
		middleware.Compress,
		middleware.Decompress,
		middleware.Cookies,
	)

	router.Get("/{linkId}", srv.handlerManager.GetLink)
	router.Post("/", srv.handlerManager.SetLink)
	router.Post("/api/shorten", srv.handlerManager.APISetLink)
	router.Get("/api/user/urls", srv.handlerManager.GetLinksByUser)
	router.Get("/ping", srv.handlerManager.PingDB)
	router.Post("/api/shorten/batch", srv.handlerManager.APISetBatchLinks)
	router.Delete("/api/user/urls", srv.handlerManager.DeleteUserLinks)

	srv.server.Handler = router
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		srv.server.Shutdown(ctxShutdown)
	}()

	go func() {
		if srv.useSSL {
			log.Println("server stated at localhost:443")
			err := srv.server.ListenAndServeTLS(srv.SSLCert, srv.SSLKey)
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("server stated at %s", srv.server.Addr))
			err := srv.server.ListenAndServe()
			log.Println(err.Error())
		}

	}()
}
