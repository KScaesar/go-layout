package pkg

import (
	"github.com/KScaesar/go-layout/pkg/utility"
)

func init() {
	defaultVersion.Store(newVersion("CRM"))
	defaultLogger.Store(utility.LoggerWhenGoTest())
	defaultShutdown.Store(NewShutdown(DefaultLogger().Logger))
}
