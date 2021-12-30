package linker

import "context"

// Linker is a business logic layer for work with links
type Linker interface {
	SetLink(ctx context.Context, link string) (string, error)
	GetLink(ctx context.Context, id string) (string, error)
}
