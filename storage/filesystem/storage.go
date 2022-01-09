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
	_ = s.loadLinks()

	return s, nil
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

	link, ok := s.links[id]
	if !ok {
		return "", errs.ErrNotFound
	}

	return link, nil
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
