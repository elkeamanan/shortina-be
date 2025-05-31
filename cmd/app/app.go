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
	config.LoadConfig(ctx)

	repository, err := InitRepositories(ctx)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	service, err := InitServices(repository)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	err = InitServer(service)
	if err != nil {
		log.Error(err.Error())
		return 1
	}

	return 0
}
