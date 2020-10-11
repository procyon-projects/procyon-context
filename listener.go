package context

type ApplicationListener interface {
	SubscribeEvents() []ApplicationEvent
	OnApplicationEvent(context ApplicationContext, event ApplicationEvent)
}
