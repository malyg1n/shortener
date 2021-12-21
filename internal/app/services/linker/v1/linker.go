package v1

import (
	"context"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/internal/app/errs"
	"github.com/malyg1n/shortener/internal/app/services/linker"
	"github.com/malyg1n/shortener/internal/app/storage"
	"net/url"
	"regexp"
)

// DefaultLinker implements Linker
type DefaultLinker struct {
	storage storage.LinksStorage
	re      *regexp.Regexp
}

var _ linker.Linker = (*DefaultLinker)(nil)

// NewDefaultLinker creates new DefaultLinker instance
func NewDefaultLinker(storage storage.LinksStorage) (*DefaultLinker, error) {
	if storage == nil {
		return nil, errs.ErrStorageInternal
	}

	return &DefaultLinker{
		storage: storage,
		re:      regexp.MustCompile(`[a-zA-Z0-9-]+`),
	}, nil
}

// GetLink get link by id from storage
func (s *DefaultLinker) GetLink(ctx context.Context, id string) (string, error) {
	if !s.re.MatchString(id) {
		return "", errs.ErrInvalidInput
	}

	return s.storage.GetLink(ctx, id)
}

// SetLink store link into storage
func (s *DefaultLinker) SetLink(ctx context.Context, link string) (string, error) {
	if _, err := url.ParseRequestURI(link); err != nil {
		return "", errs.ErrInvalidInput
	}
	linkID := uuid.New().String()
	s.storage.SetLink(ctx, linkID, link)

	return linkID, nil
}
