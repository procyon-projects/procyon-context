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
	BroadcastEvent(context ApplicationContext, event ApplicationEvent) error
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

func NewBaseApplicationEventBroadcasterWithFactory(factory peas.ConfigurablePeaFactory) *BaseApplicationEventBroadcaster {
	if factory == nil {
		panic("Pea Factory must not be null")
	}
	return &BaseApplicationEventBroadcaster{
		peaFactory:             factory,
		eventListenerRetriever: NewDefaultApplicationEventListenerRetrieverWithFactory(factory),
		mu:                     sync.RWMutex{},
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

func (broadcaster *BaseApplicationEventBroadcaster) BroadcastEvent(event ApplicationEvent) error {
	// do nothing
	return nil
}

func (broadcaster *BaseApplicationEventBroadcaster) GetApplicationListeners(context ApplicationContext, event ApplicationEvent) []ApplicationListener {
	broadcaster.mu.Lock()
	listeners := broadcaster.eventListenerRetriever.GetApplicationListeners()
	broadcaster.mu.Unlock()
	supportListeners := make([]ApplicationListener, 0)
	for _, listener := range listeners {
		if broadcaster.supportsEvent(context, listener, event) {
			supportListeners = append(supportListeners, listener)
		}
	}
	return supportListeners
}

func (broadcaster *BaseApplicationEventBroadcaster) supportsEvent(context ApplicationContext,
	listener ApplicationListener,
	event ApplicationEvent) bool {
	subscribedEvents := listener.SubscribeEvents(context)
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

func NewSimpleApplicationEventBroadcasterWithFactory(factory peas.ConfigurablePeaFactory) *SimpleApplicationEventBroadcaster {
	return &SimpleApplicationEventBroadcaster{
		BaseApplicationEventBroadcaster: NewBaseApplicationEventBroadcasterWithFactory(factory),
	}
}

func (broadcaster *SimpleApplicationEventBroadcaster) BroadcastEvent(context ApplicationContext, event ApplicationEvent) (err error) {
	listeners := broadcaster.GetApplicationListeners(context, event)
	for _, listener := range listeners {
		err = broadcaster.invokeListener(context, listener, event)
		if err != nil {
			break
		}
	}
	return nil
}

func (broadcaster *SimpleApplicationEventBroadcaster) invokeListener(context ApplicationContext,
	listener ApplicationListener,
	event ApplicationEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			broadcaster.logger.Fatal(context, fmt.Sprintf("while invoking an application listener, the error occurred : %s", r))
		}
	}()
	listener.OnApplicationEvent(context, event)
	return nil
}
