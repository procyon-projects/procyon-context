package context

import "time"

type ApplicationEvent interface {
	GetSource() interface{}
	GetTimestamp() int64
}

type BaseApplicationEvent struct {
	source    interface{}
	timestamp int64
}

func NewBaseApplicationEvent(source interface{}) BaseApplicationEvent {
	return BaseApplicationEvent{
		source:    source,
		timestamp: time.Now().Unix(),
	}
}

func (appEvent BaseApplicationEvent) GetSource() interface{} {
	return appEvent.source
}

func (appEvent BaseApplicationEvent) GetTimestamp() int64 {
	return appEvent.timestamp
}

type ApplicationContextEvent interface {
	ApplicationEvent
	GetApplicationContext() ApplicationContext
}

type BaseApplicationContextEvent struct {
	BaseApplicationEvent
}

func NewBaseApplicationContextEvent(source interface{}) BaseApplicationContextEvent {
	return BaseApplicationContextEvent{
		NewBaseApplicationEvent(source),
	}
}

func (appContextEvent BaseApplicationContextEvent) GetApplicationContext() ApplicationContext {
	return appContextEvent.source.(ApplicationContext)
}

type ApplicationContextStartedEvent struct {
	BaseApplicationContextEvent
}

func NewApplicationContextStartedEvent(source ApplicationContext) ApplicationContextStartedEvent {
	return ApplicationContextStartedEvent{
		NewBaseApplicationContextEvent(source),
	}
}

type ApplicationContextStoppedEvent struct {
	BaseApplicationContextEvent
}

func NewApplicationContextStoppedEvent(source ApplicationContext) ApplicationContextStoppedEvent {
	return ApplicationContextStoppedEvent{
		NewBaseApplicationContextEvent(source),
	}
}

type ApplicationContextRefreshedEvent struct {
	BaseApplicationContextEvent
}

func NewApplicationContextRefreshedEvent(source ApplicationContext) ApplicationContextRefreshedEvent {
	return ApplicationContextRefreshedEvent{
		NewBaseApplicationContextEvent(source),
	}
}

type ApplicationContextClosedEvent struct {
	BaseApplicationContextEvent
}

func NewApplicationContextClosedEvent(source ApplicationContext) ApplicationContextClosedEvent {
	return ApplicationContextClosedEvent{
		NewBaseApplicationContextEvent(source),
	}
}
