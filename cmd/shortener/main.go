package main

import (
	"context"
	"fmt"
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
	defer ctxCancel()
	log.Println("run main")

	storage, err := pgsql.NewLinksStoragePG(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("run storage")

	defer storage.Close()

	service, err := v1.NewDefaultLinker(storage)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Println("run service")
	cfg := config.GetConfig()
	log.Println(fmt.Sprintf("%v, %v, %v, %v, %v", cfg.Addr, cfg.BaseURL, cfg.DatabaseDSN, cfg.FileStoragePath, cfg.SecretKey))

	server, err := rest.NewAPIServer(service, cfg.Addr)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("server created")

	server.Run(ctx)

	<-ctx.Done()
}
