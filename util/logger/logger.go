package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Initialize() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}
