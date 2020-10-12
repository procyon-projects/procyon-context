package context

import "errors"

type ConfigurationPropertiesBindingProcessor struct {
	context ApplicationContext
	binder  ConfigurationPropertiesBinder
}

func NewConfigurationPropertiesBindingProcessor() *ConfigurationPropertiesBindingProcessor {
	return &ConfigurationPropertiesBindingProcessor{}
}

func (processor *ConfigurationPropertiesBindingProcessor) SetApplicationContext(context ApplicationContext) {
	processor.context = context
}

func (processor *ConfigurationPropertiesBindingProcessor) AfterProperties() {
	processor.binder = NewConfigurationPropertiesBinder(processor.context.(ConfigurableApplicationContext))
}

func (processor *ConfigurationPropertiesBindingProcessor) BeforeInitialization(peaName string, pea interface{}) (interface{}, error) {
	err := processor.binder.Bind(pea)
	if err != nil {
		return nil, errors.New("error occurred while configuration properties was being binding to pea instance : " + peaName)
	}
	return pea, err
}

func (processor *ConfigurationPropertiesBindingProcessor) Initialize() error {
	return nil
}

func (processor *ConfigurationPropertiesBindingProcessor) AfterInitialization(peaName string, pea interface{}) (interface{}, error) {
	return pea, nil
}
