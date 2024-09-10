package pkg

import (
	"github.com/KScaesar/go-layout/pkg/utility"
)

func init() {
	defaultVersion.Store(newVersion("CRM"))
	SetDefaultLogger(utility.LoggerWhenGoTest(false))
	SetDefaultShutdown(NewShutdown(DefaultLogger().Logger, 0))
}
