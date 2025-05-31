package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func InitializeLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}
