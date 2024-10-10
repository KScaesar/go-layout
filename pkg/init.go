package pkg

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync/atomic"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func init() {
	defaultShutdown.Store(utility.NewShutdown(context.Background(), -1, Logger().Logger))
	go Shutdown().Serve()
}

// Init initializes the necessary default global variables
func Init(conf *Config) io.Closer {
	var err error
	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	logWriter, err := initLogger(conf.Filename.Logger, &conf.Logger)
	if err != nil {
		Logger().Error("init logger fail", slog.Any("err", err))
		return nil
	}

	Logger().Debug("show config", slog.Any("conf", conf))
	if conf.ShowErrCode {
		ErrorRegistry().ShowErrors()
	}

	err = utility.InitO11YTracer(&conf.O11Y, Shutdown(), Version().ServiceName)
	if err != nil {
		Logger().Error("init o11y tracer fail", slog.Any("err", err))
		return nil
	}

	return logWriter
}

//

func initLogger(filename string, conf *wlog.Config) (w io.WriteCloser, err error) {
	logger, w, err := wlog.NewSmartLogger(filename, conf)
	if err != nil {
		return nil, err
	}

	logger.Logger = logger.With(slog.String("svc", Version().ServiceName))
	Logger().PointToNew(logger)
	Logger().SetStdDefaultLevel()
	Logger().SetStdDefaultLogger()
	return
}

var defaultLogger = wlog.NewStderrLoggerWhenNormal(false)

func Logger() *wlog.Logger {
	return defaultLogger
}

//

var defaultShutdown atomic.Pointer[utility.Shutdown]

func Shutdown() *utility.Shutdown {
	return defaultShutdown.Load()
}

//

var defaultErrorRegistry = utility.NewErrorRegistry()

func ErrorRegistry() *utility.ErrorRegistry {
	return defaultErrorRegistry
}
