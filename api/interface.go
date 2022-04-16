package api

import "context"

// Server base interface.
type Server interface {
	Run(ctx context.Context)
}
