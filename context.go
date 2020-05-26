package context

import (
	"github.com/Rollcomp/procyon-core"
	"github.com/Rollcomp/procyon-peas"
	"sync"
)

type ApplicationContext interface {
	peas.ConfigurablePeaFactory
	GetApplicationName() string
	GetStartupTimestamp() int64
}

type ConfigurableContext interface {
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
	GetPeaFactory() peas.ConfigurablePeaFactory
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
	name                        string
	startupTimestamp            int64
	environment                 core.ConfigurableEnvironment
	mu                          sync.RWMutex
	peaFactory                  peas.ConfigurablePeaFactory
	applicationEventBroadcaster ApplicationEventBroadcaster
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
	/* application event broadcaster */
	ctx.initApplicationEventBroadcaster()
	/* custom configure */
	ctx.OnConfigure()
	/* register application event listeners */
	ctx.registerApplicationEventListeners()
	ctx.mu.Unlock()
}

func (ctx *GenericApplicationContext) initApplicationEventBroadcaster() {
	ctx.applicationEventBroadcaster = NewSimpleApplicationEventBroadcaster(ctx.peaFactory)
}

func (ctx *GenericApplicationContext) registerApplicationEventListeners() {

}

func (ctx *GenericApplicationContext) GetPeaFactory() peas.ConfigurablePeaFactory {
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

func (ctx *GenericApplicationContext) RegisterSharedPea(peaName string, sharedObject interface{}) {

}

func (ctx *GenericApplicationContext) GetSharedPea(peaName string) interface{} {
	return nil
}

func (ctx *GenericApplicationContext) ContainsSharedPea(peaName string) bool {
	return false
}

func (ctx *GenericApplicationContext) GetSharedPeaNames() []string {
	return nil
}

func (ctx *GenericApplicationContext) GetSharedPeaCount() int {
	return 0
}
