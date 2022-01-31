package storage

import (
	"context"
	"github.com/malyg1n/shortener/model"
)

// LinksStorage interface
type LinksStorage interface {
	SetLink(ctx context.Context, id, link, userUUID string)
	GetLink(ctx context.Context, id string) (string, error)
	GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error)
	GetLinkByOriginal(ctx context.Context, url string) (string, error)
	Close() error
}
