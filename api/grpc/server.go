package grpc

import (
	"context"
	"github.com/malyg1n/shortener/api/grpc/handler"
	pb "github.com/malyg1n/shortener/api/grpc/proto"
	"github.com/malyg1n/shortener/services/linker"
	"google.golang.org/grpc"
	"log"
	"net"
)

// APIServer struct.
type APIServer struct {
	handlerManager *handler.LinkerHandler
	listener       net.Listener
}

// NewAPIServer creates new instance
func NewAPIServer(service linker.Linker, listener net.Listener) (*APIServer, error) {
	server := &APIServer{
		handlerManager: handler.NewLinkerHandler(service),
		listener:       listener,
	}
	return server, nil
}

// Run server.
func (srv *APIServer) Run(ctx context.Context) {
	s := grpc.NewServer()
	pb.RegisterLinkerServer(s, srv.handlerManager)

	go func() {
		log.Println(s.Serve(srv.listener))
	}()

	<-ctx.Done()
	s.GracefulStop()
}
