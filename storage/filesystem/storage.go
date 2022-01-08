package filesystem

import (
	"context"
	"encoding/json"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/storage"
	"os"
	"sync"
)

var _ storage.LinksStorage = (*LinksStorageFile)(nil)

// LinksStorageFile structure
type LinksStorageFile struct {
	mx       sync.RWMutex
	links    map[string]string
	filename string
}

type link struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// NewLinksStorageFile creates new LinksStorageMap instance
func NewLinksStorageFile() (*LinksStorageFile, error) {
	cfg := config.GetConfig()
	s := &LinksStorageFile{
		filename: cfg.FileStoragePath,
		links:    make(map[string]string),
	}
	err := s.loadLinks()

	return s, err
}

// SetLink store link into collection
func (s *LinksStorageFile) SetLink(ctx context.Context, id string, link string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.links[id] = link
}

// GetLink returns link from collection by id
func (s *LinksStorageFile) GetLink(ctx context.Context, id string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if link, ok := s.links[id]; ok {
		return link, nil
	}

	return "", errs.ErrNotFound
}

// Close file handler
func (s *LinksStorageFile) Close() error {
	return s.uploadLinks()
}

func (s *LinksStorageFile) loadLinks() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	dec := json.NewDecoder(file)
	var links []link

	dec.Decode(&links)

	for _, link := range links {
		s.links[link.ID] = link.URL
	}

	return nil
}

func (s *LinksStorageFile) uploadLinks() error {
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	var links []link
	for linkID, linkURL := range s.links {
		links = append(links, link{ID: linkID, URL: linkURL})
	}

	return enc.Encode(links)
}
