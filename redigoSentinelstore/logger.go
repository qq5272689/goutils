package sentinelstore

import (
	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func init() {
	Logger = log.New()
	Logger.Formatter = &log.JSONFormatter{}
}
