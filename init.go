package context

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Initialize Pools */
	initBaseApplicationContextPool()
	/* Configuration Properties Binding Processor */
	core.Register(NewConfigurationPropertiesBindingProcessor)
}
