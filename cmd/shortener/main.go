package main

import (
	"context"
	"fmt"
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
		fmt.Println("cancel global context")
		cancel()
	}()

	if err := rest.RunServer(ctx); err != nil {
		fmt.Println(fmt.Sprintf("failed to serve:+%v\n", err))
		log.Printf("failed to serve:+%v\n", err)
	}
}
