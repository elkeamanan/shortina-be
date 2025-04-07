package handler

import (
	"bytes"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("[REQ] %s %s | Query: %v | Body: %s", r.Method, r.URL.Path, r.URL.Query(), string(bodyBytes))

		next.ServeHTTP(w, r)
	})
}
