package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/services/linker"
	"github.com/malyg1n/shortener/services/linker/v1"
	"github.com/malyg1n/shortener/storage"
	"github.com/malyg1n/shortener/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type HandlerSuite struct {
	suite.Suite
	handler *LinksHandler
	service linker.Linker
	storage storage.LinksStorage
}

func (s *HandlerSuite) SetupTest() {
	s.storage = inmemory.NewLinksStorageMap()
	s.service, _ = v1.NewDefaultLinker(s.storage)
	s.handler, _ = NewLinksHandler(s.service)
}

func TestLinksHandlers(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) TestGetLink() {
	ctx := context.Background()
	shortLinkID, _ := s.service.SetLink(ctx, "https://google.com", "fake_uuid")
	tests := []struct {
		name         string
		codeExpected int
		linkExpected string
		id           string
	}{
		{
			name:         "valid link",
			codeExpected: 307,
			linkExpected: "https://google.com",
			id:           shortLinkID,
		},
		{
			name:         "empty link",
			codeExpected: 405,
			linkExpected: "",
			id:           "",
		},
		{
			name:         "undefined link",
			codeExpected: 400,
			linkExpected: "",
			id:           "undefined",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(s.getRouter())
			defer ts.Close()

			res, _ := testRequest(t, ts, http.MethodGet, "/"+tt.id, nil, map[string]string{})
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.codeExpected, res.StatusCode)
			assert.Equal(t, tt.linkExpected, res.Header.Get("Location"))
		})
	}
}

func (s *HandlerSuite) TestSetLink() {
	tests := []struct {
		name         string
		codeExpected int
		link         string
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
		s.T().Run(tt.name, func(t *testing.T) {
			payload := strings.NewReader(tt.link)

			ts := httptest.NewServer(s.getRouter())
			defer ts.Close()

			res, _ := testRequest(t, ts, http.MethodPost, "/", payload, map[string]string{})
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.codeExpected, res.StatusCode)
		})
	}
}

func (s *HandlerSuite) TestApiSetLink() {
	tests := []struct {
		name         string
		codeExpected int
		jsonBody     string
		error        string
	}{
		{
			"valid link",
			201,
			`{"url": "https://google.com"}`,
			"",
		},
		{
			"invalid link",
			400,
			`{"url": "invalid link"}`,
			"invalid input",
		},
		{
			"empty link",
			400,
			`{"url": ""}`,
			"invalid input",
		},
		{
			"invalid json",
			400,
			`{"url": "https://google.com""}`,
			`invalid character '"' after object key:value pair`,
		},
		{
			"empty body",
			400,
			``,
			"EOF",
		},
		{
			"invalid param name",
			400,
			`{"uri": "https://google.com"}`,
			"invalid input",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			payload := strings.NewReader(tt.jsonBody)

			ts := httptest.NewServer(s.getRouter())
			defer ts.Close()

			res, body := testRequest(t, ts, http.MethodPost, "/api/shorten", payload, map[string]string{})
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.codeExpected, res.StatusCode)
			if tt.error != "" {
				assert.Equal(t, tt.error, body)
			}
		})
	}
}

func (s *HandlerSuite) TestCompressAndDecompressMiddlewares() {
	tests := []struct {
		name           string
		codeExpected   int
		body           io.Reader
		headers        map[string]string
		waitingHeaders map[string]string
		error          string
	}{
		{
			"#1",
			201,
			strings.NewReader(`{"url": "https://google.com"}`),
			map[string]string{"Accept-Encoding": "gzip"},
			map[string]string{"Content-Encoding": "gzip"},
			"",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(s.getRouter())
			defer ts.Close()

			res, body := testRequest(t, ts, http.MethodPost, "/api/shorten", tt.body, tt.headers)
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.codeExpected, res.StatusCode)
			if tt.error != "" {
				assert.Equal(t, tt.error, body)
			} else {
				for k, v := range tt.waitingHeaders {
					assert.Equal(t, v, res.Header.Get(k))
				}
			}
		})
	}
}

func (s *HandlerSuite) getRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Compress)
	router.Use(middleware.Decompress)
	router.Use(middleware.Cookies)
	router.Get("/{linkId}", s.handler.GetLink)
	router.Post("/", s.handler.SetLink)
	router.Post("/api/shorten", s.handler.APISetLink)

	return router
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, payload io.Reader, headers map[string]string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, payload)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp, strings.TrimSpace(string(respBody))
}
