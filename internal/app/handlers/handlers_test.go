package handlers

import (
	"github.com/malyg1n/shortener/internal/app/services"
	"github.com/malyg1n/shortener/internal/app/storage"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var handler *BaseHandler

func TestMain(m *testing.M) {
	handler = NewBaseHandler(services.NewDefaultLinksService(storage.NewLinksStorageMap()))
	m.Run()
}

func Test_GetLink(t *testing.T) {
	shortLink, _ := handler.service.SetLink("https://google.com")
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
			r := handler.NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()
			res, _ := testRequest(t, handler.GetLink, http.MethodGet, "/"+tt.id, nil)
			if res.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, res.StatusCode)
			}
			if res.Header.Get("Location") != tt.link {
				t.Errorf("Expected header location %s, got %s", tt.link, res.Header.Get("Location"))
			}
		})
	}
}

func Test_SetLink(t *testing.T) {
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
			r := handler.NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()
			res, _ := testRequest(t, handler.SetLink, http.MethodPost, "/", payload)
			if res.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, res.StatusCode)
			}
		})
	}
}

func testRequest(t *testing.T, handler http.HandlerFunc, method, path string, payload io.Reader) (*http.Response, string) {
	req := httptest.NewRequest(method, path, payload)

	w := httptest.NewRecorder()
	h := http.HandlerFunc(handler)
	h.ServeHTTP(w, req)
	res := w.Result()
	respBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Error(err.Error())
	}

	defer res.Body.Close()

	return res, string(respBody)
}
