package handlers

import (
	v1 "github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHandlerManager(t *testing.T) {
	hm, err := NewHandlerManager(nil)
	assert.Error(t, err)
	assert.Nil(t, hm)

	storage := inmemory.NewLinksStorageMap()
	service, _ := v1.NewDefaultLinker(storage)
	hm, err = NewHandlerManager(service)
	assert.NoError(t, err)
	assert.NotNil(t, hm)
}
