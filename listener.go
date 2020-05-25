package context

type ApplicationListener interface {
	SubscribeEvents() []ApplicationEvent
	OnApplicationEvent(event ApplicationEvent)
}
