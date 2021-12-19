package handlers

import (
	"github.com/malyg1n/shortener/internal/app/services"
	"io"
	"net/http"
	"strings"
)

type BaseHandler struct {
	service services.LinksService
}

func NewHandlers(service services.LinksService) *BaseHandler {
	return &BaseHandler{
		service: service,
	}
}

// Resolve ...
func (bh *BaseHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bh.getLink(w, r)
	case http.MethodPost:
		bh.setLink(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (bh *BaseHandler) setLink(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	linkID, err := bh.service.SetLink(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + r.Host + "/" + linkID))
}

func (bh *BaseHandler) getLink(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	link, err := bh.service.GetLink(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
