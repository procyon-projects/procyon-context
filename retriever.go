package context

type ApplicationEventListenerRetriever interface {
	GetListeners() []ApplicationListener
}

type DefaultApplicationEventListenerRetriever struct {
}

func (retriever DefaultApplicationEventListenerRetriever) GetListeners() []ApplicationListener {
	return nil
}
