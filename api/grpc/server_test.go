package grpc

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/malyg1n/shortener/api/grpc/proto"
	v1 "github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

func TestAPIServer_Run(t *testing.T) {
	linker, err := v1.NewDefaultLinker(inmemory.NewLinksStorageMap())
	assert.NoError(t, err)
	listen, err := net.Listen("tcp", ":32101")
	assert.NoError(t, err)
	server, err := NewAPIServer(linker, listen)
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		server.Run(ctx)
	}()

	conn, err := grpc.Dial(`:32101`, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	assert.NoError(t, err)
	c := pb.NewLinkerClient(conn)
	testRoutes(c, t)

	cancel()
}

func testRoutes(c pb.LinkerClient, t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	newLink := &pb.SetLinkRequest{
		OriginalLink: "https://example.com",
		UserID:       userID,
	}
	link, err := c.SetLink(ctx, newLink)
	assert.NoError(t, err)

	orLink, err := c.GetLink(ctx, &pb.GetLinkRequest{ShortLink: link.ShortLink})
	assert.NoError(t, err)
	assert.Equal(t, newLink.OriginalLink, orLink.OriginalLink)
}
