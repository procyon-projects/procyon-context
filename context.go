package context

import (
	"github.com/google/uuid"
	core "github.com/procyon-projects/procyon-core"
	"github.com/procyon-projects/procyon-peas"
	"sync"
)

type Context interface {
	GetAppId() uuid.UUID
	GetContextId() uuid.UUID
	GetApplicationName() string
	GetStartupTimestamp() int64
}

type ApplicationContext interface {
	Context
	peas.ConfigurablePeaFactory
}

type ConfigurableContext interface {
	SetLogger(logger Logger)
	GetLogger() Logger
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
	GetPeaFactory() peas.ConfigurablePeaFactory
	AddApplicationListener(listener ApplicationListener)
	AddPeaFactoryProcessor(processor peas.PeaFactoryProcessor)
	Copy(cloneContext ConfigurableContext, contextId uuid.UUID)
}

type ConfigurableApplicationContext interface {
	ApplicationContext
	ConfigurableContext
}

type ConfigurableContextAdapter interface {
	Configure()
	OnConfigure()
}

type BaseApplicationContext struct {
	ConfigurableContextAdapter
	appId            uuid.UUID
	contextId        uuid.UUID
	name             string
	startupTimestamp int64
	logger           Logger
	environment      core.ConfigurableEnvironment
	mu               *sync.RWMutex
	peas.ConfigurablePeaFactory
	applicationEventBroadcaster ApplicationEventBroadcaster
	applicationListeners        []ApplicationListener
	peaFactoryProcessors        []peas.PeaFactoryProcessor
}

func NewBaseApplicationContext(appId uuid.UUID, contextId uuid.UUID, configurableContextAdapter ConfigurableContextAdapter) *BaseApplicationContext {
	if configurableContextAdapter == nil {
		panic("Configurable Context Adapter must not be null")
	}
	return &BaseApplicationContext{
		appId:                      appId,
		contextId:                  contextId,
		mu:                         &sync.RWMutex{},
		ConfigurableContextAdapter: configurableContextAdapter,
		ConfigurablePeaFactory:     peas.NewDefaultPeaFactory(nil),
		applicationListeners:       make([]ApplicationListener, 0),
		peaFactoryProcessors:       make([]peas.PeaFactoryProcessor, 0),
	}
}

func (ctx *BaseApplicationContext) SetApplicationName(name string) {
	ctx.name = name
}

func (ctx *BaseApplicationContext) GetApplicationName() string {
	return ctx.name
}

func (ctx *BaseApplicationContext) GetAppId() uuid.UUID {
	return ctx.appId
}

func (ctx *BaseApplicationContext) GetContextId() uuid.UUID {
	return ctx.contextId
}

func (ctx *BaseApplicationContext) GetStartupTimestamp() int64 {
	return ctx.startupTimestamp
}

func (ctx *BaseApplicationContext) SetEnvironment(environment core.ConfigurableEnvironment) {
	ctx.environment = environment
}

func (ctx *BaseApplicationContext) GetEnvironment() core.ConfigurableEnvironment {
	return ctx.environment
}

func (ctx *BaseApplicationContext) SetLogger(logger Logger) {
	if ctx.logger != nil {
		panic("There is already a logger, you cannot change it")
	}
	ctx.logger = logger
}

func (ctx *BaseApplicationContext) GetLogger() Logger {
	return ctx.logger
}

func (ctx *BaseApplicationContext) AddApplicationListener(listener ApplicationListener) {
	if listener == nil {
		panic("Listener must not be null")
	}
	if ctx.applicationEventBroadcaster != nil {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(listener)
	}
	ctx.applicationListeners = append(ctx.applicationListeners, listener)
}

func (ctx *BaseApplicationContext) AddPeaFactoryProcessor(processor peas.PeaFactoryProcessor) {
	if processor == nil {
		panic("Processor must not be null")
	}
	ctx.peaFactoryProcessors = append(ctx.peaFactoryProcessors, processor)
}

func (ctx *BaseApplicationContext) GetApplicationListeners() []ApplicationListener {
	return ctx.applicationListeners
}

func (ctx *BaseApplicationContext) PublishEvent(event ApplicationEvent) {
	ctx.applicationEventBroadcaster.BroadcastEvent(ctx, event)
}

func (ctx *BaseApplicationContext) Configure() {
	ctx.mu.Lock()
	ctx.preparePeaFactory()
	ctx.invokePeaFactoryPostProcessors()
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

func (ctx *BaseApplicationContext) preparePeaFactory() {
	peaFactory := ctx.GetPeaFactory()
	_ = peaFactory.RegisterSharedPea("environment", ctx.GetEnvironment())
}

func (ctx *BaseApplicationContext) invokePeaFactoryPostProcessors() {
	peaFactory := ctx.GetPeaFactory()
	if _, ok := peaFactory.(peas.PeaDefinitionRegistry); ok {

	}
}

func (ctx *BaseApplicationContext) initPeaProcessors() {

}

func (ctx *BaseApplicationContext) initApplicationEventBroadcaster() {
	ctx.applicationEventBroadcaster = NewSimpleApplicationEventBroadcasterWithFactory(ctx.logger, ctx.ConfigurablePeaFactory)
}

func (ctx *BaseApplicationContext) initApplicationEventListeners() {
	appListeners := ctx.GetApplicationListeners()
	for _, appListener := range appListeners {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(appListener)
	}
}

func (ctx *BaseApplicationContext) GetPeaFactory() peas.ConfigurablePeaFactory {
	return ctx.ConfigurablePeaFactory
}

func (ctx *BaseApplicationContext) Copy(cloneContext ConfigurableContext, contextId uuid.UUID) {
	if clone, ok := cloneContext.(*BaseApplicationContext); ok {
		clone.appId = ctx.appId
		clone.contextId = contextId
		clone.name = ctx.name
		clone.startupTimestamp = ctx.startupTimestamp
		clone.mu = ctx.mu
		clone.ConfigurableContextAdapter = ctx.ConfigurableContextAdapter
		clone.environment = ctx.environment
		clone.applicationListeners = ctx.applicationListeners
		clone.applicationEventBroadcaster = ctx.applicationEventBroadcaster
	}
}
