package utils

import (
	"fmt"
)

// 定义终端颜色的常量
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

// Info 打印信息级别的日志
func Info(format string, args ...interface{}) {
	fmt.Printf(InfoColor, fmt.Sprintf(format, args...))
	fmt.Println()
}

// Notice 打印通知级别的日志
func Notice(format string, args ...interface{}) {
	fmt.Printf(NoticeColor, fmt.Sprintf(format, args...))
	fmt.Println()
}

// Warning 打印警告级别的日志
func Warning(format string, args ...interface{}) {
	fmt.Printf(WarningColor, fmt.Sprintf(format, args...))
	fmt.Println()
}

// Error 打印错误级别的日志
func Error(format string, args ...interface{}) {
	fmt.Printf(ErrorColor, fmt.Sprintf(format, args...))
	fmt.Println()
}

// Debug 打印调试级别的日志
func Debug(format string, args ...interface{}) {
	fmt.Printf(DebugColor, fmt.Sprintf(format, args...))
	fmt.Println()
}
