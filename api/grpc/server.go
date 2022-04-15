package grpc

import (
	"context"
	pb "github.com/malyg1n/shortener/api/grpc/proto"
	"github.com/malyg1n/shortener/services/linker"
	"strconv"
)

type LinkerServer struct {
	pb.UnimplementedLinkerServer
	linker linker.Linker
}

func NewLinkerService(service linker.Linker) *LinkerServer {
	return &LinkerServer{
		linker: service,
	}
}

func (s *LinkerServer) SetLink(ctx context.Context, in *pb.SetLinkRequest) (*pb.SetLinkResponse, error) {
	var response pb.SetLinkResponse

	shortLink, err := s.linker.SetLink(ctx, in.OriginalLink, strconv.FormatUint(in.UserID, 10))
	response.ShortLink = shortLink

	if err != nil {
		response.Error = err.Error()
		response.ShortLink = ""
	}

	return &response, nil
}
