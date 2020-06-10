package context

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Configuration Properties Binding Processor */
	core.Register(NewConfigurationPropertiesBindingProcessor)
}
