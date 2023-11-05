package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
)

// LoggerManager 定义了日志管理器的结构
type LoggerManager struct {
	Logger *logrus.Logger
}

// NewLoggerManager 初始化一个LoggerManager对象
func NewLoggerManager(logPath string) *LoggerManager {
	// 创建logrus实例
	var baseLogger = logrus.New()

	// 设置日志输出格式为JSON
	baseLogger.SetFormatter(&logrus.JSONFormatter{})

	// 创建log目录
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			baseLogger.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// 设置输出
	logFile := filepath.Join(logPath, "app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		baseLogger.SetOutput(file)
	} else {
		baseLogger.Info("Failed to log to file, using default stderr")
	}

	return &LoggerManager{Logger: baseLogger}
}

func (l *LoggerManager) LogInfo(message string) {
	l.Logger.WithFields(logrus.Fields{
		"platform": runtime.GOOS,
		"arch":     runtime.GOARCH,
	}).Info(message)
}

func (l *LoggerManager) LogError(err error) {
	l.Logger.WithFields(logrus.Fields{
		"platform": runtime.GOOS,
		"arch":     runtime.GOARCH,
	}).Error(err)
}
