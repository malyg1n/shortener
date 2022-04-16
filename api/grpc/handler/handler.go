package handler

import (
	"context"
	pb "github.com/malyg1n/shortener/api/grpc/proto"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/pkg/subnet"
	"github.com/malyg1n/shortener/services/linker"
	"regexp"
)

// LinkerHandler implements LinkerServer.
type LinkerHandler struct {
	pb.UnimplementedLinkerServer
	linker linker.Linker
}

var re *regexp.Regexp

func init() {
	re, _ = regexp.Compile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
}

// NewLinkerHandler LinkerHandler constructor.
func NewLinkerHandler(service linker.Linker) *LinkerHandler {
	return &LinkerHandler{
		linker: service,
	}
}

// SetLink implements set link interface.
func (h *LinkerHandler) SetLink(ctx context.Context, in *pb.SetLinkRequest) (*pb.SetLinkResponse, error) {
	var response pb.SetLinkResponse
	if !re.MatchString(in.UserID) {
		response.Error = "invalid user id"
		return &response, nil
	}

	shortLink, err := h.linker.SetLink(ctx, in.OriginalLink, in.UserID)
	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	response.ShortLink = shortLink

	return &response, nil
}

// GetLink implements get link interface.
func (h *LinkerHandler) GetLink(ctx context.Context, in *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {
	var response pb.GetLinkResponse

	originalLink, err := h.linker.GetLink(ctx, in.ShortLink)
	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	response.OriginalLink = originalLink
	return &response, nil
}

// GetUserLinks implements get user link interface.
func (h *LinkerHandler) GetUserLinks(ctx context.Context, in *pb.GetUserLinksRequest) (*pb.BaseLinksResponse, error) {
	var response pb.BaseLinksResponse
	if !re.MatchString(in.UserID) {
		response.Error = "invalid user id"
		return &response, nil
	}

	links, err := h.linker.GetLinksByUser(ctx, in.UserID)
	if err != nil {
		response.Error = "no content"
		return &response, nil
	}

	for _, lnk := range links {
		response.Links = append(response.Links, &pb.BaseLinkResponse{
			OriginalURL: lnk.OriginalURL,
			ShortURL:    lnk.ShortURL,
		})
	}

	return &response, nil
}

// SetBatchLinks implements set batch links interface.
func (h *LinkerHandler) SetBatchLinks(ctx context.Context, in *pb.CorrelationLinksRequest) (*pb.CorrelationLinksResponse, error) {
	var response pb.CorrelationLinksResponse
	if !re.MatchString(in.UserID) {
		response.Error = "invalid user id"
		return &response, nil
	}

	canonicalLinks := make([]model.Link, len(in.Links), len(in.Links))
	for k, lnk := range in.Links {
		canonicalLinks[k] = model.Link{
			ShortURL:    "",
			OriginalURL: lnk.OriginalURL,
		}
	}

	links, err := h.linker.SetBatchLinks(ctx, canonicalLinks, in.UserID)
	canonicalLinks = nil

	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	for k, l := range in.Links {
		lnk := &pb.CorrelationLinkResponse{
			CorrelationId: l.CorrelationId,
			ShortURL:      links[k].ShortURL,
		}
		response.Links = append(response.Links, lnk)
	}

	return &response, nil
}

// DeleteLinks implements delete links interface.
func (h *LinkerHandler) DeleteLinks(ctx context.Context, in *pb.DeleteLinksRequest) (*pb.DeleteLinksResponse, error) {
	var response pb.DeleteLinksResponse
	if !re.MatchString(in.UserID) {
		response.Error = "invalid user id"
		return &response, nil
	}

	h.linker.DeleteLinks(ctx, in.ShortLinks, in.UserID)

	return &response, nil
}

// Statistic implements statistic interface.
func (h *LinkerHandler) Statistic(ctx context.Context, in *pb.StatisticRequest) (*pb.StatisticResponse, error) {
	var response pb.StatisticResponse

	sn := config.GetConfig().TrustedSubnet
	if subnet.CheckSubnet(in.IP, sn) {
		response.Error = "forbidden"
		return &response, nil
	}

	users, links, err := h.linker.Statistic(ctx)
	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	response.Users = uint64(users)
	response.Urls = uint64(links)

	return &response, nil
}
