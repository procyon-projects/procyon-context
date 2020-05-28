package context

type ApplicationEventPublisher interface {
	PublishEvent(event ApplicationEvent)
}
