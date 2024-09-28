package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func GetLogger() *log.Logger {
	return log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Prefix:          "GIN API: ",
	})
}
