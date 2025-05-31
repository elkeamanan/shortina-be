package handler

import (
	"bytes"
	"context"
	"elkeamanan/shortina/internal/user/domain"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	AuthorizationHeaderKey     = "Authorization"
	XAuthorizedUserIDHeaderKey = "x-authorized-user-id"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("[REQ] %s %s | Query: %v | Body: %s", r.Method, r.URL.Path, r.URL.Query(), string(bodyBytes))

		next.ServeHTTP(w, r)
	})
}

func authVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeaderKey)
		if authHeader == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := domain.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), XAuthorizedUserIDHeaderKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func optionalAuthVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeaderKey)
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := domain.VerifyToken(tokenString)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), XAuthorizedUserIDHeaderKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
