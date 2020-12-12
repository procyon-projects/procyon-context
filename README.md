<img src="https://procyon-projects.github.io/img/logo.png" width="128">

# Procyon Context
[![Go Report Card](https://goreportcard.com/badge/github.com/procyon-projects/procyon-context)](https://goreportcard.com/report/github.com/procyon-projects/procyon-context)
[![codecov](https://codecov.io/gh/procyon-projects/procyon-context/branch/master/graph/badge.svg?token=8Q2DVS1SZX)](https://codecov.io/gh/procyon-projects/procyon-context)
[![Build Status](https://travis-ci.com/procyon-projects/procyon-context.svg?branch=master)](https://travis-ci.com/procyon-projects/procyon-context)
[![Gitter](https://badges.gitter.im/procyon-projects/community.svg)](https://gitter.im/procyon-projects/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![PkgGoDev](https://pkg.go.dev/badge/procyon-projects/procyon)](https://pkg.go.dev/github.com/procyon-projects/procyon-context)

This gives you a basic understanding of Procyon Context Module. It covers
components provided by the framework, such as Logger, Properties, Initializers and Events.

## Logger
It provides an interface for loggers. A logger is provided by the framework.
You can get it in construction function by using the dependency injection.
If you want to implement your logger, you need to implement the interface.

```go
type Logger interface {
	Trace(ctx interface{}, message interface{})
	Debug(ctx interface{}, message interface{})
	Info(ctx interface{}, message interface{})
	Warning(ctx interface{}, message interface{})
	Error(ctx interface{}, message interface{})
	Fatal(ctx interface{}, message interface{})
	Panic(ctx interface{}, message interface{})
	Tracef(ctx interface{}, format string, args ...interface{})
	Debugf(ctx interface{}, format string, args ...interface{})
	Infof(ctx interface{}, format string, args ...interface{})
	Warningf(ctx interface{}, format string, args ...interface{})
	Errorf(ctx interface{}, format string, args ...interface{})
	Fatalf(ctx interface{}, format string, args ...interface{})
	Panicf(ctx interface{}, format string, args ...interface{})
}
```

## Configuration Properties
This interface is used to bind the command-line parameters to your struct's instance.
```go
type ConfigurationProperties interface {
	GetPrefix() string
}
```
The example is given below. Note that you need to register your properties by using the function **core.Register**.
Otherwise, their instances won't be created by the framework.
```go
type MyConfigurationProperties struct {
	    Name                string  `yaml:"name" json:"name" default:"Test Application"`
        CustomPort          uint    `yaml:"port" json:"port" default:"8090"`
}

func NewMyConfigurationProperties() *MyConfigurationProperties() {
    return &MyConfigurationProperties{}
}

func (properties *MyConfigurationProperties) GetPrefix() string  {
    return "application"
}

```
When you specify the parameters **--application.name** and **--application.port**, they will be bind to 
your instance. Otherwise, their default values will be used.

## Application Context Initializer
This interface is used to initialize the context by custom context initializer. It is invoked 
while the context is prepared. 
```go
type ApplicationContextInitializer interface {
	InitializeContext(context ConfigurableApplicationContext)
}
```
The example of a custom context initializer is given below. Note that you need to register your initializer by using the function core.Register. 
Otherwise, their instances won't be created by the framework.
```go
type CustomContextInitializer struct {

}

func NewCustomContextInitializer() CustomContextInitializer {
    return CustomContextInitializer{}
}

func (initializer CustomContextInitializer) InitializeContext(context ConfigurableApplicationContext) {
    // do whatever you want
}
```

## Application Event
All events have to implement this interface. You can have custom events by implementing
the interface. 
```go
type ApplicationEvent interface {
	GetEventId() ApplicationEventId
	GetParentEventId() ApplicationEventId
	GetSource() interface{}
	GetTimestamp() int64
}
```

* **GetEventId** returns the unique id of the event.
* **GetParentEventId** returns the unique id of the parent event of the event. Parent Event is not necessary.
* **GetSource**  returns the object with which the event is associated.
* **GetTimestamp** returns the system time in milliseconds when the event occurred.

The example is given below.

First, Your event has to have an unique id, you can get it by using the function **context.GetEventId**.

```go
var customEventId = context.GetEventId("github.com.procyon.CustomEvent")

func CustomEventId() ApplicationEventId {
	return customEventId
}
```

The second thing you need to do is to implement the interface **context.ApplicationEvent**
It's shown below.

```go
type CustomEvent struct {
	source    interface{}
	timestamp int64
}

func NewCustomEvent(obj interface) CustomEvent {
	return CustomEvent{
		source:    obj,
		timestamp: time.Now().Unix(),
	}
}

func (event CustomEvent) GetEventId() context.ApplicationEventId {
	return customEventId
}

func (event CustomEvent) GetParentEventId() context.ApplicationEventId {
	return -1
}

func (event CustomEvent) GetSource() interface{} {
	return event.source
}

func (event CustomEvent) GetTimestamp() int64 {
	return event.timestamp
}
```

## Application Listener
This interface need to be implemented by application event listeners.
Event Ids to be subscribed need to be returned by SubscribeEvents.
All events must have an unique id. Otherwise, there will be conflicts. 
That's why when you want to create a custom event, use the function **context.GetEventId**
to have an event id. 
```go
type ApplicationListener interface {
	SubscribeEvents() []ApplicationEventId
	OnApplicationEvent(context Context, event ApplicationEvent)
}
```

The example is given below. Note that you need to register your listeners by using the function **core.Register**.
Otherwise, their instances won't be created by the framework.
```go
type CustomEventListener interface {

}
	
func NewCustomEventListener() CustomEventListener {
    return CustomEventListener{}
}	

func (listener CustomEventListener) SubscribeEvents() []context.ApplicationEventId {
    return []context.ApplicationEventId {
        CustomEventId(),
    }
}

func (listener CustomEventListener) OnApplicationEvent(context context.Context, event context.ApplicationEvent) {
    // do whatever you want...
}
```

## Application Event Publisher
It is used to notify all matching listeners registered. Events might be framework events
or application-specific events. A framework event publisher is provided by the framework.
You can get it in construction function by using the dependency injection. It's recommend
to use for async execution. Also, you can have your custom event publisher by implementing
this interface.

```go
type ApplicationEventPublisher interface {
	PublishEvent(context Context, event ApplicationEvent)
}
```

The example of a custom event publisher is given below. Note that you need to register your publisher by using the function **core.Register**.
Otherwise, their instances won't be created by the framework.

```go
type CustomEventPublisher struct {
	
}

func NewCustomEventPublisher() CustomEventPublisher  {
	return CustomEventPublisher{}
}

func (publisher CustomEventPublisher) PublishEvent(context context.Context, event context.ApplicationEvent) {
    // do whatever you want...
}
```

## License
Procyon Framework is released under version 2.0 of the Apache License
