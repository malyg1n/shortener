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
	_, _ = w.Write([]byte(getFullURL(linkID)))
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
	_, _ = w.Write(result)
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
	_, _ = w.Write(result)
}

// APISetBatchLinks generate links by collection.
func (hm *HandlerManager) APISetBatchLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID := ctx.Value(middleware.ContextUserKey).(string)

	dec := json.NewDecoder(r.Body)
	inLinks := make([]models.SetBatchLinkRequest, 0)

	if err := dec.Decode(&inLinks); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	outLinks := make([]models.SetBatchLinkResponse, 0)
	for _, l := range inLinks {
		link, err := hm.service.SetLink(ctx, l.OriginalURL, userUUID)
		if err == nil {
			outLinks = append(outLinks, models.SetBatchLinkResponse{
				CorrelationID: l.CorrelationID,
				ShortURL:      getFullURL(link),
			})
		}
	}

	result, err := json.Marshal(outLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(result)
}

func getFullURL(linkID string) string {
	cfg := config.GetConfig()
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	return fmt.Sprintf("%s/%s", baseURL, linkID)
}
