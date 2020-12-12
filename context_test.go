package context

import (
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testConfigurableContextAdapter struct {
}

func (adapter testConfigurableContextAdapter) Configure() {

}

func (adapter testConfigurableContextAdapter) OnConfigure() {

}

func (adapter testConfigurableContextAdapter) FinishConfigure() {

}

func TestBaseApplicationContext_WithNilConfigurableContextAdapter(t *testing.T) {
	assert.Panics(t, func() {
		NewBaseApplicationContext("app-id", "context-id", nil)
	})
}

func TestBaseApplicationContext_SetLogger(t *testing.T) {
	baseApplicationContext := NewBaseApplicationContext("app-id", "context-id", testConfigurableContextAdapter{})
	logger := NewSimpleLogger()
	baseApplicationContext.SetLogger(logger)
	assert.Panics(t, func() {
		baseApplicationContext.SetLogger(logger)
	})
	baseApplicationContext.AddApplicationListener(testApplicationListener1{})
	baseApplicationContext.AddApplicationListener(testApplicationListener2{})

	assert.Equal(t, 2, len(baseApplicationContext.GetApplicationListeners()))
}

func TestBaseApplicationContext_AddApplicationListener(t *testing.T) {
	baseApplicationContext := NewBaseApplicationContext("app-id", "context-id", testConfigurableContextAdapter{})
	assert.Panics(t, func() {
		baseApplicationContext.AddApplicationListener(nil)
	})
}

func TestBaseApplicationContext(t *testing.T) {
	baseApplicationContext := NewBaseApplicationContext("app-id", "context-id", testConfigurableContextAdapter{})
	baseApplicationContext.Get("")

	baseApplicationContext.AddApplicationListener(testApplicationListener1{})
	baseApplicationContext.AddApplicationListener(testApplicationListener2{})

	logger := NewSimpleLogger()
	baseApplicationContext.SetLogger(logger)
	assert.Equal(t, logger, baseApplicationContext.GetLogger())

	baseApplicationContext.SetApplicationName("app-name")
	assert.Equal(t, "app-name", baseApplicationContext.GetApplicationName())

	assert.Equal(t, ApplicationId("app-id"), baseApplicationContext.GetAppId())
	assert.Equal(t, ContextId("context-id"), baseApplicationContext.GetContextId())

	baseApplicationContext.Put("test-key", "test-value")
	assert.Equal(t, "test-value", baseApplicationContext.Get("test-key"))

	env := core.NewStandardEnvironment()
	baseApplicationContext.SetEnvironment(env)
	assert.Equal(t, env, baseApplicationContext.GetEnvironment())
	baseApplicationContext.Configure()
}
