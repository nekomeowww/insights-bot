// Package logger 日志包，用于日志输出和打印
package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapField zapcore.Field

func (f ZapField) MatchValue() any {
	switch f.Type {
	case zapcore.UnknownType: // checked
		return f.Interface
	case zapcore.ArrayMarshalerType: // checked
		return f.Interface
	case zapcore.ObjectMarshalerType: // checked
		return f.Interface
	case zapcore.BinaryType: // checked
		return f.Interface
	case zapcore.BoolType: // checked
		return f.Integer != 0
	case zapcore.ByteStringType: // checked
		return f.Interface
	case zapcore.Complex128Type: // checked
		return f.Interface
	case zapcore.Complex64Type: // checked
		return f.Interface
	case zapcore.DurationType: // checked
		return time.Duration(f.Integer)
	case zapcore.Float64Type: // checked
		return f.Integer
	case zapcore.Float32Type: // checked
		return f.Integer
	case zapcore.Int64Type: // checked
		return f.Integer
	case zapcore.Int32Type: // checked
		return f.Integer
	case zapcore.Int16Type: // checked
		return f.Integer
	case zapcore.Int8Type: // checked
		return f.Integer
	case zapcore.StringType: // checked
		return f.String
	case zapcore.TimeType: // checked
		return time.Unix(0, f.Integer)
	case zapcore.TimeFullType: // checked
		return f.Interface
	case zapcore.Uint64Type: // checked
		return f.Integer
	case zapcore.Uint32Type: // checked
		return f.Integer
	case zapcore.Uint16Type: // checked
		return f.Integer
	case zapcore.Uint8Type: // checked
		return f.Integer
	case zapcore.UintptrType: // checked
		return f.Integer
	case zapcore.ReflectType: // checked
		return f.Interface
	case zapcore.NamespaceType: // checked
		return ""
	case zapcore.StringerType: // checked
		return f.Interface
	case zapcore.ErrorType: // checked
		return f.Interface
	case zapcore.SkipType: // checked
		return ""
	case zapcore.InlineMarshalerType: // checked
		return f.Interface
	}

	return f.Interface
}

type Logger struct {
	LogrusLogger *logrus.Logger
	ZapLogger    *zap.Logger

	namespace string
}

// Debug 打印 debug 级别日志。
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Debug(msg, fields...)

	entry := logrus.NewEntry(l.LogrusLogger)
	SetCallFrame(entry, l.namespace, 1)

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Debug(msg)
}

// Info 打印 info 级别日志。
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Info(msg, fields...)

	entry := logrus.NewEntry(l.LogrusLogger)
	SetCallFrame(entry, l.namespace, 1)

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Info(msg)
}

// Warn 打印 warn 级别日志。
func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Warn(msg, fields...)

	entry := logrus.NewEntry(l.LogrusLogger)
	SetCallFrame(entry, l.namespace, 1)

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Warn(msg)
}

// Error 打印错误日志。
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Error(msg, fields...)

	entry := logrus.NewEntry(l.LogrusLogger)
	SetCallFrame(entry, l.namespace, 1)

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Error(msg)
}

// Fatal 打印致命错误日志，打印后立即退出程序。
func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	l.ZapLogger.Fatal(msg, fields...)

	entry := logrus.NewEntry(l.LogrusLogger)
	SetCallFrame(entry, l.namespace, 1)

	for _, v := range fields {
		entry = entry.WithField(v.Key, ZapField(v).MatchValue())
	}

	entry.Fatal(msg)
}

// SetCallFrame 设定调用栈。
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

// SetCallerFrameWithFileAndLine 设定调用栈。
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

func zapCoreLevelToLogrusLevel(level zapcore.Level) logrus.Level {
	switch level {
	case zapcore.DebugLevel:
		return logrus.DebugLevel
	case zapcore.InfoLevel:
		return logrus.InfoLevel
	case zapcore.WarnLevel:
		return logrus.WarnLevel
	case zapcore.ErrorLevel:
		return logrus.ErrorLevel
	case zapcore.FatalLevel:
		return logrus.FatalLevel
	case zapcore.PanicLevel, zapcore.DPanicLevel:
		return logrus.PanicLevel
	case zapcore.InvalidLevel:
		return logrus.InfoLevel
	default:
		return logrus.InfoLevel
	}
}

// NewLogger 按需创建 logger 实例。
func NewLogger(level zapcore.Level, namespace string, logFilePath string, hook []logrus.Hook) (*Logger, error) {
	var err error
	if logFilePath == "" {
		logFilePath, err = autoCreateLogFile(logFilePath)
		if err != nil {
			return nil, err
		}
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.OutputPaths = []string{logFilePath}
	config.ErrorOutputPaths = []string{logFilePath}
	config.InitialFields = map[string]interface{}{
		"app_name": namespace,
	}

	zapLogger, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	logrusLogger := logrus.New()
	if len(hook) > 0 {
		for _, h := range hook {
			logrusLogger.Hooks.Add(h)
		}
	}

	logrusLogger.SetFormatter(NewLogFileFormatter())
	logrusLogger.SetReportCaller(true)
	logrusLogger.Level = zapCoreLevelToLogrusLevel(level)

	l := &Logger{
		LogrusLogger: logrusLogger,
		ZapLogger:    zapLogger,
		namespace:    namespace,
	}

	l.Info("logger init success for both logrus and zap",
		zap.String("logFilePath", logFilePath),
		zap.String("current_level", level.String()),
	)

	return l, nil
}

func autoCreateLogFile(logFilePathStr string) (string, error) {
	if logFilePathStr != "" {
		return logFilePathStr, nil
	}

	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	logDir := filepath.Join(filepath.Dir(execPath), "logs")
	logFilePath := filepath.Join(logDir, "insights-bot.log")

	err = os.MkdirAll(logDir, 0744)
	if err != nil {
		return "", fmt.Errorf("failed to create %s directory: %w", logDir, err)
	}

	stat, err := os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err2 := os.Create(logFilePath)
			if err2 != nil {
				return "", fmt.Errorf("failed to create %s log file: %w", logFilePath, err)
			}
		} else {
			return "", err
		}
	}
	if stat != nil && stat.IsDir() {
		return "", errors.New("path exists but it is a directory")
	}

	return logFilePath, nil
}
