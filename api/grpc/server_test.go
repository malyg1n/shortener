package grpc

import (
	"context"
	pb "github.com/malyg1n/shortener/api/grpc/proto"
	v1 "github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkerService_SetLink(t *testing.T) {
	tests := []struct {
		name string
		link string
		err  string
	}{
		{
			name: "valid",
			link: "https://google.com",
			err:  "",
		},
		{
			name: "invalid",
			link: "invalid link",
			err:  "invalid input",
		},
		{
			name: "empty",
			link: "",
			err:  "invalid input",
		},
	}
	linker, _ := v1.NewDefaultLinker(inmemory.NewLinksStorageMap())
	server := NewLinkerService(linker)
	for _, tt := range tests {
		reg := &pb.SetLinkRequest{
			OriginalLink: tt.link,
			UserID:       1,
		}
		rsp, err := server.SetLink(context.Background(), reg)
		assert.NoError(t, err)
		assert.Equal(t, tt.err, rsp.Error)
	}
}
