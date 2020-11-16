package context

import (
	"errors"
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
)

type ConfigurationPropertiesBinder struct {
	env                  core.Environment
	typeConverterService core.TypeConverterService
}

func newConfigurationPropertiesBinder(env core.Environment, typeConverterService core.TypeConverterService) ConfigurationPropertiesBinder {
	return ConfigurationPropertiesBinder{
		env,
		typeConverterService,
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
		if !goo.GetType(target).IsPointer() {
			return errors.New("configuration properties cannot be bound as it is not a pointer of type")
		}
		return binder.bindTargetFields(prefix, target)
	}
	return nil
}

func (binder ConfigurationPropertiesBinder) bindTargetFields(prefix string, target interface{}) error {
	targetTyp := goo.GetType(target).ToStructType()
	exportedFields := targetTyp.GetExportedFields()
	for _, field := range exportedFields {
		var bindTag, defaultTag goo.Tag
		var err error
		bindTag, err = field.GetTagByName("json")
		if err != nil {
			bindTag, err = field.GetTagByName("yaml")
		}
		if err == nil {
			value := bindTag.Value
			defaultTag, err = field.GetTagByName("default")
			if err != nil && value != "" {
				binder.bindTargetField(field, target, binder.getFullPropertyName(prefix, value), "")
			} else if err == nil && defaultTag.Value != "" {
				binder.bindTargetField(field, target, binder.getFullPropertyName(prefix, value), defaultTag.Value)
			}
		}
	}
	return nil
}

func (binder ConfigurationPropertiesBinder) bindTargetField(field goo.Field, instance interface{}, propertyName string, defaultValue string) error {
	propertyValue := binder.env.GetProperty(propertyName, defaultValue)
	if propertyValue != nil {
		if field.CanSet() {
			propertyValueType := goo.GetType(propertyValue)
			fieldType := field.GetType()
			if propertyValueType.GetGoType() == fieldType.GetGoType() {
				field.SetValue(instance, propertyValue)
			} else if binder.typeConverterService.CanConvert(propertyValueType, fieldType) {
				value, err := binder.typeConverterService.Convert(propertyValue, propertyValueType, fieldType)
				if err != nil {
					return err
				}
				field.SetValue(instance, value)
			}
		}
	}
	return nil
}

func (binder ConfigurationPropertiesBinder) getFullPropertyName(prefix string, tagValue string) string {
	return prefix + "." + tagValue
}
