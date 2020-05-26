package context

import (
	core "github.com/Rollcomp/procyon-core"
	peas "github.com/Rollcomp/procyon-peas"
)

const defaultEventListenerRetrieverSize = 16

type ApplicationEventListenerRetriever interface {
	GetApplicationListeners() []ApplicationListener
	AddApplicationListener(listener ApplicationListener)
	AddApplicationListenerByPeaName(peaName string)
	RemoveApplicationListener(listener ApplicationListener)
	RemoveApplicationListenerByPeaName(peaName string)
	RemoveAllApplicationListeners()
}

type DefaultApplicationEventListenerRetriever struct {
	appEventListeners    map[string]ApplicationListener
	appEventListenerPeas map[string]interface{}
	peaFactory           peas.ConfigurablePeaFactory
}

func NewDefaultApplicationEventListenerRetriever(factory peas.ConfigurablePeaFactory) DefaultApplicationEventListenerRetriever {
	if factory == nil {
		panic("Pea Factory must not be null")
	}
	return DefaultApplicationEventListenerRetriever{
		appEventListeners:    make(map[string]ApplicationListener, defaultEventListenerRetrieverSize),
		appEventListenerPeas: make(map[string]interface{}, defaultEventListenerRetrieverSize),
		peaFactory:           factory,
	}
}

func (retriever DefaultApplicationEventListenerRetriever) GetApplicationListeners() []ApplicationListener {
	listeners := make([]ApplicationListener, defaultEventListenerRetrieverSize)
	for key := range retriever.appEventListeners {
		listeners = append(listeners, retriever.appEventListeners[key])
	}
	for peaName := range retriever.appEventListenerPeas {
		if peaName != "" {
			peaObj, err := retriever.peaFactory.GetPeaByNameAndType(peaName, core.GetType((*ApplicationListener)(nil)))
			if err != nil {
				listeners = append(listeners, peaObj.(ApplicationListener))
			}
		}
	}
	return listeners
}

func (retriever DefaultApplicationEventListenerRetriever) AddApplicationListener(listener ApplicationListener) {
	typ := core.GetType(listener)
	if core.IsStruct(typ) {
		retriever.appEventListeners[typ.String()] = listener
	} else {
		panic("It must be struct")
	}
}

func (retriever DefaultApplicationEventListenerRetriever) AddApplicationListenerByPeaName(peaName string) {
	retriever.appEventListeners[peaName] = nil
}

func (retriever DefaultApplicationEventListenerRetriever) RemoveApplicationListener(listener ApplicationListener) {
	typ := core.GetType(listener)
	if core.IsStruct(typ) {
		_, ok := retriever.appEventListeners[typ.String()]
		if ok {
			delete(retriever.appEventListeners, typ.String())
		}
	} else {
		panic("It must be struct")
	}
}

func (retriever DefaultApplicationEventListenerRetriever) RemoveApplicationListenerByPeaName(peaName string) {
	_, ok := retriever.appEventListenerPeas[peaName]
	if ok {
		delete(retriever.appEventListenerPeas, peaName)
	}
}

func (retriever DefaultApplicationEventListenerRetriever) RemoveAllApplicationListeners() {
	retriever.appEventListeners = make(map[string]ApplicationListener, defaultEventListenerRetrieverSize)
	retriever.appEventListenerPeas = make(map[string]interface{}, defaultEventListenerRetrieverSize)
}
