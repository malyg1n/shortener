package handlers

import (
	"github.com/malyg1n/shortener/internal/app/services"
	"github.com/malyg1n/shortener/internal/app/storage"
	"io"
	"net/http"
	"strings"
)

var (
	service services.LinksService
)

func init() {
	service = services.NewDefaultLinksService(storage.NewLinksStorageMap())
}

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
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	linkID, err := service.SetLink(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + r.Host + "/" + linkID))
}

func getLink(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	link, err := service.GetLink(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
