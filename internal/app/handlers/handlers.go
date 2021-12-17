package handlers

import (
	"fmt"
	"github.com/malyg1n/shortener/internal/app/support"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const linkPattern = `[a-zA-Z0-9]+`

var links = make(map[string]string)

// BaseHandler ...
func BaseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getLink(w, r)
	case http.MethodPost:
		setLink(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func setLink(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	urlParam := string(b)
	//_, err = url.ParseRequestURI(urlParam)
	//if err != nil {
	//	http.Error(w, "Incorrect url param", http.StatusBadRequest)
	//	return
	//}
	randString := support.RandomString(6)
	links[randString] = urlParam
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("http://" + r.Host + "/" + randString))
	if err != nil {
		log.Fatal(err)
	}
}

func getLink(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	matched, _ := regexp.MatchString(linkPattern, id)
	if !matched {
		http.Error(w, "Invalid link ID", http.StatusBadRequest)
		return
	}
	link, ok := links[id]
	if !ok {
		http.Error(w, fmt.Sprintf("Not found link with ID %s", id), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
