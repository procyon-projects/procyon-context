package context

import (
	"errors"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
)

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

type ConfigurationPropertiesBindingProcessor struct {
	binder ConfigurationPropertiesBinder
}

func NewConfigurationPropertiesBindingProcessor(env core.Environment, typeConverterService core.TypeConverterService) ConfigurationPropertiesBindingProcessor {
	return ConfigurationPropertiesBindingProcessor{
		newConfigurationPropertiesBinder(env, typeConverterService),
	}
}

func (processor ConfigurationPropertiesBindingProcessor) BeforePeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	err := processor.binder.Bind(pea)
	if err != nil {
		return nil, errors.New("error occurred while configuration properties was being binding to pea instance : " + peaName)
	}
	return pea, nil
}

func (processor ConfigurationPropertiesBindingProcessor) InitializePea() error {
	return nil
}

func (processor ConfigurationPropertiesBindingProcessor) AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	return pea, nil
}
