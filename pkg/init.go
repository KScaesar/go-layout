package pkg

import (
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func init() {
	SetLogger(wlog.NewLoggerWhenNormalRun(false))
	SetShutdown(NewShutdownWhenDefault(Logger().Logger, 0))
}

// Init initializes the necessary default global variables
func Init(conf *configs.Config) {
	logger := wlog.NewLogger(os.Stdout, &conf.Logger, wlog.DefaultFormat()...)
	logger.Logger = logger.With(slog.String("svc", Version().ServiceName))
	logger.SetStdDefaultLevel()
	logger.SetStdDefaultLogger()
	SetLogger(logger)
	SetShutdown(NewShutdownWhenDefault(logger.Logger, 0))
	return
}
