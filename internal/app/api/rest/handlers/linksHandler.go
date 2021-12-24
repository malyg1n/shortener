package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/internal/app/errs"
	"github.com/malyg1n/shortener/internal/app/services/linker"
	"io"
	"net/http"
)

// LinksHandler is a base handler
type LinksHandler struct {
	service linker.Linker
}

// NewLinksHandler creates new LinksHandler instance
func NewLinksHandler(service linker.Linker) (*LinksHandler, error) {
	if service == nil {
		return nil, errs.ErrLinkerInternal
	}
	return &LinksHandler{
		service: service,
	}, nil
}

// SetLink get and store url
func (lh *LinksHandler) SetLink(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	linkID, err := lh.service.SetLink(ctx, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + r.Host + "/" + linkID))
}

// GetLink redirects ro url
func (lh *LinksHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "linkId")
	ctx := r.Context()
	link, err := lh.service.GetLink(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}