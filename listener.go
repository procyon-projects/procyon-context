package context

type ApplicationListener interface {
	SubscribeEvents() []ApplicationEventId
	OnApplicationEvent(context Context, event ApplicationEvent)
}
