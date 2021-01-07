package context

import (
	"fmt"
	"github.com/procyon-projects/procyon-configure"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type LoggingConfiguration interface {
	ApplyLoggingProperties(properties configure.LoggingProperties)
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
	Print(ctx interface{}, message interface{})
	Tracef(ctx interface{}, format string, args ...interface{})
	Debugf(ctx interface{}, format string, args ...interface{})
	Infof(ctx interface{}, format string, args ...interface{})
	Warningf(ctx interface{}, format string, args ...interface{})
	Errorf(ctx interface{}, format string, args ...interface{})
	Fatalf(ctx interface{}, format string, args ...interface{})
	Panicf(ctx interface{}, format string, args ...interface{})
	Printf(ctx interface{}, format string, args ...interface{})
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

func (l *SimpleLogger) ApplyLoggingProperties(properties configure.LoggingProperties) {
	loggerLevel := logrus.DebugLevel
	switch properties.Level {
	case "TRACE":
		loggerLevel = logrus.TraceLevel
	case "DEBUG":
		loggerLevel = logrus.DebugLevel
	case "INFO":
		loggerLevel = logrus.InfoLevel
	case "ERROR":
		loggerLevel = logrus.ErrorLevel
	case "WARNING":
		loggerLevel = logrus.WarnLevel
	case "FATAL":
		loggerLevel = logrus.FatalLevel
	case "PANIC":
		loggerLevel = logrus.PanicLevel
	default:
		loggerLevel = logrus.DebugLevel
	}
	l.log.Level = loggerLevel

	fullLogPath := properties.FilePath + properties.FileName
	if fullLogPath != "" {
		logFile, err := os.OpenFile(fullLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("Could Not Open Log File : " + err.Error())
		}
		l.log.SetOutput(logFile)
	} else {
		l.log.Formatter.(*LogFormatter).isTerminal = true
	}
}

func (l *SimpleLogger) Trace(ctx interface{}, message interface{}) {
	l.logMessage(ctx, TraceLevel, message)
}

func (l *SimpleLogger) Debug(ctx interface{}, message interface{}) {
	l.logMessage(ctx, DebugLevel, message)
}

func (l *SimpleLogger) Info(ctx interface{}, message interface{}) {
	l.logMessage(ctx, InfoLevel, message)
}

func (l SimpleLogger) Warning(ctx interface{}, message interface{}) {
	l.logMessage(ctx, WarnLevel, message)
}

func (l *SimpleLogger) Error(ctx interface{}, message interface{}) {
	l.logMessage(ctx, ErrorLevel, message)
}

func (l *SimpleLogger) Fatal(ctx interface{}, message interface{}) {
	l.logMessage(ctx, FatalLevel, message)
}

func (l *SimpleLogger) Panic(ctx interface{}, message interface{}) {
	l.logMessage(ctx, PanicLevel, message)
}

func (l *SimpleLogger) Print(ctx interface{}, message interface{}) {
	l.log.Print(message)
}

func (l *SimpleLogger) Tracef(ctx interface{}, format string, args ...interface{}) {
	l.logMessageWithFormat(ctx, TraceLevel, format, args...)
}

func (l *SimpleLogger) Debugf(ctx interface{}, format string, args ...interface{}) {
	l.logMessageWithFormat(ctx, DebugLevel, format, args...)
}

func (l *SimpleLogger) Infof(ctx interface{}, format string, args ...interface{}) {
	l.logMessageWithFormat(ctx, InfoLevel, format, args...)
}

func (l SimpleLogger) Warningf(ctx interface{}, format string, args ...interface{}) {
	l.logMessageWithFormat(ctx, WarnLevel, format, args...)
}

func (l *SimpleLogger) Errorf(ctx interface{}, format string, args ...interface{}) {
	l.logMessageWithFormat(ctx, ErrorLevel, format, args...)
}

func (l *SimpleLogger) Fatalf(ctx interface{}, format string, args ...interface{}) {
	l.logMessageWithFormat(ctx, FatalLevel, format, args...)
}

func (l *SimpleLogger) Panicf(ctx interface{}, format string, args ...interface{}) {
	l.logMessageWithFormat(ctx, PanicLevel, format, args...)
}

func (l *SimpleLogger) Printf(ctx interface{}, format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

func (l *SimpleLogger) logMessage(ctx interface{}, level LogLevel, message interface{}) {
	entry := l.getLogEntry(ctx, level)
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

func (l *SimpleLogger) logMessageWithFormat(ctx interface{}, level LogLevel, format string, args ...interface{}) {
	entry := l.getLogEntry(ctx, level)
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

func (l *SimpleLogger) getLogEntry(ctx interface{}, level LogLevel) *logrus.Entry {
	if ctx == nil {
		panic("Context cannot be nil")
	}

	var entry *logrus.Entry
	switch ctx.(type) {
	case Context:
		entry = l.log.WithFields(logrus.Fields{
			"CONTEXT_ID":  ctx.(Context).GetContextId(),
			"LEVEL_COLOR": l.getLevelColor(level),
		})
	case ContextId:
		entry = l.log.WithFields(logrus.Fields{
			"CONTEXT_ID":  ctx,
			"LEVEL_COLOR": l.getLevelColor(level),
		})
	default:
		panic("First parameter must be Context or Context Id")
	}
	return entry
}

func (l *SimpleLogger) getLevelColor(logLevel LogLevel) int {
	var levelColor int
	switch logLevel {
	case DebugLevel, TraceLevel:
		levelColor = 37 // gray
	case WarnLevel:
		levelColor = 33 // yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = 31 // red
	case InfoLevel:
		levelColor = 36 // blue
	default:
		levelColor = 37
	}
	return levelColor
}

type LogFormatter struct {
	logrus.TextFormatter
	isTerminal bool
	logFormat  string
}

func NewLogFormatter() *LogFormatter {
	formatter := &LogFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	formatter.logFormat = "[%s] %s %s : %s\n"
	return formatter
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	contextId := entry.Data["CONTEXT_ID"]
	if contextId == nil {
		return []byte(entry.Message + "\n"), nil
	}

	return []byte(
		fmt.Sprintf(f.logFormat, entry.Time.Format(f.TimestampFormat),
			f.GetLevelString(entry),
			f.GetSumContextId(contextId.(ContextId)),
			entry.Message)), nil
}

func (f *LogFormatter) GetSumContextId(contextId ContextId) string {
	var sumContextId = string(contextId)
	return sumContextId[:8]
}

func (f *LogFormatter) GetLevelString(entry *logrus.Entry) string {
	levelString := strings.ToUpper(entry.Level.String())
	if !f.isTerminal {
		return fmt.Sprintf("%-7s", levelString)
	}
	return fmt.Sprintf("\x1b[%dm%-7s\x1b[0m", entry.Data["LEVEL_COLOR"], levelString)
}
