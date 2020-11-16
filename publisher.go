package context

type ApplicationEventPublisher interface {
	PublishEvent(context Context, event ApplicationEvent)
}
