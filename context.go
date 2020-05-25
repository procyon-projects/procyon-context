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
	GetPeaFactory() peas.PeaFactory
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
	peaFactory       peas.PeaFactory
	ConfigurableContextAdapter
}

func NewGenericApplicationContext(configurableContextAdapter ConfigurableContextAdapter) *GenericApplicationContext {
	if configurableContextAdapter == nil {
		panic("Configurable Context Adapter must not be null")
	}
	return &GenericApplicationContext{
		mu:                         sync.RWMutex{},
		ConfigurableContextAdapter: configurableContextAdapter,
		peaFactory:                 peas.NewDefaultPeaFactory(),
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

func (ctx *GenericApplicationContext) GetPeaFactory() peas.PeaFactory {
	return ctx.peaFactory
}

func (ctx *GenericApplicationContext) GetPea(name string) (interface{}, error) {
	return nil, nil
}

func (ctx *GenericApplicationContext) GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error) {
	return nil, nil
}

func (ctx *GenericApplicationContext) GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (ctx *GenericApplicationContext) GetPeaByType(typ *core.Type) (interface{}, error) {
	return nil, nil
}

func (ctx *GenericApplicationContext) ContainsPea(name string) (interface{}, error) {
	return nil, nil
}
