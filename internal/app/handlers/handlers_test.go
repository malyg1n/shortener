package handlers

import (
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
	links["valid_link"] = "https://google.com"
	links["invalid_link"] = "some text"
	tests := []struct {
		name     string
		code     int
		value    string
		location string
	}{
		{
			name:     "valid link",
			code:     307,
			value:    links["valid_link"],
			location: links["valid_link"],
		},
		{
			name:     "invalid link",
			code:     400,
			value:    links["invalid_link"],
			location: "",
		},
		{
			name:     "empty link",
			code:     400,
			value:    "",
			location: "",
		},
		{
			name:     "undefined link",
			code:     400,
			value:    "undefined",
			location: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
