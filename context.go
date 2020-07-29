package context

import (
	"github.com/google/uuid"
	"github.com/procyon-projects/procyon-core"
	"github.com/procyon-projects/procyon-peas"
	"sync"
)

var (
	baseApplicationContextPool sync.Pool
)

func initBaseApplicationContextPool() {
	baseApplicationContextPool = sync.Pool{
		New: newBaseApplicationContext,
	}
}

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
	CloneContext(contextId uuid.UUID, factory peas.ConfigurablePeaFactory) ConfigurableContext
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

func newBaseApplicationContext() interface{} {
	return &BaseApplicationContext{}
}

func NewBaseApplicationContext(appId uuid.UUID, contextId uuid.UUID, configurableContextAdapter ConfigurableContextAdapter) *BaseApplicationContext {
	if configurableContextAdapter == nil {
		panic("Configurable Context Adapter must not be null")
	}
	return &BaseApplicationContext{
		appId:                      appId,
		contextId:                  contextId,
		mu:                         sync.RWMutex{},
		ConfigurableContextAdapter: configurableContextAdapter,
		peaFactory:                 peas.NewDefaultPeaFactory(nil),
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

func (ctx *BaseApplicationContext) SetLogger(logger core.Logger) {
	ctx.logger = logger
}

func (ctx *BaseApplicationContext) GetLogger() core.Logger {
	return ctx.logger
}

func (ctx *BaseApplicationContext) AddApplicationListener(listener ApplicationListener) {
	if ctx.applicationEventBroadcaster != nil {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(listener)
	}
	ctx.applicationListeners = append(ctx.applicationListeners, listener)
}

func (ctx *BaseApplicationContext) GetApplicationListeners() []ApplicationListener {
	return ctx.applicationListeners
}

func (ctx *BaseApplicationContext) PublishEvent(event ApplicationEvent) {
	_ = ctx.applicationEventBroadcaster.BroadcastEvent(event)

}

func (ctx *BaseApplicationContext) Configure() {
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

func (ctx *BaseApplicationContext) initPeaProcessors() {

}

func (ctx *BaseApplicationContext) initApplicationEventBroadcaster() {
	ctx.applicationEventBroadcaster = NewSimpleApplicationEventBroadcasterWithFactory(ctx.peaFactory)
}

func (ctx *BaseApplicationContext) initApplicationEventListeners() {
	appListeners := ctx.GetApplicationListeners()
	for _, appListener := range appListeners {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(appListener)
	}
}

func (ctx *BaseApplicationContext) GetPeaFactory() peas.ConfigurablePeaFactory {
	return ctx.peaFactory
}

func (ctx *BaseApplicationContext) GetPea(name string) (interface{}, error) {
	return nil, nil
}

func (ctx *BaseApplicationContext) GetPeaByNameAndType(name string, typ *core.Type) (interface{}, error) {
	return nil, nil
}

func (ctx *BaseApplicationContext) GetPeaByNameAndArgs(name string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (ctx *BaseApplicationContext) GetPeaByType(typ *core.Type) (interface{}, error) {
	return nil, nil
}

func (ctx *BaseApplicationContext) ContainsPea(name string) (interface{}, error) {
	return nil, nil
}

func (ctx *BaseApplicationContext) RegisterSharedPea(peaName string, sharedObject interface{}) error {
	return nil
}

func (ctx *BaseApplicationContext) GetSharedPea(peaName string) interface{} {
	return nil
}

func (ctx *BaseApplicationContext) ContainsSharedPea(peaName string) bool {
	return false
}

func (ctx *BaseApplicationContext) GetSharedPeaNames() []string {
	return nil
}

func (ctx *BaseApplicationContext) GetSharedPeaCount() int {
	return 0
}

func (ctx *BaseApplicationContext) AddPeaProcessor(processor peas.PeaProcessor) error {
	return nil
}

func (ctx *BaseApplicationContext) GetPeaProcessors() []peas.PeaProcessor {
	return nil
}

func (ctx *BaseApplicationContext) GetPeaProcessorsCount() int {
	return 0
}

func (ctx *BaseApplicationContext) RegisterScope(scopeName string, scope peas.PeaScope) error {
	return nil
}

func (ctx *BaseApplicationContext) GetRegisteredScopes() []string {
	return nil
}

func (ctx *BaseApplicationContext) GetRegisteredScope(scopeName string) peas.PeaScope {
	return nil
}

func (ctx *BaseApplicationContext) SetParentPeaFactory(parent peas.PeaFactory) {

}

func (ctx *BaseApplicationContext) ClonePeaFactory() peas.PeaFactory {
	return nil
}

func (ctx *BaseApplicationContext) RegisterTypeToScope(typ *core.Type, scope peas.PeaScope) error {
	return nil
}

func (ctx *BaseApplicationContext) CloneContext(contextId uuid.UUID, factory peas.ConfigurablePeaFactory) ConfigurableContext {
	cloneContext := baseApplicationContextPool.Get().(*BaseApplicationContext)
	cloneContext.appId = ctx.appId
	cloneContext.contextId = contextId
	cloneContext.name = ctx.name
	cloneContext.startupTimestamp = ctx.startupTimestamp
	cloneContext.mu = ctx.mu
	cloneContext.ConfigurableContextAdapter = ctx.ConfigurableContextAdapter
	cloneContext.peaFactory = factory
	cloneContext.environment = ctx.environment
	cloneContext.applicationListeners = ctx.applicationListeners
	cloneContext.applicationEventBroadcaster = ctx.applicationEventBroadcaster
	return cloneContext
}

func (ctx *BaseApplicationContext) PutToPool() {
	baseApplicationContextPool.Put(ctx)
}
