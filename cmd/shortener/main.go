package main

import (
	"context"
	"github.com/malyg1n/shortener/api/rest"
	"github.com/malyg1n/shortener/pkg/config"
	v1 "github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/pgsql"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	storage, err := pgsql.NewLinksStoragePG(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	service, err := v1.NewDefaultLinker(storage)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cfg := config.GetConfig()

	server, err := rest.NewAPIServer(service, cfg.Addr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	server.Run(ctx)

	<-ctx.Done()

	storage.Close()
	ctxCancel()
}
