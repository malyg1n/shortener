package filesystem

import (
	"context"
	"encoding/json"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/storage"
	"os"
	"sync"
)

var _ storage.LinksStorage = (*LinksStorageFile)(nil)

// LinksStorageFile structure.
type LinksStorageFile struct {
	mx       sync.RWMutex
	links    linkCollection
	filename string
}

type linkCollection struct {
	UserLinks map[string][]model.Link
	Links     map[string]string
}

// NewLinksStorageFile creates new LinksStorageMap instance.
func NewLinksStorageFile() (*LinksStorageFile, error) {
	cfg := config.GetConfig()
	s := &LinksStorageFile{
		filename: cfg.FileStoragePath,
		links: linkCollection{
			Links:     map[string]string{},
			UserLinks: map[string][]model.Link{},
		},
	}
	_ = s.loadLinks()

	return s, nil
}

// SetLink store link into collection.
func (s *LinksStorageFile) SetLink(ctx context.Context, id, link, userUUID string) {
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
func (s *LinksStorageFile) GetLink(ctx context.Context, id string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	link, ok := s.links.Links[id]
	if !ok {
		return "", errs.ErrNotFound
	}

	return link, nil
}

// GetLinksByUser returns links by user.
func (s *LinksStorageFile) GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	links, ok := s.links.UserLinks[userUUID]
	if !ok {
		return nil, errs.ErrNotFound
	}

	return links, nil
}

// GetLinkByOriginal returns links by original link.
func (s *LinksStorageFile) GetLinkByOriginal(ctx context.Context, url string) (string, error) {
	for k, v := range s.links.Links {
		if url == v {
			return k, nil
		}
	}

	return "", errs.ErrNotFound
}

// Close file handler.
func (s *LinksStorageFile) Close() error {
	return s.uploadLinks()
}

func (s *LinksStorageFile) loadLinks() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	dec := json.NewDecoder(file)

	if err := dec.Decode(&s.links); err != nil {
		return nil
	}

	return nil
}

func (s *LinksStorageFile) uploadLinks() error {
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()
	enc := json.NewEncoder(file)

	return enc.Encode(s.links)
}
