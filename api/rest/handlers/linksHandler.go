package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/api/rest/models"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/pkg/errs"
	"io"
	"net/http"
	"strings"
)

// SetLink get and store url.
func (hm *HandlerManager) SetLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userUUID, ok := ctx.Value(middleware.ContextUserKey).(string)
	if !ok {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	exitStatus := http.StatusCreated
	linkID, err := hm.service.SetLink(ctx, string(b), userUUID)

	if err != nil {
		if !errors.Is(errs.ErrLinkExists, err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		linkID, err = hm.service.GetLinkByOriginal(ctx, string(b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		exitStatus = http.StatusConflict
	}

	w.WriteHeader(exitStatus)
	_, err = w.Write([]byte(getFullURL(linkID)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetLink redirects ro url.
func (hm *HandlerManager) GetLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "linkId")
	link, err := hm.service.GetLink(ctx, id)

	if err != nil {
		if errors.Is(errs.ErrLinkRemoved, err) {
			http.Error(w, err.Error(), http.StatusGone)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// APISetLink get and store url.
func (hm *HandlerManager) APISetLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dec := json.NewDecoder(r.Body)
	s := models.SetLinkRequest{}
	if err := dec.Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userUUID, ok := ctx.Value(middleware.ContextUserKey).(string)
	if !ok {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	exitStatus := http.StatusCreated
	linkID, err := hm.service.SetLink(ctx, s.URL, userUUID)

	if err != nil {
		if !errors.Is(errs.ErrLinkExists, err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		linkID, err = hm.service.GetLinkByOriginal(ctx, s.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		exitStatus = http.StatusConflict
	}

	res := models.SetLinkResponse{Result: getFullURL(linkID)}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(exitStatus)
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetLinksByUser returns links bu user cookie.
func (hm *HandlerManager) GetLinksByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userUUID, ok := ctx.Value(middleware.ContextUserKey).(string)
	if !ok {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

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
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// APISetBatchLinks generate links by collection.
func (hm *HandlerManager) APISetBatchLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userUUID, ok := ctx.Value(middleware.ContextUserKey).(string)
	if !ok {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	dec := json.NewDecoder(r.Body)
	inLinks := make([]models.SetBatchLinkRequest, 0)

	if err := dec.Decode(&inLinks); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	canonicalLinks := make([]model.Link, len(inLinks))
	for k, lnk := range inLinks {
		canonicalLinks[k] = model.Link{
			ShortURL:    "",
			OriginalURL: lnk.OriginalURL,
		}
	}

	links, err := hm.service.SetBatchLinks(ctx, canonicalLinks, userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	outLinks := make([]models.SetBatchLinkResponse, 0)
	for k, l := range inLinks {
		link := links[k]
		outLinks = append(outLinks, models.SetBatchLinkResponse{
			CorrelationID: l.CorrelationID,
			ShortURL:      getFullURL(link.ShortURL),
		})
	}

	result, err := json.Marshal(outLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteUserLinks delete links.
func (hm *HandlerManager) DeleteUserLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userUUID, ok := ctx.Value(middleware.ContextUserKey).(string)
	if !ok {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	ids := make([]string, 0)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	hm.service.DeleteLinks(ctx, ids, userUUID)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func getFullURL(linkID string) string {
	cfg := config.GetConfig()
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	return fmt.Sprintf("%s/%s", baseURL, linkID)
}
