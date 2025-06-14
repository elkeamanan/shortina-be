package app

import (
	"context"
	"elkeamanan/shortina/config"
	"elkeamanan/shortina/util"

	log "github.com/sirupsen/logrus"
)

func Run() int {
	util.InitializeLogger()
	ctx := context.Background()
	err := config.LoadConfig(ctx)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	log.Info("successfully init config")

	repository, err := InitRepositories(ctx)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	log.Info("successfully init repositories")

	service, err := InitServices(repository)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	log.Info("successfully init services")

	err = InitServer(service)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	return 0
}
