package app

import (
	"elkeamanan/shortina/config"
	"elkeamanan/shortina/handler"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func InitServer(serviceContainer ServiceContainer) error {
	server := handler.NewServer(mux.NewRouter())

	port := config.Cfg.HttpPort
	log.Info(fmt.Sprintf("Server is listening on port %d", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), enableCORS(server))
	return err
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
