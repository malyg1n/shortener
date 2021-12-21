package v1

import (
	"context"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/internal/app/errs"
	"github.com/malyg1n/shortener/internal/app/services/linker"
	"github.com/malyg1n/shortener/internal/app/storage"
	"net/http"
	"net/url"
	"regexp"
)

// DefaultLinker implements Linker
type DefaultLinker struct {
	storage storage.LinksStorage
}

var (
	_  linker.Linker = (*DefaultLinker)(nil)
	re               = regexp.MustCompile(`[a-zA-Z0-9-]+`)
)

// NewDefaultLinker creates new DefaultLinker instance
func NewDefaultLinker(storage storage.LinksStorage) *DefaultLinker {
	if storage == nil {
		panic("storage is not defined")
	}

	return &DefaultLinker{
		storage: storage,
	}
}

// GetLink get link by id from storage
func (s *DefaultLinker) GetLink(ctx context.Context, id string) (string, error) {
	if !re.MatchString(id) {
		return "", &errs.ValidationError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    `url is not valid`,
		}
	}

	return s.storage.GetLink(ctx, id)
}

// SetLink store link into storage
func (s *DefaultLinker) SetLink(ctx context.Context, link string) (string, error) {
	if _, err := url.ParseRequestURI(link); err != nil {
		return "", &errs.ValidationError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    `url is not valid`,
		}
	}
	linkID := uuid.New().String()
	s.storage.SetLink(ctx, linkID, link)

	return linkID, nil
}
