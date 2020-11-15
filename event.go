package context

import (
	"time"
)

type ApplicationEventId uint64

var applicationContextEventId = GetEventId("github.com.procyon.ApplicationContextEvent")
var applicationContextStartedEventId = GetEventId("github.com.procyon.ApplicationContextStartedEvent")
var applicationContextStoppedEventId = GetEventId("github.com.procyon.ApplicationContextStoppedEvent")
var applicationContextRefreshedEventId = GetEventId("github.com.procyon.ApplicationContextRefreshedEvent")
var applicationContextClosedEventId = GetEventId("github.com.procyon.ApplicationContextClosedEvent")

func ApplicationContextEventId() ApplicationEventId {
	return applicationContextEventId
}

func ApplicationContextStartedEventId() ApplicationEventId {
	return applicationContextStartedEventId
}

func ApplicationContextStoppedEventId() ApplicationEventId {
	return applicationContextStoppedEventId
}

func ApplicationContextRefreshedEventId() ApplicationEventId {
	return applicationContextRefreshedEventId
}

func ApplicationContextClosedEventId() ApplicationEventId {
	return applicationContextClosedEventId
}

func GetEventId(eventName string) ApplicationEventId {
	if len(eventName) > 0 {
		hash := uint64(0)
		for _, character := range eventName {
			hash = uint64(31)*hash + uint64(character)
		}
		return ApplicationEventId(hash)
	}
	return 0
}

type ApplicationEvent interface {
	GetEventId() ApplicationEventId
	GetParentEventId() ApplicationEventId
	GetSource() interface{}
	GetTimestamp() int64
}

type ApplicationContextEvent interface {
	ApplicationEvent
	GetApplicationContext() ApplicationContext
}

type ApplicationContextStartedEvent struct {
	source    ApplicationContext
	timestamp int64
}

func NewApplicationContextStartedEvent(source ApplicationContext) ApplicationContextStartedEvent {
	return ApplicationContextStartedEvent{
		source:    source,
		timestamp: time.Now().Unix(),
	}
}

func (event ApplicationContextStartedEvent) GetEventId() ApplicationEventId {
	return applicationContextStartedEventId
}

func (event ApplicationContextStartedEvent) GetParentEventId() ApplicationEventId {
	return applicationContextEventId
}

func (event ApplicationContextStartedEvent) GetSource() interface{} {
	return event.source
}

func (event ApplicationContextStartedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationContextStartedEvent) GetApplicationContext() ApplicationContext {
	return event.source
}

type ApplicationContextStoppedEvent struct {
	source    ApplicationContext
	timestamp int64
}

func NewApplicationContextStoppedEvent(source ApplicationContext) ApplicationContextStoppedEvent {
	return ApplicationContextStoppedEvent{
		source:    source,
		timestamp: time.Now().Unix(),
	}
}

func (event ApplicationContextStoppedEvent) GetEventId() ApplicationEventId {
	return applicationContextStoppedEventId
}

func (event ApplicationContextStoppedEvent) GetParentEventId() ApplicationEventId {
	return applicationContextEventId
}

func (event ApplicationContextStoppedEvent) GetSource() interface{} {
	return event.source
}

func (event ApplicationContextStoppedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationContextStoppedEvent) GetApplicationContext() ApplicationContext {
	return event.source
}

type ApplicationContextRefreshedEvent struct {
	source    ApplicationContext
	timestamp int64
}

func NewApplicationContextRefreshedEvent(source ApplicationContext) ApplicationContextRefreshedEvent {
	return ApplicationContextRefreshedEvent{
		source:    source,
		timestamp: time.Now().Unix(),
	}
}

func (event ApplicationContextRefreshedEvent) GetEventId() ApplicationEventId {
	return applicationContextRefreshedEventId
}

func (event ApplicationContextRefreshedEvent) GetParentEventId() ApplicationEventId {
	return applicationContextEventId
}

func (event ApplicationContextRefreshedEvent) GetSource() interface{} {
	return event.source
}

func (event ApplicationContextRefreshedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationContextRefreshedEvent) GetApplicationContext() ApplicationContext {
	return event.source
}

type ApplicationContextClosedEvent struct {
	source    ApplicationContext
	timestamp int64
}

func NewApplicationContextClosedEvent(source ApplicationContext) ApplicationContextClosedEvent {
	return ApplicationContextClosedEvent{
		source:    source,
		timestamp: time.Now().Unix(),
	}
}

func (event ApplicationContextClosedEvent) GetEventId() ApplicationEventId {
	return applicationContextClosedEventId
}

func (event ApplicationContextClosedEvent) GetParentEventId() ApplicationEventId {
	return applicationContextEventId
}

func (event ApplicationContextClosedEvent) GetSource() interface{} {
	return event.source
}

func (event ApplicationContextClosedEvent) GetTimestamp() int64 {
	return event.timestamp
}

func (event ApplicationContextClosedEvent) GetApplicationContext() ApplicationContext {
	return event.source
}
