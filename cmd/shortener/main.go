package main

import (
	"github.com/malyg1n/shortener/internal/app/api/rest"
	"log"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		log.Fatal(rest.RunServer())
	}()

	<-c
}
