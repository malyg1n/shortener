package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/services/linker"
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

// APISetLink get and store url
func (lh *LinksHandler) APISetLink(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := struct {
		URL string `json:"url"`
	}{}
	if err = json.Unmarshal(b, &s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	link, err := lh.service.SetLink(ctx, s.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := struct {
		Result string `json:"result"`
	}{Result: fmt.Sprintf("http://%s/%s", r.Host, link)}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}
