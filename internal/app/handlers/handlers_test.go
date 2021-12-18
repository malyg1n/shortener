package handlers

import (
	"github.com/malyg1n/shortener/internal/app/services"
	"github.com/malyg1n/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		code   int
	}{
		{
			"get",
			"GET",
			400,
		},
		{
			"post",
			"POST",
			400,
		},
		{
			"put",
			"PUT",
			405,
		},
		{
			"delete",
			"DELETE",
			405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(BaseHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}
		})
	}
}

func Test_getLink(t *testing.T) {
	srv := services.NewDefaultLinksService(storage.NewLinksStorageMap())
	shortLink, _ := srv.SetLink("https://google.com")
	tests := []struct {
		name string
		code int
		link string
		id   string
	}{
		{
			name: "valid link",
			code: 307,
			link: "https://google.com",
			id:   shortLink,
		},
		{
			name: "empty link",
			code: 400,
			link: "",
			id:   "",
		},
		{
			name: "undefined link",
			code: 400,
			link: "",
			id:   "undefined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(getLink)
			h.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.code, result.StatusCode)
		})
	}
}

func Test_setLink(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		value       string
		body        string
		contentType string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
