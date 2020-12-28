package context

import core "github.com/procyon-projects/procyon-core"

func init() {
	core.Register(newLoggingProperties)
	/* Configuration Properties Binding Processor */
	core.Register(NewConfigurationPropertiesBindingProcessor)
}
