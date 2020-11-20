package context

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type LogLevel uint32

const (
	PanicLevel LogLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger interface {
	Trace(ctx interface{}, args ...interface{})
	Debug(ctx interface{}, args ...interface{})
	Info(ctx interface{}, args ...interface{})
	Warning(ctx interface{}, args ...interface{})
	Error(ctx interface{}, args ...interface{})
	Fatal(ctx interface{}, args ...interface{})
	Panic(ctx interface{}, args ...interface{})
}

type SimpleLogger struct {
	log *logrus.Logger
}

func NewSimpleLogger() *SimpleLogger {
	log := &SimpleLogger{
		&logrus.Logger{
			Out:       os.Stdout,
			Formatter: NewLogFormatter(),
			Level:     logrus.InfoLevel,
		},
	}
	return log
}

func (l *SimpleLogger) Trace(ctx interface{}, args ...interface{}) {
	l.logCtx(ctx, TraceLevel, args)
}

func (l *SimpleLogger) Debug(ctx interface{}, args ...interface{}) {
	l.logCtx(ctx, DebugLevel, args)
}

func (l *SimpleLogger) Info(ctx interface{}, args ...interface{}) {
	l.logCtx(ctx, InfoLevel, args)
}

func (l SimpleLogger) Warning(ctx interface{}, args ...interface{}) {
	l.logCtx(ctx, WarnLevel, args)
}

func (l *SimpleLogger) Error(ctx interface{}, args ...interface{}) {
	l.logCtx(ctx, ErrorLevel, args)
}

func (l *SimpleLogger) Fatal(ctx interface{}, args ...interface{}) {
	l.logCtx(ctx, FatalLevel, args)
}

func (l *SimpleLogger) Panic(ctx interface{}, args ...interface{}) {
	l.logCtx(ctx, PanicLevel, args)
}

func (l *SimpleLogger) logCtx(ctx interface{}, level LogLevel, args ...interface{}) {
	if ctx == nil {
		panic("Context cannot be nil")
	}
	var entry *logrus.Entry
	switch ctx.(type) {
	case Context:
		entry = l.log.WithFields(logrus.Fields{
			"context_id": ctx.(Context).GetContextId(),
		})
	case ContextId:
		entry = l.log.WithFields(logrus.Fields{
			"context_id": ctx,
		})
	default:
		panic("First parameter must be Context or Context Id")
	}
	switch level {
	case TraceLevel:
		entry.Trace(args)
	case InfoLevel:
		entry.Info(args)
	case WarnLevel:
		entry.Warn(args)
	case ErrorLevel:
		entry.Error(args)
	case FatalLevel:
		entry.Fatal(args)
	case PanicLevel:
		entry.Panic(args)
	}
}

type LogFormatter struct {
	logrus.TextFormatter
}

func NewLogFormatter() *LogFormatter {
	formatter := &LogFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	return formatter
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 37 // gray
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}

	var logContextId = ""
	contextId := entry.Data["context_id"].(string)
	if contextId != "" {
		separatorIndex := strings.Index(contextId, "-")
		logContextId = logContextId[:separatorIndex]
	}

	return []byte(
		fmt.Sprintf("[%s] \x1b[%dm%-7s\x1b[0m %s : %s\n", entry.Time.Format(f.TimestampFormat), levelColor, strings.ToUpper(entry.Level.String()), logContextId, entry.Message)), nil
}
