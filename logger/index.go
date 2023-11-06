package logger

import (
	"github.com/charmbracelet/log"
	"io"
	"os"
	"time"
)

var LogWriter io.Writer
var Logger *log.Logger

func init() {
	LogWriter = os.Stderr
	Logger = log.NewWithOptions(LogWriter, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})
}

func GetLogger() *log.Logger {
	return Logger
}
