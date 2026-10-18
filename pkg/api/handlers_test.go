package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnsupportedMethods(t *testing.T) {
	mux := http.NewServeMux()
	Init(mux, "", nil)

	cases := []struct {
		path, method, allow string
	}{
		{"/api/signin", http.MethodGet, "POST"},
		{"/api/nextdate", http.MethodPost, "GET"},
		{"/api/tasks", http.MethodPost, "GET"},
		{"/api/task/done", http.MethodGet, "POST"},
		{"/api/task", http.MethodPatch, "GET, POST, PUT, DELETE"},
	}
	for _, tc := range cases {
		t.Run(tc.path+"_"+tc.method, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			result := httptest.NewRecorder()
			mux.ServeHTTP(result, req)
			if result.Code != http.StatusMethodNotAllowed {
				t.Fatalf("want 405, got %d: %s", result.Code, result.Body.String())
			}
			if got := result.Header().Get("Allow"); got != tc.allow {
				t.Fatalf("Allow header: want %q, got %q", tc.allow, got)
			}
			var response map[string]string
			if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil || response["error"] == "" {
				t.Fatalf("invalid JSON error body: %q (%v)", result.Body.String(), err)
			}
		})
	}
}

func TestBadRequests(t *testing.T) {
	mux := http.NewServeMux()
	Init(mux, "", nil)
	cases := []struct {
		method, path, body string
	}{
		{http.MethodGet, "/api/task", ""},
		{http.MethodDelete, "/api/task", ""},
		{http.MethodPost, "/api/task", "bad-json"},
		{http.MethodPut, "/api/task", "bad-json"},
	}
	for _, tc := range cases {
		result := httptest.NewRecorder()
		mux.ServeHTTP(result, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
		if result.Code != http.StatusBadRequest {
			t.Fatalf("%s %s: want 400, got %d", tc.method, tc.path, result.Code)
		}
	}
}

func TestJSONEncodeFailureReturns500(t *testing.T) {
	result := httptest.NewRecorder()
	writeJSONStatus(result, http.StatusOK, map[string]any{"unsupported": make(chan int)})
	if result.Code != http.StatusInternalServerError {
		t.Fatalf("want 500, got %d", result.Code)
	}
}
