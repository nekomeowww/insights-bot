// Package logger 日志包，用于日志输出和打印
package logger

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

// LogFileFormatter defines the format for log file.
type LogFileFormatter struct {
	logrus.TextFormatter
	MinimumCallerDepth int
}

// NewLogFileFormatter return the log format for log file.
//
// eg: 2023-06-01T12:00:00 [info] [controllers/some_controller/code_file.go:99] foo key=value
func NewLogFileFormatter() *LogFileFormatter {
	return &LogFileFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		},
		MinimumCallerDepth: 0,
	}
}

// Format renders a single log entry for log file
//
// the original file log format is defined here: github.com/sirupsen/logrus/text_formatter.TextFormatter{}.Format().
func (f *LogFileFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(map[string]any)
	for k, v := range entry.Data {
		data[k] = v
	}

	keys := make([]string, 0, len(data))

	for k := range data {
		if k == "file" {
			continue
		}

		keys = append(keys, k)
	}

	if !f.DisableSorting {
		if nil != f.SortingFunc {
			f.SortingFunc(keys)
		} else {
			sort.Strings(keys)
		}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixStr := entry.Time.Format(timestampFormat) + " "
	var renderFunc func(a ...any) string

	switch entry.Level {
	case logrus.TraceLevel:
		renderFunc = color.FgGray.Render
	case logrus.DebugLevel:
		renderFunc = color.FgGreen.Render
	case logrus.InfoLevel:
		renderFunc = color.FgCyan.Render
	case logrus.WarnLevel:
		renderFunc = color.FgYellow.Render
	case logrus.ErrorLevel:
		renderFunc = color.FgRed.Render
	case logrus.FatalLevel:
		renderFunc = color.FgMagenta.Render
	case logrus.PanicLevel:
		renderFunc = color.FgMagenta.Render
	default:
		renderFunc = color.FgGray.Render
	}

	prefixStr += renderFunc("[", entry.Level.String(), "]")

	b.WriteString(prefixStr)
	if data["file"] != nil {
		b.WriteString(fmt.Sprintf(" [%s]", data["file"]))
		delete(data, "file")
	} else if entry.Context != nil {
		caller, _ := entry.Context.Value(runtimeCaller).(*runtime.Frame)
		if caller != nil {
			b.WriteString(fmt.Sprintf(" [%s:%d]", caller.File, caller.Line))
		}
	}

	if "" != entry.Message {
		b.WriteString(" " + entry.Message)
	}

	for _, key := range keys {
		value := data[key]
		appendKeyValue(b, key, value, f.QuoteEmptyFields)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

// appendKeyValue append value with key to data that to be appended to log file.
func appendKeyValue(b *bytes.Buffer, key string, value interface{}, QuoteEmptyFields bool) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	b.WriteString(key)
	b.WriteByte('=')
	appendValue(b, value, QuoteEmptyFields)
}

// appendValue append value to data used for method appendKeyValue.
func appendValue(b *bytes.Buffer, value interface{}, QuoteEmptyFields bool) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !needsQuoting(stringVal, QuoteEmptyFields) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

// needsQuoting check where text needs to be quoted.
func needsQuoting(text string, QuoteEmptyFields bool) bool {
	if QuoteEmptyFields && len(text) == 0 {
		return true
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}

	return false
}
