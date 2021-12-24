package inmemory

import (
	"context"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/storage"
	"sync"
)

var _ storage.LinksStorage = (*LinksStorageMap)(nil)

// LinksStorageMap structure
type LinksStorageMap struct {
	// key - short url id, value = original url
	links sync.Map
}

// NewLinksStorageMap creates new LinksStorageMap instance
func NewLinksStorageMap() *LinksStorageMap {
	return &LinksStorageMap{}
}

// SetLink store link into collection
func (s *LinksStorageMap) SetLink(ctx context.Context, id string, link string) {
	s.links.Store(id, link)
}

// GetLink returns link from collection by id
func (s *LinksStorageMap) GetLink(ctx context.Context, id string) (string, error) {
	if link, ok := s.links.Load(id); ok {
		return link.(string), nil
	}

	return "", errs.ErrNotFound
}
