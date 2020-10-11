package context

import (
	"fmt"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	"sync"
)

type ApplicationEventBroadcaster interface {
	RegisterApplicationListener(listener ApplicationListener)
	RegisterApplicationListenerByPeaName(peaName string)
	UnregisterApplicationListener(listener ApplicationListener)
	UnregisterApplicationListenerByPeaName(peaName string)
	RemoveAllApplicationListeners()
	BroadcastEvent(context ApplicationContext, event ApplicationEvent)
}

type BaseApplicationEventBroadcaster struct {
	logger                 Logger
	peaFactory             peas.ConfigurablePeaFactory
	eventListenerRetriever ApplicationEventListenerRetriever
	mu                     sync.RWMutex
}

func NewBaseApplicationEventBroadcaster(logger Logger) *BaseApplicationEventBroadcaster {
	return &BaseApplicationEventBroadcaster{
		eventListenerRetriever: NewDefaultApplicationEventListenerRetriever(),
		mu:                     sync.RWMutex{},
		logger:                 logger,
	}
}

func NewBaseApplicationEventBroadcasterWithFactory(logger Logger, factory peas.ConfigurablePeaFactory) *BaseApplicationEventBroadcaster {
	if factory == nil {
		panic("Pea Factory must not be null")
	}
	return &BaseApplicationEventBroadcaster{
		peaFactory:             factory,
		eventListenerRetriever: NewDefaultApplicationEventListenerRetrieverWithFactory(factory),
		mu:                     sync.RWMutex{},
		logger:                 logger,
	}
}

func (broadcaster *BaseApplicationEventBroadcaster) RegisterApplicationListener(listener ApplicationListener) {
	broadcaster.mu.Lock()
	broadcaster.eventListenerRetriever.AddApplicationListener(listener)
	broadcaster.mu.Unlock()
}

func (broadcaster *BaseApplicationEventBroadcaster) RegisterApplicationListenerByPeaName(peaName string) {
	broadcaster.mu.Lock()
	broadcaster.eventListenerRetriever.AddApplicationListenerByPeaName(peaName)
	broadcaster.mu.Unlock()
}

func (broadcaster *BaseApplicationEventBroadcaster) UnregisterApplicationListener(listener ApplicationListener) {
	broadcaster.mu.Lock()
	broadcaster.eventListenerRetriever.RemoveApplicationListener(listener)
	broadcaster.mu.Unlock()
}

func (broadcaster *BaseApplicationEventBroadcaster) UnregisterApplicationListenerByPeaName(peaName string) {
	broadcaster.mu.Lock()
	broadcaster.eventListenerRetriever.RemoveApplicationListenerByPeaName(peaName)
	broadcaster.mu.Unlock()
}

func (broadcaster *BaseApplicationEventBroadcaster) RemoveAllApplicationListeners() {
	broadcaster.mu.Lock()
	broadcaster.eventListenerRetriever.RemoveAllApplicationListeners()
	broadcaster.mu.Unlock()
}

func (broadcaster *BaseApplicationEventBroadcaster) BroadcastEvent(context ApplicationContext, event ApplicationEvent) {
	// do nothing
}

func (broadcaster *BaseApplicationEventBroadcaster) GetApplicationListeners(event ApplicationEvent) []ApplicationListener {
	broadcaster.mu.Lock()
	listeners := broadcaster.eventListenerRetriever.GetApplicationListeners()
	broadcaster.mu.Unlock()
	supportListeners := make([]ApplicationListener, 0)
	for _, listener := range listeners {
		if broadcaster.supportsEvent(listener, event) {
			supportListeners = append(supportListeners, listener)
		}
	}
	return supportListeners
}

func (broadcaster *BaseApplicationEventBroadcaster) supportsEvent(listener ApplicationListener, event ApplicationEvent) bool {
	subscribedEvents := listener.SubscribeEvents()
	for _, subscribedEvent := range subscribedEvents {
		subscribedEventType := core.GetType(subscribedEvent)
		eventType := core.GetType(event)
		if core.IsInterface(subscribedEventType) && eventType.Typ.Implements(subscribedEventType.Typ) {
			return true
		} else if subscribedEventType.Typ == eventType.Typ {
			return true
		} else if core.IsStruct(subscribedEventType) && core.IsEmbeddedStruct(eventType, subscribedEventType) {
			return true
		}
	}
	return false
}

type SimpleApplicationEventBroadcaster struct {
	*BaseApplicationEventBroadcaster
}

func NewSimpleApplicationEventBroadcaster(logger Logger) *SimpleApplicationEventBroadcaster {
	return &SimpleApplicationEventBroadcaster{
		BaseApplicationEventBroadcaster: NewBaseApplicationEventBroadcaster(logger),
	}
}

func NewSimpleApplicationEventBroadcasterWithFactory(logger Logger, factory peas.ConfigurablePeaFactory) *SimpleApplicationEventBroadcaster {
	return &SimpleApplicationEventBroadcaster{
		BaseApplicationEventBroadcaster: NewBaseApplicationEventBroadcasterWithFactory(logger, factory),
	}
}

func (broadcaster *SimpleApplicationEventBroadcaster) BroadcastEvent(context ApplicationContext, event ApplicationEvent) {
	listeners := broadcaster.GetApplicationListeners(event)
	for _, listener := range listeners {
		broadcaster.invokeListener(context, listener, event)
	}
}

func (broadcaster *SimpleApplicationEventBroadcaster) invokeListener(context ApplicationContext,
	listener ApplicationListener,
	event ApplicationEvent) {
	defer func() {
		if r := recover(); r != nil {
			broadcaster.logger.Fatal(context, fmt.Sprintf("while invoking an application listener, the error occurred : %s", r))
		}
	}()
	listener.OnApplicationEvent(context, event)
}
