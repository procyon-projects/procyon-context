package context

import (
	"sync"
)

type ApplicationEventBroadcaster interface {
	RegisterApplicationListener(listener ApplicationListener)
	UnregisterApplicationListener(listener ApplicationListener)
	RemoveAllApplicationListeners()
	BroadcastEvent(context ApplicationContext, event ApplicationEvent)
}

type SimpleApplicationEventBroadcaster struct {
	eventListenerMap map[ApplicationEventId][]ApplicationListener
	mu               sync.RWMutex
}

func NewSimpleApplicationEventBroadcaster() *SimpleApplicationEventBroadcaster {
	return &SimpleApplicationEventBroadcaster{
		mu:               sync.RWMutex{},
		eventListenerMap: make(map[ApplicationEventId][]ApplicationListener, 0),
	}
}

func (broadcaster *SimpleApplicationEventBroadcaster) RegisterApplicationListener(listener ApplicationListener) {
	broadcaster.mu.Lock()
	for _, eventId := range listener.SubscribeEvents() {
		for _, eventListener := range broadcaster.eventListenerMap[eventId] {
			if eventListener.GetApplicationListenerName() == listener.GetApplicationListenerName() {
				broadcaster.mu.Unlock()
				return
			}
		}
		broadcaster.eventListenerMap[eventId] = append(broadcaster.eventListenerMap[eventId], listener)
	}
	broadcaster.mu.Unlock()
}

func (broadcaster *SimpleApplicationEventBroadcaster) UnregisterApplicationListener(listener ApplicationListener) {
	broadcaster.mu.Lock()
	for _, eventId := range listener.SubscribeEvents() {
		for registeredEventId, events := range broadcaster.eventListenerMap {
			if eventId == registeredEventId {
				broadcaster.deleteEventListener(registeredEventId, listener, events)
			}
		}
	}
	broadcaster.mu.Unlock()
}

func (broadcaster *SimpleApplicationEventBroadcaster) deleteEventListener(eventId ApplicationEventId, listener ApplicationListener, eventListeners []ApplicationListener) {
	tempListeners := eventListeners
	for index, eventListener := range tempListeners {
		if eventListener.GetApplicationListenerName() == listener.GetApplicationListenerName() {
			tempListeners = append(tempListeners[:index], tempListeners[index+1:]...)
		}
	}
	if len(tempListeners) == 0 {
		delete(broadcaster.eventListenerMap, eventId)
	} else {
		broadcaster.eventListenerMap[eventId] = tempListeners
	}
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
		listener.OnApplicationEvent(context, event)
	}
}
