package context

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type logWriter struct {
	logMessage string
}

func (writer *logWriter) Write(p []byte) (n int, err error) {
	writer.logMessage = string(p)
	return len(p), nil
}

func testLogMessage(t *testing.T, writer *logWriter, level string, message string) {
	assert.True(t, strings.Contains(writer.logMessage, level))
	assert.True(t, strings.Contains(writer.logMessage, message))
}

func TestSimpleLogger(t *testing.T) {
	contextId := ContextId("test-context-id")
	simpleLogger := NewSimpleLogger()
	logWriter := &logWriter{}
	simpleLogger.log.Out = logWriter
	simpleLogger.log.Level = logrus.TraceLevel
	simpleLogger.log.ExitFunc = func(i int) {

	}

	simpleLogger.Trace(contextId, "test message")
	testLogMessage(t, logWriter, "TRACE", "test message")

	simpleLogger.Debug(contextId, "test message")
	testLogMessage(t, logWriter, "DEBUG", "test message")

	simpleLogger.Error(contextId, "test message")
	testLogMessage(t, logWriter, "ERROR", "test message")

	simpleLogger.Info(contextId, "test message")
	testLogMessage(t, logWriter, "INFO", "test message")

	simpleLogger.Warning(contextId, "test message")
	testLogMessage(t, logWriter, "WARNING", "test message")

	simpleLogger.Fatal(contextId, "test message")
	testLogMessage(t, logWriter, "FATAL", "test message")

	assert.Panics(t, func() {
		simpleLogger.Panic(contextId, "test message")
	})
	testLogMessage(t, logWriter, "PANIC", "test message")

	simpleLogger.Tracef(contextId, "test message")
	testLogMessage(t, logWriter, "TRACE", "test message")

	simpleLogger.Debugf(contextId, "test message")
	testLogMessage(t, logWriter, "DEBUG", "test message")

	simpleLogger.Errorf(contextId, "test message")
	testLogMessage(t, logWriter, "ERROR", "test message")

	simpleLogger.Infof(contextId, "test message")
	testLogMessage(t, logWriter, "INFO", "test message")

	simpleLogger.Warningf(contextId, "test message")
	testLogMessage(t, logWriter, "WARNING", "test message")

	simpleLogger.Fatalf(contextId, "test message")
	testLogMessage(t, logWriter, "FATAL", "test message")

	assert.Panics(t, func() {
		simpleLogger.Panicf(contextId, "test message")
	})
	testLogMessage(t, logWriter, "PANIC", "test message")
}
