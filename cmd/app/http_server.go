package app

import (
	"elkeamanan/shortina/config"
	"elkeamanan/shortina/handler"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func InitServer(serviceContainer *ServiceContainer) error {
	server := handler.NewServer(mux.NewRouter(), serviceContainer.LinkService, serviceContainer.UserService)

	port := config.Cfg.HttpPort
	log.Info(fmt.Sprintf("Server is listening on port %d", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	return err
}
