package context

import (
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
	"github.com/procyon-projects/procyon-peas"
	"sync"
)

const bootstrapProcessor = "github.com.procyon.context.bootstrapProcessor"
const eventListenerProcessor = "github.com.procyon.context.eventListenerProcessor"

type ApplicationId string
type ContextId string

type Context interface {
	GetContextId() ContextId
	Get(key string) interface{}
	Put(key string, value interface{})
}

type ApplicationContext interface {
	Context
	peas.ConfigurablePeaFactory
	GetAppId() ApplicationId
	GetApplicationName() string
	GetStartupTimestamp() int64
}

type ConfigurableContext interface {
	SetLogger(logger Logger)
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
	GetPeaFactory() peas.ConfigurablePeaFactory
	AddApplicationListener(listener ApplicationListener)
}

type ConfigurableApplicationContext interface {
	ApplicationContext
	ConfigurableContext
}

type ConfigurableContextAdapter interface {
	Configure()
	OnConfigure()
	FinishConfigure()
}

type BaseApplicationContext struct {
	ConfigurableContextAdapter
	peas.ConfigurablePeaFactory
	appId                       ApplicationId
	contextId                   ContextId
	name                        string
	startupTimestamp            int64
	logger                      Logger
	environment                 core.ConfigurableEnvironment
	mu                          *sync.RWMutex
	applicationEventBroadcaster ApplicationEventBroadcaster
	applicationListeners        []ApplicationListener
	bag                         map[string]interface{}
}

func NewBaseApplicationContext(appId ApplicationId, contextId ContextId, configurableContextAdapter ConfigurableContextAdapter) *BaseApplicationContext {
	if configurableContextAdapter == nil {
		panic("Configurable Context Adapter must not be null")
	}
	ctx := &BaseApplicationContext{
		appId:                      appId,
		contextId:                  contextId,
		mu:                         &sync.RWMutex{},
		ConfigurableContextAdapter: configurableContextAdapter,
		ConfigurablePeaFactory:     peas.NewDefaultPeaFactory(nil),
		applicationListeners:       make([]ApplicationListener, 0),
		bag:                        make(map[string]interface{}, 0),
	}
	ctx.initContext()
	return ctx
}

func (ctx *BaseApplicationContext) initContext() {
	peaDefinitionRegistry := ctx.GetPeaFactory().(peas.PeaDefinitionRegistry)
	if !peaDefinitionRegistry.ContainsPeaDefinition(bootstrapProcessor) {
		peaDefinition := peas.NewSimplePeaDefinition(goo.GetType(NewBootstrapProcessor))
		peaDefinitionRegistry.RegisterPeaDefinition(bootstrapProcessor, peaDefinition)
	}
	if !peaDefinitionRegistry.ContainsPeaDefinition(eventListenerProcessor) {
		peaDefinition := peas.NewSimplePeaDefinition(goo.GetType(NewEventListenerProcessor))
		peaDefinitionRegistry.RegisterPeaDefinition(eventListenerProcessor, peaDefinition)
	}
}

func (ctx *BaseApplicationContext) SetLogger(logger Logger) {
	if ctx.logger != nil {
		panic("There is already a logger, you cannot change it")
	}
	ctx.logger = logger
}

func (ctx *BaseApplicationContext) SetApplicationName(name string) {
	ctx.name = name
}

func (ctx *BaseApplicationContext) GetApplicationName() string {
	return ctx.name
}

func (ctx *BaseApplicationContext) GetAppId() ApplicationId {
	return ctx.appId
}

func (ctx *BaseApplicationContext) GetContextId() ContextId {
	return ctx.contextId
}

func (ctx *BaseApplicationContext) Get(key string) interface{} {
	return ctx.bag[key]
}

