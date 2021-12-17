package server

import (
	"github.com/malyg1n/shortener/internal/app/handlers"
	"log"
	"net/http"
)

// Run ...
func Run() {
	http.HandleFunc("/", handlers.BaseHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
