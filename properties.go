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
	context              ConfigurableApplicationContext
	env                  core.Environment
	typeConverterService core.TypeConverterService
}

func NewConfigurationPropertiesBinder(context ConfigurableApplicationContext) ConfigurationPropertiesBinder {
	return ConfigurationPropertiesBinder{
		context,
		context.GetEnvironment(),
		context.GetEnvironment().GetTypeConverterService(),
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
		fieldType := &core.Type{Val: field}
		if jsonTagValue != "" {
			binder.bindTargetField(fieldType, binder.getFullPropertyName(prefix, jsonTagValue), defaultValue)
		} else if yamlTagValue != "" {
			binder.bindTargetField(fieldType, binder.getFullPropertyName(prefix, yamlTagValue), defaultValue)
		}
	}
	return nil
}

func (binder ConfigurationPropertiesBinder) bindTargetField(fieldType *core.Type, propertyName string, defaultValue string) {
	propertyValue := binder.env.GetProperty(propertyName, defaultValue)
	if propertyValue != nil {
		if fieldType.Val.IsValid() && fieldType.Val.CanSet() {
			if binder.typeConverterService.CanConvert(core.GetType(propertyValue), fieldType) {
				value := binder.typeConverterService.Convert(propertyValue, core.GetType(propertyValue), fieldType)
				fieldType.Val.Set(reflect.ValueOf(value))
			}
		}
	}
}

func (binder ConfigurationPropertiesBinder) getFullPropertyName(prefix string, tagValue string) string {
	return prefix + "." + tagValue
}
