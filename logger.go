package context

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type LoggingProperties struct {
	Level string `yaml:"level" json:"level" default:"TRACE"`
	File  string `yaml:"file" json:"file"`
	Path  string `yaml:"path" json:"path"`
}

func newLoggingProperties() *LoggingProperties {
	return &LoggingProperties{}
}

func (properties *LoggingProperties) GetConfigurationPrefix() string {
	return "logging"
}

type LoggerConfiguration interface {
	ApplyLoggerProperties(properties LoggingProperties)
}

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
	Trace(ctx interface{}, message interface{})
	Debug(ctx interface{}, message interface{})
	Info(ctx interface{}, message interface{})
	Warning(ctx interface{}, message interface{})
	Error(ctx interface{}, message interface{})
	Fatal(ctx interface{}, message interface{})
	Panic(ctx interface{}, message interface{})
	Tracef(ctx interface{}, format string, args ...interface{})
	Debugf(ctx interface{}, format string, args ...interface{})
	Infof(ctx interface{}, format string, args ...interface{})
	Warningf(ctx interface{}, format string, args ...interface{})
	Errorf(ctx interface{}, format string, args ...interface{})
	Fatalf(ctx interface{}, format string, args ...interface{})
	Panicf(ctx interface{}, format string, args ...interface{})
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

func (l *SimpleLogger) ApplyLoggerProperties(properties LoggingProperties) {
	defaultLoggerLevel := logrus.DebugLevel
	switch properties.Level {
	case "TRACE":
		defaultLoggerLevel = logrus.TraceLevel
	case "DEBUG":
		defaultLoggerLevel = logrus.DebugLevel
	case "INFO":
		defaultLoggerLevel = logrus.InfoLevel
	case "ERROR":
		defaultLoggerLevel = logrus.ErrorLevel
	case "WARNING":
		defaultLoggerLevel = logrus.WarnLevel
	case "FATAL":
		defaultLoggerLevel = logrus.FatalLevel
	case "PANIC":
		defaultLoggerLevel = logrus.PanicLevel
	default:
		defaultLoggerLevel = logrus.TraceLevel
	}
	l.log.Level = defaultLoggerLevel
}

func (l *SimpleLogger) Trace(ctx interface{}, message interface{}) {
	l.logCtxMessage(ctx, TraceLevel, message)
}

func (l *SimpleLogger) Debug(ctx interface{}, message interface{}) {
	l.logCtxMessage(ctx, DebugLevel, message)
}

func (l *SimpleLogger) Info(ctx interface{}, message interface{}) {
	l.logCtxMessage(ctx, InfoLevel, message)
}

func (l SimpleLogger) Warning(ctx interface{}, message interface{}) {
	l.logCtxMessage(ctx, WarnLevel, message)
}

func (l *SimpleLogger) Error(ctx interface{}, message interface{}) {
	l.logCtxMessage(ctx, ErrorLevel, message)
}

func (l *SimpleLogger) Fatal(ctx interface{}, message interface{}) {
	l.logCtxMessage(ctx, FatalLevel, message)
}

func (l *SimpleLogger) Panic(ctx interface{}, message interface{}) {
	l.logCtxMessage(ctx, PanicLevel, message)
}

func (l *SimpleLogger) Tracef(ctx interface{}, format string, args ...interface{}) {
	l.logCtxMessageFormat(ctx, TraceLevel, format, args...)
}

func (l *SimpleLogger) Debugf(ctx interface{}, format string, args ...interface{}) {
	l.logCtxMessageFormat(ctx, DebugLevel, format, args...)
}

func (l *SimpleLogger) Infof(ctx interface{}, format string, args ...interface{}) {
	l.logCtxMessageFormat(ctx, InfoLevel, format, args...)
}

func (l SimpleLogger) Warningf(ctx interface{}, format string, args ...interface{}) {
	l.logCtxMessageFormat(ctx, WarnLevel, format, args...)
}

func (l *SimpleLogger) Errorf(ctx interface{}, format string, args ...interface{}) {
	l.logCtxMessageFormat(ctx, ErrorLevel, format, args...)
}

func (l *SimpleLogger) Fatalf(ctx interface{}, format string, args ...interface{}) {
	l.logCtxMessageFormat(ctx, FatalLevel, format, args...)
}

func (l *SimpleLogger) Panicf(ctx interface{}, format string, args ...interface{}) {
	l.logCtxMessageFormat(ctx, PanicLevel, format, args...)
}

func (l *SimpleLogger) logCtxMessage(ctx interface{}, level LogLevel, message interface{}) {
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
		entry.Trace(message)
	case DebugLevel:
		entry.Debug(message)
	case InfoLevel:
		entry.Info(message)
	case WarnLevel:
		entry.Warn(message)
	case ErrorLevel:
		entry.Error(message)
	case FatalLevel:
		entry.Fatal(message)
	case PanicLevel:
		entry.Panic(message)
	}
}

func (l *SimpleLogger) logCtxMessageFormat(ctx interface{}, level LogLevel, format string, args ...interface{}) {
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
		entry.Tracef(format, args...)
	case DebugLevel:
		entry.Debugf(format, args...)
	case InfoLevel:
		entry.Infof(format, args...)
	case WarnLevel:
		entry.Warnf(format, args...)
	case ErrorLevel:
		entry.Errorf(format, args...)
	case FatalLevel:
		entry.Fatalf(format, args...)
	case PanicLevel:
		entry.Panicf(format, args...)
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
	contextId := entry.Data["context_id"].(ContextId)
	if contextId != "" {
		contextIdStr := string(contextId)
		separatorIndex := strings.Index(contextIdStr, "-")
		logContextId = contextIdStr[:separatorIndex]
	}

	return []byte(
		fmt.Sprintf("[%s] \x1b[%dm%-7s\x1b[0m %s : %s\n", entry.Time.Format(f.TimestampFormat), levelColor, strings.ToUpper(entry.Level.String()), logContextId, entry.Message)), nil
}
