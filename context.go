package context

import (
	"github.com/Rollcomp/procyon-core"
	web "github.com/Rollcomp/procyon-web"
)

type ApplicationContext interface {
	GetApplicationName() string
	GetStartupTimeStamp() int64
}

type WebApplicationContext interface {
	ApplicationContext
}

type ConfigurableApplicationContext interface {
	SetEnvironment(environment core.ConfigurableEnvironment)
	GetEnvironment() core.ConfigurableEnvironment
}

type WebServerApplicationContext struct {
	webServer web.Server
}

func NewWebServerApplicationContext() *WebServerApplicationContext {
	return &WebServerApplicationContext{}
}

func (ctx *WebServerApplicationContext) createWebServer() {
	if ctx.webServer == nil {
		ctx.webServer, _ = web.GetWebServer()
	}
}

type DefaultWebServerApplicationContext struct {
}
