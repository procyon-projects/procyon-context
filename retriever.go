package context

import (
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
)

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

func NewDefaultApplicationEventListenerRetriever() DefaultApplicationEventListenerRetriever {
	return DefaultApplicationEventListenerRetriever{
		appEventListeners:    make(map[string]ApplicationListener, 0),
		appEventListenerPeas: make(map[string]interface{}, 0),
	}
}

func NewDefaultApplicationEventListenerRetrieverWithFactory(factory peas.ConfigurablePeaFactory) DefaultApplicationEventListenerRetriever {
	if factory == nil {
		core.Log.Fatal("Pea Factory must not be null")
	}
	return DefaultApplicationEventListenerRetriever{
		appEventListeners:    make(map[string]ApplicationListener, 0),
		appEventListenerPeas: make(map[string]interface{}, 0),
		peaFactory:           factory,
	}
}

func (retriever DefaultApplicationEventListenerRetriever) GetApplicationListeners() []ApplicationListener {
	listeners := make([]ApplicationListener, 0)
	for key := range retriever.appEventListeners {
		listeners = append(listeners, retriever.appEventListeners[key])
	}
	if retriever.peaFactory != nil {
		for peaName := range retriever.appEventListenerPeas {
			if peaName != "" {
				peaObj, err := retriever.peaFactory.GetPeaByNameAndType(peaName, core.GetType((*ApplicationListener)(nil)))
				if err != nil {
					listeners = append(listeners, peaObj.(ApplicationListener))
				}
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
		core.Log.Error("It must be struct")
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
		core.Log.Error("It must be struct")
	}
}

func (retriever DefaultApplicationEventListenerRetriever) RemoveApplicationListenerByPeaName(peaName string) {
	_, ok := retriever.appEventListenerPeas[peaName]
	if ok {
		delete(retriever.appEventListenerPeas, peaName)
	}
}

func (retriever DefaultApplicationEventListenerRetriever) RemoveAllApplicationListeners() {
	retriever.appEventListeners = make(map[string]ApplicationListener, 0)
	retriever.appEventListenerPeas = make(map[string]interface{}, 0)
}
