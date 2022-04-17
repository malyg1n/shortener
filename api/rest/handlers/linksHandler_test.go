package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/malyg1n/shortener/api/rest/middleware"
	"github.com/malyg1n/shortener/pkg/config"
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
	"time"
)

type HandlerSuite struct {
	suite.Suite
	handler *HandlerManager
	service linker.Linker
	storage storage.LinksStorage
}

func (s *HandlerSuite) SetupTest() {
	s.storage = inmemory.NewLinksStorageMap()
	s.service, _ = v1.NewDefaultLinker(s.storage)
	s.handler, _ = NewHandlerManager(s.service)
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

func (s *HandlerSuite) TestUserContextErrors() {
	endpoints := []struct {
		method  string
		url     string
		handler http.HandlerFunc
	}{
		{
			method:  "POST",
			url:     "/",
			handler: s.handler.SetLink,
		},
		{
			method:  "POST",
			url:     "/api/shorten",
			handler: s.handler.APISetLink,
		},
		{
			method:  "GET",
			url:     "/user/urls",
			handler: s.handler.GetLinksByUser,
		},
		{
			method:  "DELETE",
			url:     "/user/urls",
			handler: s.handler.DeleteUserLinks,
		},
		{
			method:  "POST",
			url:     "/api/shorten/batch",
			handler: s.handler.APISetBatchLinks,
		},
	}
	for _, en := range endpoints {
		req, err := http.NewRequest(en.method, en.url, strings.NewReader(""))
		assert.NoError(s.T(), err)

		rr := httptest.NewRecorder()
		en.handler.ServeHTTP(rr, req)
		assert.Equal(s.T(), rr.Code, http.StatusInternalServerError)
		assert.Equal(s.T(), strings.TrimSpace(rr.Body.String()), "could not parse user from token as string")
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

func (s *HandlerSuite) TestApiSetBatchLinks() {
	tests := []struct {
		name         string
		codeExpected int
		jsonBody     string
		error        string
	}{
		{
			"400",
			400,
			`{"url": "https://google.com"}`,
			"",
		},
		{
			"201",
			201,
			`[{"correlation_id": "1", "original_url": "https://google.com"}]`,
			"",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			payload := strings.NewReader(tt.jsonBody)

			ts := httptest.NewServer(s.getRouter())
			defer ts.Close()

			res, body := testRequest(t, ts, http.MethodPost, "/api/shorten/batch", payload, map[string]string{})
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

func (s *HandlerSuite) TestGetLinksByUserNoContent() {
	req, err := http.NewRequest("GET", "/user/urls", nil)
	assert.NoError(s.T(), err)

	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserKey, "abc123")
	req = req.WithContext(ctx)

	s.handler.GetLinksByUser(rr, req)
	assert.Equal(s.T(), rr.Code, http.StatusNoContent)
}

func (s *HandlerSuite) TestGetLinksByUser() {
	req, err := http.NewRequest("GET", "/user/urls", nil)
	assert.NoError(s.T(), err)
	userUUID := uuid.New().String()

	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.ContextUserKey, userUUID)
	req = req.WithContext(ctx)

	s.service.SetLink(ctx, "https://ya.ru", userUUID)

	s.handler.GetLinksByUser(rr, req)
	assert.Equal(s.T(), rr.Code, http.StatusOK)
}

func (s *HandlerSuite) TestCompressMiddleware() {
	s.T().Run("compress", func(t *testing.T) {
		ts := httptest.NewServer(s.getRouter())
		defer ts.Close()

		res, _ := testRequest(
			t,
			ts,
			http.MethodPost,
			"/api/shorten",
			strings.NewReader(`{"url": "https://google.com"}`),
			map[string]string{"Accept-Encoding": "gzip"},
		)
		defer func() {
			_ = res.Body.Close()
		}()

		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	})
}

func (s *HandlerSuite) TestDecompressMiddleware() {
	s.T().Run("compress", func(t *testing.T) {
		ts := httptest.NewServer(s.getRouter())
		defer ts.Close()

		res, body := testRequest(
			t,
			ts,
			http.MethodPost,
			"/",
			strings.NewReader(`https://google.com`),
			map[string]string{"Content-Encoding": "gzip"},
		)
		defer func() {
			_ = res.Body.Close()
		}()

		assert.Equal(t, 400, res.StatusCode)
		assert.Equal(t, "gzip: invalid header", body)
	})
}

func (s *HandlerSuite) TestExistsLink() {
	tests := []struct {
		name         string
		path         string
		codeExpected int
		link         string
	}{
		{
			"#1",
			"/",
			201,
			"https://google.com",
		},
		{
			"#2",
			"/",
			409,
			"https://google.com",
		},
		{
			"#3",
			"/api/shorten",
			409,
			`{"url": "https://google.com"}`,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			payload := strings.NewReader(tt.link)

			ts := httptest.NewServer(s.getRouter())
			defer ts.Close()

			res, _ := testRequest(t, ts, http.MethodPost, tt.path, payload, map[string]string{})
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.codeExpected, res.StatusCode)
		})
	}
}

func (s *HandlerSuite) TestDeleteLinks() {
	s.T().Run("delete", func(t *testing.T) {
		ts := httptest.NewServer(s.getRouter())
		defer ts.Close()

		res, _ := testRequest(
			t,
			ts,
			http.MethodDelete,
			"/api/user/urls",
			strings.NewReader(`["fake-string"]`),
			map[string]string{},
		)

		assert.Equal(t, 202, res.StatusCode)
		_ = res.Body.Close()

		ctx := context.Background()
		shortLinkID, _ := s.service.SetLink(ctx, "https://google.com", "fake_uuid")

		res, _ = testRequest(t, ts, http.MethodGet, "/"+shortLinkID, nil, map[string]string{})
		assert.Equal(t, 307, res.StatusCode)
		_ = res.Body.Close()

		s.service.DeleteLinks(ctx, []string{shortLinkID}, "fake_uuid")

		time.Sleep(time.Millisecond * 200)
		res, _ = testRequest(t, ts, http.MethodGet, "/"+shortLinkID, nil, map[string]string{})
		assert.Equal(t, 410, res.StatusCode)
		_ = res.Body.Close()
	})
}

func (s *HandlerSuite) TestStatistics() {
	tests := []struct {
		name         string
		ip           string
		net          string
		expectedCode int
	}{
		{
			name:         "valid",
			ip:           "127.0.0.1",
			net:          "127.0.0.0/16",
			expectedCode: 200,
		},
		{
			name:         "invalid ip",
			ip:           "127.1.0.1",
			net:          "127.0.0.0/16",
			expectedCode: 403,
		},
		{
			name:         "empty subnet",
			ip:           "127.0.0.1",
			net:          "",
			expectedCode: 403,
		},
		{
			name:         "empty ip",
			ip:           "",
			net:          "127.0.0.0/16",
			expectedCode: 403,
		},
	}
	cfg := config.GetConfig()
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(s.getRouter())
			defer ts.Close()
			headers := map[string]string{"X-Real-IP": tt.ip}
			cfg.TrustedSubnet = tt.net
			res, _ := testRequest(t, ts, http.MethodGet, "/api/internal/stats", nil, headers)
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.expectedCode, res.StatusCode)
		})
	}
}

func (s *HandlerSuite) getRouter() chi.Router {
	router := chi.NewRouter().With(middleware.Compress, middleware.Decompress, middleware.Cookies)
	router.Get("/{linkId}", s.handler.GetLink)
	router.Post("/", s.handler.SetLink)
	router.Post("/api/shorten", s.handler.APISetLink)
	router.Get("/user/urls", s.handler.GetLinksByUser)
	router.Post("/api/shorten/batch", s.handler.APISetBatchLinks)
	router.Delete("/api/user/urls", s.handler.DeleteUserLinks)
	router.With(middleware.CheckSubnet).Get("/api/internal/stats", s.handler.Statistic)

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
