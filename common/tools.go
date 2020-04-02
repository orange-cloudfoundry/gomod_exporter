package common

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// InitLogs - initialize log facility
func InitLogs(config *BaseConfig) {
	log.SetOutput(os.Stderr)
	if config.Log.JSON {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: config.Log.NoColor,
		})
	}
	log.SetLevel(log.ErrorLevel)
	switch strings.ToUpper(config.Log.Level) {
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	}
}
