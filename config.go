package context

import "github.com/codnect/go-one"

type Configuration interface {
	Register() []one.Func
}

type ActivateConfigurationProperties interface {
	ConfigurationProperties() []one.Func
}

type ImportConfiguration interface {
	Imports() []one.Func
}

type ConfigurationCollection interface {
	Configurations() []one.Func
}
