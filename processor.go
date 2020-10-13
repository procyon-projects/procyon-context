package context

import peas "github.com/procyon-projects/procyon-peas"

type BootstrapProcessor struct {
}

func NewBootstrapProcessor() BootstrapProcessor {
	return BootstrapProcessor{}
}

func (processor BootstrapProcessor) AfterPeaDefinitionRegistryInitialization(registry peas.PeaDefinitionRegistry) {
	processor.processPeaDefinitions(registry)
}

func (processor BootstrapProcessor) AfterPeaFactoryInitialization(factory peas.ConfigurablePeaFactory) {
	// do something
}

func (processor BootstrapProcessor) processPeaDefinitions(registry peas.PeaDefinitionRegistry) {
	scanner := NewComponentPeaDefinitionScanner(registry)
	scanner.DoScan()
}

type EventListenerProcessor struct {
}

func NewEventListenerProcessor() EventListenerProcessor {
	return EventListenerProcessor{}
}

func (processor EventListenerProcessor) AfterPeaDefinitionRegistryInitialization(registry peas.PeaDefinitionRegistry) {
	// do something
}

func (processor EventListenerProcessor) AfterPeaFactoryInitialization(factory peas.ConfigurablePeaFactory) {
	// do something
}