func (ctx *BaseApplicationContext) Put(key string, value interface{}) {
	ctx.bag[key] = value
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

func (ctx *BaseApplicationContext) AddApplicationListener(listener ApplicationListener) {
	if listener == nil {
		panic("Listener must not be null")
	}
	if ctx.applicationEventBroadcaster != nil {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(listener)
	}
	ctx.applicationListeners = append(ctx.applicationListeners, listener)
}

func (ctx *BaseApplicationContext) GetApplicationListeners() []ApplicationListener {
	return ctx.applicationListeners
}

func (ctx *BaseApplicationContext) PublishEvent(event ApplicationEvent) {
	ctx.applicationEventBroadcaster.BroadcastEvent(ctx, event)
}

func (ctx *BaseApplicationContext) Configure() {
	ctx.mu.Lock()
	err := ctx.preparePeaFactory()
	if err != nil {
		panic(err)
	}
	/* pea processors */
	ctx.initPeaProcessors()
	/* application event broadcaster */
	ctx.initApplicationEventBroadcaster()
	/* custom configure */
	ctx.OnConfigure()
	/* application event listeners */
	ctx.initApplicationEventListeners()
	/* finish pea factory initialization */
	ctx.finishPeaFactoryInitialization()
	/* finish the configure */
	ctx.FinishConfigure()
	ctx.mu.Unlock()
}

func (ctx *BaseApplicationContext) preparePeaFactory() (err error) {
	peaFactory := ctx.GetPeaFactory()
	err = peaFactory.RegisterSharedPea("environment", ctx.GetEnvironment())
	if err != nil {
		return err
	}
	err = peaFactory.RegisterSharedPea("logger", ctx.logger)
	if err != nil {
		return err
	}
	peaFactory.RegisterTypeAsOnlyReadable(goo.GetType((*ConfigurationProperties)(nil)))
	return
}

func (ctx *BaseApplicationContext) initPeaProcessors() {
	peaFactory := ctx.GetPeaFactory()
	if peaDefinitionRegistry, ok := peaFactory.(peas.PeaDefinitionRegistry); ok {
		// pea definition registry processor
		ctx.invokePeaDefinitionRegistryProcessors(ctx.getPeaDefinitionRegistryProcessors(peaDefinitionRegistry), peaDefinitionRegistry)
		// pea factory processor
		ctx.invokePeaFactoryProcessors(ctx.getPeaFactoryProcessors(peaDefinitionRegistry), peaFactory)
		// pea processors
		ctx.registerPeaProcessors(peaDefinitionRegistry)
	}
}

func (ctx *BaseApplicationContext) getPeaDefinitionRegistryProcessors(peaDefinitionRegistry peas.PeaDefinitionRegistry) []peas.PeaDefinitionRegistryProcessor {
	processors := make([]peas.PeaDefinitionRegistryProcessor, 0)
	peaFactory := ctx.GetPeaFactory()
	processorType := goo.GetType((*peas.PeaDefinitionRegistryProcessor)(nil))
	processorNames := peaDefinitionRegistry.GetPeaNamesForType(processorType)
	for _, processorName := range processorNames {
		instance, err := peaFactory.GetPeaByNameAndType(processorName, processorType)
		if err != nil {
			panic(err)
		}
		if instance != nil {
			processors = append(processors, instance.(peas.PeaDefinitionRegistryProcessor))
		}
	}
	return processors
}

func (ctx *BaseApplicationContext) invokePeaDefinitionRegistryProcessors(processors []peas.PeaDefinitionRegistryProcessor,
	registry peas.PeaDefinitionRegistry) {
	for _, processor := range processors {
		processor.AfterPeaDefinitionRegistryInitialization(registry)
	}
}

func (ctx *BaseApplicationContext) getPeaFactoryProcessors(peaDefinitionRegistry peas.PeaDefinitionRegistry) []peas.PeaFactoryProcessor {
	processors := make([]peas.PeaFactoryProcessor, 0)
	peaFactory := ctx.GetPeaFactory()
	processorType := goo.GetType((*peas.PeaFactoryProcessor)(nil))
	processorNames := peaDefinitionRegistry.GetPeaNamesForType(processorType)
	for _, processorName := range processorNames {
		instance, err := peaFactory.GetPeaByNameAndType(processorName, processorType)
		if err != nil {
			panic(err)
		}
		if instance != nil {
			processors = append(processors, instance.(peas.PeaFactoryProcessor))
		}
	}
	return processors
}

func (ctx *BaseApplicationContext) invokePeaFactoryProcessors(processors []peas.PeaFactoryProcessor,
	factory peas.ConfigurablePeaFactory) {
	for _, processor := range processors {
		processor.AfterPeaFactoryInitialization(factory)
	}
}

func (ctx *BaseApplicationContext) registerPeaProcessors(peaDefinitionRegistry peas.PeaDefinitionRegistry) {
	processors := make([]peas.PeaProcessor, 0)
	peaFactory := ctx.GetPeaFactory()
	processorType := goo.GetType((*peas.PeaProcessor)(nil))
	processorNames := peaDefinitionRegistry.GetPeaNamesForType(processorType)
	for _, processorName := range processorNames {
		instance, err := peaFactory.GetPeaByNameAndType(processorName, processorType)
		if err != nil {
			panic(err)
		}
		if instance != nil {
			processors = append(processors, instance.(peas.PeaProcessor))
		}
	}
	for _, processor := range processors {
		ctx.AddPeaProcessor(processor)
	}
}

func (ctx *BaseApplicationContext) initApplicationEventBroadcaster() {
	ctx.applicationEventBroadcaster = NewSimpleApplicationEventBroadcaster()
}

func (ctx *BaseApplicationContext) initApplicationEventListeners() {
	appListeners := ctx.GetApplicationListeners()
	for _, appListener := range appListeners {
		ctx.applicationEventBroadcaster.RegisterApplicationListener(appListener)
	}
}

func (ctx *BaseApplicationContext) finishPeaFactoryInitialization() {
	ctx.GetPeaFactory().PreInstantiateSharedPeas()
}

func (ctx *BaseApplicationContext) GetPeaFactory() peas.ConfigurablePeaFactory {
	return ctx.ConfigurablePeaFactory
}
