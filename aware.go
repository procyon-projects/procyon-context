package context

type ApplicationContextAware interface {
	SetApplicationContext(context ApplicationContext)
}

type ApplicationEventPublisherAware interface {
	SetApplicationEventPublisher(publisher ApplicationEventPublisher)
}
