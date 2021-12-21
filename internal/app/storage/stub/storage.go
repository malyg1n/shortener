package stub

import (
	"context"
	"github.com/malyg1n/shortener/internal/app/errs"
	"github.com/malyg1n/shortener/internal/app/storage"
	"net/http"
	"sync"
)

var _ storage.LinksStorage = (*LinksStorageStub)(nil)

// LinksStorageStub structure
type LinksStorageStub struct {
	// key - short url id, value = original url
	links sync.Map
}

// NewLinksStorageStub creates new LinksStorageStub instance
func NewLinksStorageStub() *LinksStorageStub {
	return &LinksStorageStub{}
}

// SetLink store link into collection
func (s *LinksStorageStub) SetLink(ctx context.Context, id string, link string) {
	s.links.Store(id, link)
}

// GetLink returns link from collection by id
func (s *LinksStorageStub) GetLink(ctx context.Context, id string) (string, error) {
	if link, ok := s.links.Load(id); ok {
		return link.(string), nil
	}

	return "", &errs.NotFoundError{
		StatusCode: http.StatusNotFound,
		Message:    "link not found",
	}
}
