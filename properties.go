package context

import (
	"errors"
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
)

type ConfigurationProperties interface {
	GetPrefix() string
}

type ConfigurationPropertiesBinder struct {
	env                  core.Environment
	typeConverterService core.TypeConverterService
}

func NewConfigurationPropertiesBinder(env core.Environment, typeConverterService core.TypeConverterService) ConfigurationPropertiesBinder {
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
			if value != "" {

			} else if err == nil && defaultTag.Value != "" {

			}
		}
	}
	return nil
}

func (binder ConfigurationPropertiesBinder) bindTargetField(fieldType goo.Type, propertyName string, defaultValue string) error {
	/*propertyValue := binder.env.GetProperty(propertyName, defaultValue)
	if propertyValue != nil {
		if fieldType.Val.IsValid() && fieldType.Val.CanSet() {
			if binder.typeConverterService.CanConvert(goo.GetType(propertyValue), fieldType) {
				value, err := binder.typeConverterService.Convert(propertyValue, goo.GetType(propertyValue), fieldType)
				if err != nil {
					return err
				}
				fieldType.Val.Set(reflect.ValueOf(value))
			}
		}
	}*/
	return nil
}

func (binder ConfigurationPropertiesBinder) getFullPropertyName(prefix string, tagValue string) string {
	return prefix + "." + tagValue
}
