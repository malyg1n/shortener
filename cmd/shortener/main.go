package main

import (
	"context"
	"github.com/malyg1n/shortener/api/rest"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	go func() {
		<-c
		cancel()
	}()

	if err := rest.RunServer(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
