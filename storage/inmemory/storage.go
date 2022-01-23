package inmemory

import (
	"context"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/storage"
	"sync"
)

var _ storage.LinksStorage = (*LinksStorageMap)(nil)

// LinksStorageMap structure
type LinksStorageMap struct {
	mx    sync.RWMutex
	links linkCollection
}

type linkCollection struct {
	UserLinks map[string][]model.Link
	Links     map[string]string
}

// NewLinksStorageMap creates new LinksStorageMap instance
func NewLinksStorageMap() *LinksStorageMap {
	return &LinksStorageMap{
		links: linkCollection{
			Links:     map[string]string{},
			UserLinks: map[string][]model.Link{},
		},
	}
}

// SetLink store link into collection.
func (s *LinksStorageMap) SetLink(ctx context.Context, id, link, userUUID string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.links.Links[id] = link
	linkModel := model.Link{
		ShortURL:    id,
		OriginalURL: link,
	}
	s.links.UserLinks[userUUID] = append(s.links.UserLinks[userUUID], linkModel)
}

// GetLink returns link from collection by id.
func (s *LinksStorageMap) GetLink(ctx context.Context, id string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	link, ok := s.links.Links[id]
	if !ok {
		return "", errs.ErrNotFound
	}

	return link, nil
}

// GetLinksByUser returns links by user.
func (s *LinksStorageMap) GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	links, ok := s.links.UserLinks[userUUID]
	if !ok {
		return nil, errs.ErrNotFound
	}

	return links, nil
}

// Close storage
func (s *LinksStorageMap) Close() error {
	return nil
}
