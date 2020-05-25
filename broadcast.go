package context

type ApplicationEventBroadcaster interface {
	RegisterApplicationListener(listener ApplicationListener)
	RegisterApplicationListenerByPeaName(peaName string)
	UnregisterApplicationListener(listener ApplicationListener)
	UnregisterApplicationListenerByPeaName(peaName string)
	RemoveAllApplicationListeners()
	BroadcastEvent(event ApplicationEvent)
}

type BaseApplicationEventBroadcaster struct {
}

func NewBaseApplicationEventBroadcaster() *BaseApplicationEventBroadcaster {
	return &BaseApplicationEventBroadcaster{}
}

func (broadcaster *BaseApplicationEventBroadcaster) RegisterApplicationListener(listener ApplicationListener) {

}

func (broadcaster *BaseApplicationEventBroadcaster) RegisterApplicationListenerByPeaName(peaName string) {

}

func (broadcaster *BaseApplicationEventBroadcaster) UnregisterApplicationListener(listener ApplicationListener) {

}

func (broadcaster *BaseApplicationEventBroadcaster) UnregisterApplicationListenerByPeaName(peaName string) {

}

func (broadcaster *BaseApplicationEventBroadcaster) RemoveAllApplicationListeners() {

}

func (broadcaster *BaseApplicationEventBroadcaster) BroadcastEvent(event ApplicationEvent) {
	// do nothing
}
