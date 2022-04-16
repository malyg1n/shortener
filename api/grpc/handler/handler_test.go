package handler

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/malyg1n/shortener/api/grpc/proto"
	v1 "github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testHandler *LinkerHandler
var svc, _ = v1.NewDefaultLinker(inmemory.NewLinksStorageMap())

func init() {
	testHandler = NewLinkerHandler(svc)
}

func TestLinkerHandler_SetLink(t *testing.T) {
	userID := uuid.NewString()
	tests := []struct {
		name   string
		link   string
		err    string
		userID string
	}{
		{
			name:   "valid",
			link:   "https://google.com",
			err:    "",
			userID: userID,
		},
		{
			name:   "double",
			link:   "https://google.com",
			err:    "link already saved",
			userID: userID,
		},
		{
			name:   "invalid",
			link:   "invalid link",
			err:    "invalid input",
			userID: userID,
		},
		{
			name:   "empty",
			link:   "",
			err:    "invalid user id",
			userID: "",
		},
	}
	for _, tt := range tests {
		reg := &pb.SetLinkRequest{
			OriginalLink: tt.link,
			UserID:       tt.userID,
		}
		rsp, err := testHandler.SetLink(context.Background(), reg)
		assert.NoError(t, err)
		assert.Equal(t, tt.err, rsp.Error)
	}
}

func TestLinkerHandler_GetLink(t *testing.T) {
	originalLink := "https://ya.ru"
	shortLink, err := svc.SetLink(context.Background(), originalLink, "123")
	assert.NoError(t, err)
	tests := []struct {
		name string
		link string
		err  string
	}{
		{
			name: "valid",
			link: shortLink,
			err:  "",
		},
		{
			name: "not found",
			link: "some-link",
			err:  "not found",
		},
	}
	for _, tt := range tests {
		reg := &pb.GetLinkRequest{
			ShortLink: tt.link,
		}
		rsp, err := testHandler.GetLink(context.Background(), reg)
		assert.NoError(t, err)
		assert.Equal(t, tt.err, rsp.Error)
	}
}

func TestLinkerHandler_SetBatchLinks(t *testing.T) {
	userID := uuid.NewString()
	cLinks := make([]*pb.CorrelationLinkRequest, 0, 2)
	cLinks = append(cLinks, &pb.CorrelationLinkRequest{
		CorrelationId: "1",
		OriginalURL:   "https://ya.me",
	})
	cLinks = append(cLinks, &pb.CorrelationLinkRequest{
		CorrelationId: "2",
		OriginalURL:   "https://ya.com",
	})
	batch := &pb.CorrelationLinksRequest{
		UserID: userID,
		Links:  cLinks,
	}

	rsp, err := testHandler.SetBatchLinks(context.Background(), batch)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rsp.Links))
	assert.Equal(t, "1", rsp.Links[0].CorrelationId)
	assert.Equal(t, "2", rsp.Links[1].CorrelationId)
}

func TestLinkerHandler_GetUserLinks(t *testing.T) {
	userID := uuid.NewString()
	svc.SetLink(context.Background(), "https://github.com", userID)

	rq1 := &pb.GetUserLinksRequest{UserID: userID}
	rq2 := &pb.GetUserLinksRequest{UserID: uuid.NewString()}

	rsp, err := testHandler.GetUserLinks(context.Background(), rq1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rsp.Links))
	assert.Equal(t, "https://github.com", rsp.Links[0].OriginalURL)

	rsp, err = testHandler.GetUserLinks(context.Background(), rq2)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(rsp.Links))
	assert.Equal(t, "no content", rsp.Error)
}

func TestLinkerHandler_DeleteLinks(t *testing.T) {
	userID := uuid.NewString()
	rq := &pb.DeleteLinksRequest{UserID: userID}
	rsp, err := testHandler.DeleteLinks(context.Background(), rq)
	assert.NoError(t, err)
	assert.Equal(t, "", rsp.Error)
}

func TestLinkerHandler_Statistic(t *testing.T) {
	sv, _ := v1.NewDefaultLinker(inmemory.NewLinksStorageMap())
	th := NewLinkerHandler(sv)
	userID1 := uuid.NewString()
	userID2 := uuid.NewString()
	sv.SetLink(context.Background(), "https://mail.ru11", userID1)
	sv.SetLink(context.Background(), "https://mail.ru22", userID1)
	sv.SetLink(context.Background(), "https://mail.ru33", userID2)

	tests := []struct {
		name   string
		subnet string
		ip     string
		users  uint64
		urls   uint64
		err    string
	}{
		{
			name:   "valid",
			subnet: "127.0.0.0/16",
			ip:     "127.0.0.1",
			users:  2,
			urls:   3,
			err:    "",
		},
	}
	for _, tt := range tests {
		reg := &pb.StatisticRequest{IP: tt.ip}
		rsp, err := th.Statistic(context.Background(), reg)
		assert.NoError(t, err)
		assert.Equal(t, tt.err, rsp.Error)
		assert.Equal(t, tt.users, rsp.Users)
		assert.Equal(t, tt.urls, rsp.Urls)
	}
}
