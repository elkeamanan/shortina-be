package handler

import (
	"bytes"
	"context"
	"elkeamanan/shortina/internal/user/domain"
	"encoding/json"
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
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				ErrorCode: 500,
				Message:   "Authorization header is required",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := domain.VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{
				ErrorCode: 401,
				Message:   "Invalid token",
			})
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
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{
				ErrorCode: 401,
				Message:   "Invalid token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), XAuthorizedUserIDHeaderKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
