package handlers

import (
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/services/linker"
)

// HandlerManager manager of handlers.
type HandlerManager struct {
	service linker.Linker
}

// NewHandlerManager instance.
func NewHandlerManager(linker linker.Linker) (*HandlerManager, error) {
	if linker == nil {
		return nil, errs.ErrLinkerInternal
	}
	return &HandlerManager{
		service: linker,
	}, nil
}
