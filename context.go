package context

import (
	"github.com/google/uuid"
	"github.com/procyon-projects/procyon-core"
	"github.com/procyon-projects/procyon-peas"
	"sync"
)

type ApplicationContext interface {
	peas.ConfigurablePeaFactory
	GetAppId() uuid.UUID
	GetContextId() uuid.UUID
	GetApplicationName() string
	GetStartupTimestamp() int64
}

type ConfigurableContext interface {
	SetLogger(logger core.Logger)
	GetLogger() core.Logger
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
	GetPeaFactory() peas.ConfigurablePeaFactory
	AddApplicationListener(listener ApplicationListener)
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
	ConfigurableContextAdapter
	appId                       uuid.UUID
	contextId                   uuid.UUID
	name                        string
	startupTimestamp            int64
	logger                      core.Logger
	environment                 core.ConfigurableEnvironment
	mu                          sync.RWMutex
	peaFactory                  peas.ConfigurablePeaFactory
	applicationEventBroadcaster ApplicationEventBroadcaster
	applicationListeners        []ApplicationListener
}

func NewGenericApplicationContext(appId uuid.UUID, contextId uuid.UUID, configurableContextAdapter ConfigurableContextAdapter) *GenericApplicationContext {
	if configurableContextAdapter == nil {
		panic("Configurable Context Adapter must not be null")
	}
	return &GenericApplicationContext{
		appId:                      appId,
		contextId:                  contextId,
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

func (ctx *GenericApplicationContext) GetAppId() uuid.UUID {
	return ctx.appId
}

func (ctx *GenericApplicationContext) GetContextId() uuid.UUID {
	return ctx.contextId
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

func (ctx *GenericApplicationContext) SetLogger(logger core.Logger) {
	ctx.logger = logger
}

func (ctx *GenericApplicationContext) GetLogger() core.Logger {
	return ctx.logger
}

func (ctx *GenericApplicationContext) AddApplicationListener(listener ApplicationListener) {
	if ctx.applicationEventBroadcaster != nil {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(listener)
	}
	ctx.applicationListeners = append(ctx.applicationListeners, listener)
}

func (ctx *GenericApplicationContext) GetApplicationListeners() []ApplicationListener {
	return ctx.applicationListeners
}

func (ctx *GenericApplicationContext) PublishEvent(event ApplicationEvent) {
	_ = ctx.applicationEventBroadcaster.BroadcastEvent(event)

}

func (ctx *GenericApplicationContext) Configure() {
	ctx.mu.Lock()
	/* pea processors */
	ctx.initPeaProcessors()
	/* application event broadcaster */
	ctx.initApplicationEventBroadcaster()
	/* custom configure */
	ctx.OnConfigure()
	/* application event listeners */
	ctx.initApplicationEventListeners()
	ctx.mu.Unlock()
}

func (ctx *GenericApplicationContext) initPeaProcessors() {

}

func (ctx *GenericApplicationContext) initApplicationEventBroadcaster() {
	ctx.applicationEventBroadcaster = NewSimpleApplicationEventBroadcasterWithFactory(ctx.peaFactory)
}

func (ctx *GenericApplicationContext) initApplicationEventListeners() {
	appListeners := ctx.GetApplicationListeners()
	for _, appListener := range appListeners {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(appListener)
	}
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

func (ctx *GenericApplicationContext) RegisterSharedPea(peaName string, sharedObject interface{}) error {
	return nil
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
