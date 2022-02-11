package v1

import (
	"context"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/services/linker"
	"github.com/malyg1n/shortener/storage"
	"net/url"
	"regexp"
)

// DefaultLinker implements Linker.
type DefaultLinker struct {
	storage storage.LinksStorage
	re      *regexp.Regexp
}

var _ linker.Linker = (*DefaultLinker)(nil)

// NewDefaultLinker creates new DefaultLinker instance.
func NewDefaultLinker(storage storage.LinksStorage) (*DefaultLinker, error) {
	if storage == nil {
		return nil, errs.ErrStorageInternal
	}

	return &DefaultLinker{
		storage: storage,
		re:      regexp.MustCompile(`[a-zA-Z0-9-]+`),
	}, nil
}

// GetLink get link by id from storage.
func (s *DefaultLinker) GetLink(ctx context.Context, id string) (string, error) {
	if !s.re.MatchString(id) {
		return "", errs.ErrInvalidInput
	}

	return s.storage.GetLink(ctx, id)
}

// SetLink store link into storage.
func (s *DefaultLinker) SetLink(ctx context.Context, link, userUUID string) (string, error) {
	if _, err := url.ParseRequestURI(link); err != nil {
		return "", errs.ErrInvalidInput
	}

	existsLink, err := s.GetLinkByOriginal(ctx, link)
	if err == nil && existsLink != "" {
		return "", errs.ErrLinkExists
	}

	linkID := uuid.New().String()
	err = s.storage.SetLink(ctx, linkID, link, userUUID)
	if err != nil {
		return "", err
	}

	return linkID, nil
}

func (s *DefaultLinker) SetBatchLinks(ctx context.Context, links []model.Link, userUUID string) ([]model.Link, error) {
	insertLinks := make([]model.Link, 0)
	for k, lnk := range links {
		shortLink, err := s.GetLinkByOriginal(ctx, lnk.OriginalURL)
		if err != nil {
			shortLink = uuid.New().String()
			insertLinks = append(insertLinks, model.Link{
				ShortURL: shortLink, OriginalURL: lnk.OriginalURL,
			})
		}
		links[k].ShortURL = shortLink
	}

	if len(insertLinks) > 0 {
		err := s.storage.SetBatchLinks(ctx, insertLinks, userUUID)
		if err != nil {
			return nil, err
		}
	}

	return links, nil
}

// GetLinksByUser returns links bu user uuid.
func (s *DefaultLinker) GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error) {
	return s.storage.GetLinksByUser(ctx, userUUID)
}

// GetLinkByOriginal returns short link by original url.
func (s *DefaultLinker) GetLinkByOriginal(ctx context.Context, url string) (string, error) {
	return s.storage.GetLinkByOriginal(ctx, url)
}

// PingStorage check availability storage.
func (s *DefaultLinker) PingStorage() error {
	return s.storage.Ping()
}
