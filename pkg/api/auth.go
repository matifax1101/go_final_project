package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const tokenCookieName = "token"
const tokenLifetime = 8 * time.Hour

func signinHandler(w http.ResponseWriter, r *http.Request) {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		writeJSON(w, map[string]string{"token": ""})
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}
	if req.Password != password {
		writeError(w, fmt.Errorf("неверный пароль"))
		return
	}

	token, err := makeToken(password, time.Now().Add(tokenLifetime))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{"token": token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie(tokenCookieName)
		if err != nil || !validToken(cookie.Value, password) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func makeToken(password string, expires time.Time) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := struct {
		PasswordHash string `json:"ph"`
		Exp          int64  `json:"exp"`
	}{
		PasswordHash: passwordHash(password),
		Exp:          expires.Unix(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := header + "." + payloadEncoded
	signature := sign(unsigned, password)
	return unsigned + "." + signature, nil
}

func validToken(token, password string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	expected := sign(parts[0]+"."+parts[1], password)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var payload struct {
		PasswordHash string `json:"ph"`
		Exp          int64  `json:"exp"`
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return false
	}
	if payload.PasswordHash != passwordHash(password) {
		return false
	}
	return time.Now().Unix() < payload.Exp
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func sign(data, password string) string {
	mac := hmac.New(sha256.New, []byte(password))
	_, _ = mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
