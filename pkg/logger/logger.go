// Package logger 日志包，用于日志输出和打印
package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger

	namespace string
}

// Debug 打印 debug 级别日志
func (l *Logger) Debug(args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Debug(args...)
}

// Debugf 格式化字符串后打印 debug 级别日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Debugf(format, args...)
}

// Info 打印 info 级别日志
func (l *Logger) Info(args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Info(args...)
}

// Infof 格式化字符串后打印 info 级别日志
func (l *Logger) Infof(format string, args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Infof(format, args...)
}

// Warn 打印 warn 级别日志
func (l *Logger) Warn(args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Warn(args...)
}

// Warnf 格式化字符串后打印 warn 级别日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Warnf(format, args...)
}

// Error 打印错误日志
func (l *Logger) Error(args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Error(args...)
}

// Errorf 格式化字符串后打印错误日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Errorf(format, args...)
}

// Fatal 打印致命错误日志，打印后立即退出程序
func (l *Logger) Fatal(args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Fatal(args...)
}

// Fatalf 格式化字符串后打印致命错误日志，打印后立即退出程序
func (l *Logger) Fatalf(format string, args ...interface{}) {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	entry.Fatalf(format, args...)
}

// Fields type, used to pass to `WithFields`.
type Fields logrus.Fields

// WithField adds a field to the log entry, note that it doesn't log until you
// call Debug, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	return entry.WithField(key, value)
}

// WithFields adds a struct of fields to the log entry. All it does is call
// `WithField` for each `Field`.
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	entry := logrus.NewEntry(l.Logger)
	SetCallFrame(entry, l.namespace, 1)
	return entry.WithFields(fields)
}

// SetCallFrame 设定调用栈
func SetCallFrame(entry *logrus.Entry, namespace string, skip int) {
	// 获取调用栈的 文件、行号
	_, file, line, _ := runtime.Caller(skip + 1)
	pc, _, _, _ := runtime.Caller(skip + 2)
	funcDetail := runtime.FuncForPC(pc)
	var funcName string
	if funcDetail != nil {
		funcName = funcDetail.Name()
	}

	SetCallerFrameWithFileAndLine(entry, namespace, funcName, file, line)
}

type contextKey string

const (
	runtimeCaller contextKey = "ContextKeyRuntimeCaller"
)

// SetCallerFrameWithFileAndLine 设定调用栈
func SetCallerFrameWithFileAndLine(entry *logrus.Entry, namespace, functionName, file string, line int) {
	splitTarget := filepath.FromSlash("/" + namespace + "/")
	// 拆解文件名，移除项目所在路径和项目名称，只保留到项目内的文件路径
	filename := strings.SplitN(file, splitTarget, 2)
	// 如果拆解后出现问题，回退到完整路径
	if len(filename) < 2 {
		filename = []string{"", file}
	}

	// 设定 logrus.Entry 的上下文信息
	entry.Context = context.WithValue(context.Background(), runtimeCaller, &runtime.Frame{
		File:     filename[1],
		Line:     line,
		Function: functionName,
	})
}

// NewLogger 按需创建 logger 实例
func NewLogger(level logrus.Level, namespace string, logFilePath string, hook []logrus.Hook) *Logger {
	// 创建 logrus 实例
	log := logrus.New()
	if len(hook) > 0 {
		for _, h := range hook {
			log.Hooks.Add(h)
		}
	}

	// 设置日志级别
	log.SetFormatter(NewLogFileFormatter())
	log.SetReportCaller(true)
	// 设置日志级别
	log.Level = level

	if logFilePath != "" {
		// 初始化日志文件
		err := initLoggerFile(log, logFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Logger{Logger: log, namespace: namespace}
}

// initLoggerFile 初始化日志文件
func initLoggerFile(logger *logrus.Logger, logPath string) error {
	execPath, _ := os.Executable()
	// 获取日志文件目录
	logDir := filepath.Join(filepath.Dir(execPath), filepath.Dir(logPath))
	logPath = filepath.Join(logDir, filepath.Base(logPath))
	// 创建并设定日志目录权限为 755（用户完全权限，组不可写，其他无权限）
	err := os.MkdirAll(logDir, 0744)
	if err != nil {
		// 如果错误，则返回
		return fmt.Errorf("failed to create %s directory: %w", logDir, err)
	}

	// 创建日志文件
	// 获取日志文件路径状态
	stat, err := os.Stat(logPath)
	if err != nil {
		// 如果路径不存在，则创建文件
		if os.IsNotExist(err) {
			_, err2 := os.Create(logPath)
			if err2 != nil {
				return fmt.Errorf("failed to create %s log file: %w", logPath, err)
			}
			// 检查是否创建完毕
			_, err2 = os.Stat(logPath)
			if err2 != nil {
				return fmt.Errorf("failed to check %s log file: %w", logPath, err)
			}
		} else {
			// 否则返回错误
			return err
		}
	}
	if stat != nil && stat.IsDir() {
		return errors.New("path exists but it is a directory")
	}

	// 打开文件
	// 设定打开方式：如文件不存在，则创建；只写入模式；追加写入模式
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s log file: %w", logPath, err)
	}

	// 设定多重输出流：一个是标准输出，一个是文件写入输出
	mw := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(mw)
	return nil
}
