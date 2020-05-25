package context

import (
	"github.com/Rollcomp/procyon-core"
	"github.com/Rollcomp/procyon-peas"
	"sync"
)

type ApplicationContext interface {
	peas.PeaFactory
	GetApplicationName() string
	GetStartupTimestamp() int64
}

type ConfigurableContext interface {
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
}

type ConfigurableContextAdapter interface {
	Configure()
	OnConfigure()
}

type ConfigurableApplicationContext interface {
	ApplicationContext
	ConfigurableContext
}

type GenericApplicationContext struct {
	name             string
	startupTimestamp int64
	environment      core.ConfigurableEnvironment
	mu               sync.RWMutex
	ConfigurableContextAdapter
}

func NewGenericApplicationContext(configurableContextAdapter ConfigurableContextAdapter) *GenericApplicationContext {
	if configurableContextAdapter == nil {
		panic("Configurable Context Adapter must not be null")
	}
	return &GenericApplicationContext{
		mu:                         sync.RWMutex{},
		ConfigurableContextAdapter: configurableContextAdapter,
	}
}

func (ctx *GenericApplicationContext) SetApplicationName(name string) {
	ctx.name = name
}

func (ctx *GenericApplicationContext) GetApplicationName() string {
	return ctx.name
}

func (ctx *GenericApplicationContext) GetStartupTimestamp() int64 {
	return ctx.startupTimestamp
}

func (ctx *GenericApplicationContext) SetEnvironment(environment core.ConfigurableEnvironment) {
	ctx.environment = environment
}

func (ctx *GenericApplicationContext) GetEnvironment() core.ConfigurableEnvironment {
	return ctx.environment
}

func (ctx *GenericApplicationContext) Configure() {
	ctx.mu.Lock()
	// TODO complete this part
	ctx.OnConfigure()
	ctx.mu.Unlock()
}
