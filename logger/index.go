package logger

import (
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"path/filepath"
	"time"
)

var LogWriter io.Writer
var Logger *log.Logger

func InitLogger(filename string) {
	// 获取程序执行路径
	exePath, err := os.Executable()
	if err != nil {
		panic("failed to get executable path: " + err.Error())
	}
	// 构造日志文件路径为程序所在目录下的 log.txt
	logName := fmt.Sprintf("%s.txt", filename)
	logPath := filepath.Join(filepath.Dir(exePath), logName)

	// 打开或创建日志文件
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	// 将 LogWriter 设置为文件
	LogWriter := logFile
	Logger = log.NewWithOptions(LogWriter, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})
}

func GetLogger() *log.Logger {
	return Logger
}
