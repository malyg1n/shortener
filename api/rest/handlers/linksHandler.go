package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/api/rest/models"
	"github.com/malyg1n/shortener/pkg/config"
	"io"
	"net/http"
	"strings"
)

// SetLink get and store url.
func (hm *HandlerManager) SetLink(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userUUID := ctx.Value(middleware.ContextUserKey).(string)
	linkID, err := hm.service.SetLink(ctx, string(b), userUUID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(getFullURL(linkID)))
}

// GetLink redirects ro url.
func (hm *HandlerManager) GetLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "linkId")
	ctx := r.Context()
	link, err := hm.service.GetLink(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// APISetLink get and store url.
func (hm *HandlerManager) APISetLink(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	s := models.SetLinkRequest{}
	if err := dec.Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userUUID := ctx.Value(middleware.ContextUserKey).(string)
	linkID, err := hm.service.SetLink(ctx, s.URL, userUUID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := models.SetLinkResponse{Result: getFullURL(linkID)}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// GetLinksByUser returns links bu user cookie.
func (hm *HandlerManager) GetLinksByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID := ctx.Value(middleware.ContextUserKey).(string)

	links, err := hm.service.GetLinksByUser(ctx, userUUID)
	if err != nil {
		http.Error(w, "No content", http.StatusNoContent)
		return
	}

	responseLinks := make([]models.LinkResponse, len(links))
	for k, link := range links {
		responseLinks[k] = models.LinkResponse{
			ShortURL: getFullURL(link.ShortURL), OriginalURL: link.OriginalURL,
		}
	}

	result, err := json.Marshal(responseLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getFullURL(linkID string) string {
	cfg := config.GetConfig()
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	return fmt.Sprintf("%s/%s", baseURL, linkID)
}
