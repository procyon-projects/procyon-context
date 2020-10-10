package context

type ApplicationListener interface {
	SubscribeEvents(context ApplicationContext) []ApplicationEvent
	OnApplicationEvent(context ApplicationContext, event ApplicationEvent)
}
