package context

import (
	"github.com/Rollcomp/procyon-core"
	"sync"
)

type ApplicationContext interface {
	GetApplicationName() string
	GetStartupTimestamp() int64
}

type ConfigurableContext interface {
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
	Configure()
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
}

func NewGenericApplicationContext() *GenericApplicationContext {
	return &GenericApplicationContext{
		mu: sync.RWMutex{},
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

func (ctx *GenericApplicationContext) CreateEnvironment() core.ConfigurableEnvironment {
	return core.NewStandardEnvironment()
}

func (ctx *GenericApplicationContext) Configure() {
	ctx.mu.Lock()
	// TODO complete this part
	ctx.mu.Unlock()
}
