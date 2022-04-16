package v1

import (
	"context"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	svc, _       = NewDefaultLinker(inmemory.NewLinksStorageMap())
	user1, user2 = uuid.NewString(), uuid.NewString()
)

func TestDefaultLinker_PingStorage(t *testing.T) {
	err := svc.PingStorage()
	assert.NoError(t, err)
}

func TestDefaultLinker_SetLink(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		hasError bool
	}{
		{
			name:     "valid",
			link:     "https://google.info",
			hasError: false,
		},
		{
			name:     "double",
			link:     "https://google.info",
			hasError: true,
		},
		{
			name:     "invalid",
			link:     "123....",
			hasError: true,
		},
	}
	for _, tt := range tests {
		_, err := svc.SetLink(context.Background(), tt.link, user1)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDefaultLinker_SetBatchLinks(t *testing.T) {
	tests := []struct {
		name  string
		links []model.Link
	}{
		{
			name: "valid",
			links: []model.Link{{
				OriginalURL: "https://exmo.com",
			}},
		},
		{
			name: "double",
			links: []model.Link{{
				OriginalURL: "https://google.info",
			}},
		},
	}

	for _, tt := range tests {
		res, err := svc.SetBatchLinks(context.Background(), tt.links, user2)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, tt.links[0].OriginalURL, res[0].OriginalURL)
	}
}

func TestDefaultLinker_GetLink(t *testing.T) {
	rLink, err := svc.SetLink(context.Background(), "https://booble.kom", user1)
	svc.DeleteLinks(context.Background(), []string{rLink}, user1)
	time.Sleep(time.Millisecond * 500)

	shLink, err := svc.SetLink(context.Background(), "https://gmail.com", user1)
	assert.NoError(t, err)
	tests := []struct {
		name  string
		link  string
		error string
	}{
		{
			name:  "valid",
			link:  shLink,
			error: "",
		},
		{
			name:  "invalid id",
			link:  "#$@!@##!",
			error: "invalid input",
		},
		{
			name:  "not exist",
			link:  uuid.NewString(),
			error: "not found",
		},
		{
			name:  "removed link",
			link:  rLink,
			error: "link was removed",
		},
	}

	for _, tt := range tests {
		_, err := svc.GetLink(context.Background(), tt.link)
		if tt.error != "" {
			assert.Error(t, err)
			assert.Equal(t, tt.error, err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDefaultLinker_GetLinkByOriginal(t *testing.T) {
	_, err := svc.SetLink(context.Background(), "https://mail.info", user1)
	assert.NoError(t, err)
	tests := []struct {
		name     string
		link     string
		hasError bool
	}{
		{
			name:     "valid",
			link:     "https://mail.info",
			hasError: false,
		},
		{
			name:     "invalid link",
			link:     "fake link 1",
			hasError: true,
		},
		{
			name:     "not exist",
			link:     "https://habr.com",
			hasError: true,
		},
	}

	for _, tt := range tests {
		_, err := svc.GetLinkByOriginal(context.Background(), tt.link)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDefaultLinker_GetLinksByUser(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		hasError bool
	}{
		{
			name:     "valid",
			user:     user1,
			hasError: false,
		},
		{
			name:     "invalid link",
			user:     uuid.NewString(),
			hasError: true,
		},
	}

	for _, tt := range tests {
		_, err := svc.GetLinksByUser(context.Background(), tt.user)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDefaultLinker_Statistic(t *testing.T) {
	x, y, err := svc.Statistic(context.Background())
	assert.NoError(t, err)
	assert.Greater(t, x, 0)
	assert.Greater(t, y, 0)
}

func TestDefaultLinker_NewDefaultLinker(t *testing.T) {
	_, err := NewDefaultLinker(nil)
	assert.Error(t, err)
}
