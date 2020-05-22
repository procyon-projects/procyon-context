package context

import "github.com/codnect/go-one"

type Configuration interface {
	Register() []one.Func
}

type ActivateConfigurationProperties interface {
	ConfigurationProperties() []ConfigurationProperties
}
