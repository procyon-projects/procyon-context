package context

import (
	"errors"
	core "github.com/procyon-projects/procyon-core"
	"reflect"
)

type ConfigurationProperties interface {
	GetPrefix() string
}

type ConfigurationPropertiesBinder struct {
	context         ConfigurableApplicationContext
	propertySources *core.PropertySources
}

func NewConfigurationPropertiesBinder(context ConfigurableApplicationContext) ConfigurationPropertiesBinder {
	return ConfigurationPropertiesBinder{
		context:         context,
		propertySources: context.GetEnvironment().GetPropertySources(),
	}
}

func (binder ConfigurationPropertiesBinder) Bind(target interface{}) error {
	if target == nil {
		return nil
	}
	if properties, ok := target.(ConfigurationProperties); ok {
		prefix := properties.GetPrefix()
		if prefix == "" {
			return errors.New("prefix must not be null")
		}
		if !core.IsPtr(target) {
			return errors.New("this object cannot be bound the configuration properties")
		}
		return binder.bindTargetFields(prefix, target)
	}
	return errors.New("it must implement ConfigurationProperties")
}

func (binder ConfigurationPropertiesBinder) bindTargetFields(prefix string, target interface{}) error {
	targetTyp := core.GetType(target)
	numField := core.GetNumField(targetTyp)
	for index := 0; index < numField; index++ {
		structField := core.GetStructFieldByIndex(targetTyp, index)
		defaultValue := structField.Tag.Get("default")
		jsonTagValue := structField.Tag.Get("json")
		yamlTagValue := structField.Tag.Get("yaml")
		field := core.GetFieldValueByIndex(targetTyp, index)
		if jsonTagValue != "" {
			binder.bindTargetField(field, binder.getFullPropertyName(prefix, jsonTagValue), defaultValue)
		} else if yamlTagValue != "" {
			binder.bindTargetField(field, binder.getFullPropertyName(prefix, yamlTagValue), defaultValue)
		}
	}
	return nil
}

func (binder ConfigurationPropertiesBinder) bindTargetField(field reflect.Value, propertyName string, defaultValue string) {
	// conversion for value will be added
	if field.IsValid() && field.CanSet() {

	}
}

func (binder ConfigurationPropertiesBinder) getFullPropertyName(prefix string, tagValue string) string {
	return prefix + "." + tagValue
}
