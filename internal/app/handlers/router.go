package handlers

import "github.com/go-chi/chi/v5"

// NewRouter ...
func (bh *BaseHandler) NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/{linkID}", bh.GetLink)
	r.Post("/", bh.SetLink)

	return r
}
