package api

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenCookieName = "token"
const tokenIssuer = "go_final_project"
const tokenSubject = "scheduler"

func signinHandler(password string, signingKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w, http.MethodPost)
			return
		}
		if password == "" {
			writeJSON(w, map[string]string{"token": ""})
			return
		}

		var request struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeErrorStatus(w, http.StatusBadRequest, err)
			return
		}
		if subtle.ConstantTimeCompare([]byte(request.Password), []byte(password)) != 1 {
			writeErrorStatus(w, http.StatusUnauthorized, errors.New("неверный пароль"))
			return
		}

		claims := jwt.RegisteredClaims{
			Issuer:   tokenIssuer,
			Subject:  tokenSubject,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// ExpiresAt отсутствует по рекомендации ревьюера.
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(signingKey)
		if err != nil {
			writeErrorStatus(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, map[string]string{"token": signedToken})
	}
}

func auth(password string, signingKey []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if password == "" {
			next(w, r)
			return
		}
		cookie, err := r.Cookie(tokenCookieName)
		if err != nil || cookie.Value == "" || !validToken(cookie.Value, signingKey) {
			writeErrorStatus(w, http.StatusUnauthorized, errors.New("требуется авторизация"))
			return
		}
		next(w, r)
	}
}

func validToken(value string, signingKey []byte) bool {
	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(
		value,
		claims,
		func(token *jwt.Token) (any, error) {
			return signingKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(tokenIssuer),
		jwt.WithSubject(tokenSubject),
		jwt.WithIssuedAt(),
	)
	return err == nil && parsed != nil && parsed.Valid
}
