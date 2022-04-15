package storage

import (
	"context"
	"github.com/malyg1n/shortener/model"
)

// LinksStorage interface
type LinksStorage interface {
	SetLink(ctx context.Context, id, link, userUUID string) error
	GetLink(ctx context.Context, id string) (model.Link, error)
	GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error)
	GetLinkByOriginal(ctx context.Context, url string) (string, error)
	SetBatchLinks(ctx context.Context, links []model.Link, userUUID string) error
	MarkLinkAsRemoved(ctx context.Context, link model.Link) error
	Statistic(ctx context.Context) (int, int, error)
	Close() error
	Ping() error
}
