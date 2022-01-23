package linker

import (
	"context"
	"github.com/malyg1n/shortener/model"
)

// Linker is a business logic layer for work with links
type Linker interface {
	SetLink(ctx context.Context, link, userUUID string) (string, error)
	GetLink(ctx context.Context, id string) (string, error)
	GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error)
}
