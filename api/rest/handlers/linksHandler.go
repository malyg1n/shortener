package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/models"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/services/linker"
	"io"
	"net/http"
	"strings"
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

	ctx := r.Context()
	userUUID := ctx.Value("user_uuid").(string)
	linkID, err := lh.service.SetLink(ctx, string(b), userUUID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(getFullURL(linkID)))
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
	dec := json.NewDecoder(r.Body)
	s := models.SetLinkRequest{}
	if err := dec.Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userUUID := ctx.Value("user_uuid").(string)
	linkID, err := lh.service.SetLink(ctx, s.URL, userUUID)

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

func (lh *LinksHandler) GetLinksByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID := ctx.Value("user_uuid").(string)

	links, err := lh.service.GetLinksByUser(ctx, userUUID)
	if err != nil {
		http.Error(w, "No content", http.StatusNoContent)
		return
	}

	responseLinks := make([]models.LinkResponse, len(links))
	for _, link := range links {
		fmt.Println(link)
		responseLinks = append(responseLinks, models.LinkResponse{
			ShortURL: getFullURL(link.ShortURL), OriginalURL: link.OriginalURL,
		})
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
