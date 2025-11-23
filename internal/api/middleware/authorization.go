package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"
	"strings"
)

func NewAuthorizationMiddleware(u string, p string) *Authorization {
	unifiEncoding := base64.StdEncoding.EncodeToString([]byte(u + ":" + p))

	return &Authorization{
		token: unifiEncoding,
	}
}

type Authorization struct {
	token string
}

func (a Authorization) Secure(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		bearerParts := strings.Split(authHeader, " ")
		if len(bearerParts) != 2 {
			http.Error(w, "malformed authorization header", http.StatusUnauthorized)
			return
		}

		if subtle.ConstantTimeCompare([]byte(a.token), []byte(bearerParts[1])) != 1 {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			slog.Warn("Authorization failed", slog.String("ip", r.RemoteAddr))
			return
		}

		next.ServeHTTP(w, r)
	}
}
