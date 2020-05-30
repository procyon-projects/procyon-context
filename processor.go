package context

type ConfigurationPropertiesBindingProcessor struct {
	context ApplicationContext
}

func NewConfigurationPropertiesBindingProcessor() ConfigurationPropertiesBindingProcessor {
	return ConfigurationPropertiesBindingProcessor{}
}

func (processor ConfigurationPropertiesBindingProcessor) SetApplicationContext(context ApplicationContext) {
	processor.context = context
}

func (processor ConfigurationPropertiesBindingProcessor) BeforeInitialization(peaName string, pea interface{}) (interface{}, error) {
	panic("implement me")
}

func (processor ConfigurationPropertiesBindingProcessor) Initialize() error {
	panic("implement me")
}

func (processor ConfigurationPropertiesBindingProcessor) AfterInitialization(peaName string, pea interface{}) (interface{}, error) {
	panic("implement me")
}
