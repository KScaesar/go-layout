package pkg

import (
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func init() {
	defaultVersion.Store(newVersion("CRM"))
	SetDefaultLogger(wlog.LoggerWhenGoTest(false))
	SetDefaultShutdown(NewShutdownWhenInit(DefaultLogger().Logger, 0))
}

// Init initializes the necessary default global variables
func Init(conf *configs.Config) {
	writer := os.Stdout
	logger := wlog.NewLogger(writer, &conf.Logger)
	logger.Logger = logger.With(slog.String("svc", DefaultVersion().ServiceName))
	logger.SetStdDefaultLevel()
	logger.SetStdDefaultLogger()
	SetDefaultLogger(logger)
	SetDefaultShutdown(NewShutdownWhenInit(logger.Logger, 0))
	return
}
