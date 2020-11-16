package context

import (
	"errors"
	core "github.com/procyon-projects/procyon-core"
)

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
