package storage

import "context"

// LinksStorage interface
type LinksStorage interface {
	SetLink(ctx context.Context, id string, link string)
	GetLink(ctx context.Context, id string) (string, error)
}
