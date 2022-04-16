package main

import (
	"context"
	"fmt"
	"github.com/malyg1n/shortener/api"
	"github.com/malyg1n/shortener/api/grpc"
	"github.com/malyg1n/shortener/api/rest"
	"github.com/malyg1n/shortener/pkg/config"
	v1 "github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/pgsql"
	"log"
	"net"
	"os/signal"
	"syscall"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func init() {
	if buildVersion == "" {
		buildVersion = "NA"
	}
	if buildDate == "" {
		buildDate = "NA"
	}
	if buildCommit == "" {
		buildCommit = "NA"
	}

	fmt.Println("Build version: " + buildVersion)
	fmt.Println("Build date: " + buildDate)
	fmt.Println("Build commit: " + buildCommit)
}

func main() {
	ctx, ctxCancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	storage, err := pgsql.NewLinksStoragePG(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	service, err := v1.NewDefaultLinker(storage)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cfg := config.GetConfig()
	var server api.Server

	if cfg.APIType == "grpc" {
		listen, err := net.Listen("tcp", cfg.Addr)
		if err != nil {
			log.Fatalf("%v", err)
		}

		server, err = grpc.NewAPIServer(service, listen)
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		server, err = rest.NewAPIServer(service, cfg.Addr, cfg.EnableHTTPS, cfg.SSLCert, cfg.SSLPrivateKey)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	server.Run(ctx)
	<-ctx.Done()

	storage.Close()
	ctxCancel()
}
