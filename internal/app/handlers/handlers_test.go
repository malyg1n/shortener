package handlers

import (
	"github.com/malyg1n/shortener/internal/app/services"
	"github.com/malyg1n/shortener/internal/app/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var handlers *BaseHandler

func TestMain(m *testing.M) {
	handlers = NewHandlers(services.NewDefaultLinksService(storage.NewLinksStorageMap()))
	m.Run()
}

func TestResolve(t *testing.T) {
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
			h := http.HandlerFunc(handlers.Resolve)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}
		})
	}
}

func Test_getLink(t *testing.T) {
	shortLink, _ := handlers.service.SetLink("https://google.com")
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
			h := http.HandlerFunc(handlers.getLink)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}
			if res.Header.Get("Location") != tt.link {
				t.Errorf("Expected header location %s, got %s", tt.link, res.Header.Get("Location"))
			}
		})
	}
}

func Test_setLink(t *testing.T) {
	tests := []struct {
		name string
		code int
		link string
	}{
		{
			"valid link",
			201,
			"https://google.com",
		},
		{
			"invalid link",
			400,
			"invalid link",
		},
		{
			"empty link",
			400,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := strings.NewReader(tt.link)
			request := httptest.NewRequest(http.MethodPost, "/", payload)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.setLink)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}
		})
	}
}
