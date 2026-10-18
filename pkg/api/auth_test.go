package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestSignInAndJWT(t *testing.T) {
	password := "test-password"
	signingKey := []byte("unit-test-key-32-bytes-long-abcdef")
	handler := signinHandler(password, signingKey)

	wrongMethod := httptest.NewRecorder()
	handler(wrongMethod, httptest.NewRequest(http.MethodGet, "/api/signin", nil))
	if wrongMethod.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /api/signin: want 405, got %d", wrongMethod.Code)
	}

	wrongPassword := httptest.NewRecorder()
	handler(wrongPassword, httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`{"password":"wrong"}`)))
	if wrongPassword.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password: want 401, got %d", wrongPassword.Code)
	}

	success := httptest.NewRecorder()
	handler(success, httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`{"password":"test-password"}`)))
	if success.Code != http.StatusOK {
		t.Fatalf("correct password: want 200, got %d", success.Code)
	}
	var payload struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(success.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Token == "" || !validToken(payload.Token, signingKey) {
		t.Fatal("JWT was not issued or verification failed")
	}
	claims := &jwt.RegisteredClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(payload.Token, claims)
	if err != nil {
		t.Fatal(err)
	}
	if claims.ExpiresAt != nil {
		t.Fatal("JWT unexpectedly has an expiration")
	}

	if validToken(payload.Token, []byte("another-32-byte-key-xxxxxxxxxxxx")) {
		t.Fatal("JWT accepted with a different signing key")
	}
	if validToken(payload.Token+"corruption", signingKey) {
		t.Fatal("corrupted JWT was accepted")
	}

	secureHandler := auth(password, signingKey, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	noCookie := httptest.NewRecorder()
	secureHandler(noCookie, httptest.NewRequest(http.MethodGet, "/api/tasks", nil))
	if noCookie.Code != http.StatusUnauthorized {
		t.Fatalf("request without token: want 401, got %d", noCookie.Code)
	}

	withCookie := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	withCookie.AddCookie(&http.Cookie{Name: tokenCookieName, Value: payload.Token})
	authorised := httptest.NewRecorder()
	secureHandler(authorised, withCookie)
	if authorised.Code != http.StatusNoContent {
		t.Fatalf("valid token: want 204, got %d", authorised.Code)
	}
}
