package logger

import (
	log "github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

func InitLogger() {
	formatter := &prefixed.TextFormatter{}
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	formatter.ForceFormatting = true
	formatter.DisableColors = false
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05 000000"
}
