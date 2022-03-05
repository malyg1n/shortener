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
	UserLinks    map[string][]model.Link
	Links        map[string]string
	DeletedLinks map[string]string
}

// NewLinksStorageMap creates new LinksStorageMap instance
func NewLinksStorageMap() *LinksStorageMap {
	return &LinksStorageMap{
		links: linkCollection{
			Links:        map[string]string{},
			UserLinks:    map[string][]model.Link{},
			DeletedLinks: map[string]string{},
		},
	}
}

// SetLink store link into collection.
func (s *LinksStorageMap) SetLink(ctx context.Context, id, link, userUUID string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.links.Links[id] = link
	linkModel := model.Link{
		ShortURL:    id,
		OriginalURL: link,
	}
	s.links.UserLinks[userUUID] = append(s.links.UserLinks[userUUID], linkModel)

	return nil
}

// GetLink returns link from collection by id.
func (s *LinksStorageMap) GetLink(ctx context.Context, id string) (model.Link, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	_, deleted := s.links.DeletedLinks[id]
	link := model.Link{
		ShortURL:  id,
		IsDeleted: deleted,
	}

	lnk, ok := s.links.Links[id]
	if !ok {
		return link, errs.ErrNotFound
	}
	link.OriginalURL = lnk

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

// GetLinkByOriginal returns links by original link.
func (s *LinksStorageMap) GetLinkByOriginal(ctx context.Context, url string) (string, error) {
	for k, v := range s.links.Links {
		if url == v {
			return k, nil
		}
	}

	return "", errs.ErrNotFound
}

// SetBatchLinks set links from collection.
func (s *LinksStorageMap) SetBatchLinks(ctx context.Context, links []model.Link, userUUID string) error {
	for _, link := range links {
		_, err := s.GetLinkByOriginal(ctx, link.OriginalURL)
		if err != nil {
			err := s.SetLink(ctx, link.ShortURL, link.OriginalURL, userUUID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// MarkLinkAsRemoved marks links as removed.
func (s *LinksStorageMap) MarkLinkAsRemoved(ctx context.Context, link model.Link) error {
	if _, ok := s.links.Links[link.ShortURL]; ok {
		if uLinks, ok := s.links.UserLinks[link.UserUUID]; ok {
			for _, mLink := range uLinks {
				if mLink.ShortURL == link.ShortURL {
					s.links.DeletedLinks[link.ShortURL] = mLink.OriginalURL
				}
			}
		}
	}
	return nil
}

// Close storage.
func (s *LinksStorageMap) Close() error {
	return nil
}

// Ping storage.
func (s *LinksStorageMap) Ping() error {
	return nil
}
