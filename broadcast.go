package context

import (
	"fmt"
	"sync"
)

type ApplicationEventBroadcaster interface {
	RegisterApplicationListener(listener ApplicationListener)
	UnregisterApplicationListener(listener ApplicationListener)
	RemoveAllApplicationListeners()
	BroadcastEvent(context ApplicationContext, event ApplicationEvent)
}

type SimpleApplicationEventBroadcaster struct {
	logger           Logger
	eventListenerMap map[ApplicationEventId][]ApplicationListener
	mu               sync.RWMutex
}

func NewSimpleApplicationEventBroadcaster(logger Logger) *SimpleApplicationEventBroadcaster {
	return &SimpleApplicationEventBroadcaster{
		mu:               sync.RWMutex{},
		eventListenerMap: make(map[ApplicationEventId][]ApplicationListener, 0),
		logger:           logger,
	}
}

func (broadcaster *SimpleApplicationEventBroadcaster) RegisterApplicationListener(listener ApplicationListener) {
	broadcaster.mu.Lock()
	for _, eventId := range listener.SubscribeEvents() {
		broadcaster.eventListenerMap[eventId] = append(broadcaster.eventListenerMap[eventId], listener)
	}
	broadcaster.mu.Unlock()
}

func (broadcaster *SimpleApplicationEventBroadcaster) UnregisterApplicationListener(listener ApplicationListener) {
	broadcaster.mu.Lock()
	//.....
	broadcaster.mu.Unlock()
}

func (broadcaster *SimpleApplicationEventBroadcaster) RemoveAllApplicationListeners() {
	broadcaster.mu.Lock()
	broadcaster.eventListenerMap = make(map[ApplicationEventId][]ApplicationListener, 0)
	broadcaster.mu.Unlock()
}

func (broadcaster *SimpleApplicationEventBroadcaster) BroadcastEvent(context ApplicationContext, event ApplicationEvent) {
	broadcaster.mu.Lock()
	listeners := broadcaster.eventListenerMap[event.GetEventId()]
	broadcaster.mu.Unlock()
	for _, listener := range listeners {
		broadcaster.invokeListener(context, listener, event)
	}
}

func (broadcaster *SimpleApplicationEventBroadcaster) invokeListener(context ApplicationContext, listener ApplicationListener, event ApplicationEvent) {
	defer func() {
		if r := recover(); r != nil {
			broadcaster.logger.Fatal(context, fmt.Sprintf("while invoking an application listener, the error occurred : %s", r))
		}
	}()
	listener.OnApplicationEvent(context, event)
}
