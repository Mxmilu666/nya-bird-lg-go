package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

// InitLogger 初始化日志
func InitLogger() *log.Logger {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "nya-bird-lg",
	})
	return Logger
}

// 添加包级别的便捷函数
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	Logger.Fatal(msg, args...)
}
