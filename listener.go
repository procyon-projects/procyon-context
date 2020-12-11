package context

type ApplicationListener interface {
	GetApplicationListenerName() string
	SubscribeEvents() []ApplicationEventId
	OnApplicationEvent(context Context, event ApplicationEvent)
}
