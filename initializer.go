package context

type ApplicationContextInitializer interface {
	InitializeContext(context ConfigurableApplicationContext)
}
