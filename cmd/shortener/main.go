package main

import (
	"context"
	"github.com/malyg1n/shortener/api/rest"
	"log"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	if err := rest.RunServer(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
