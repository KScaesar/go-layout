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

	logWriter, err := initLogger(&conf.Logger)
	if err != nil {
		Logger().Error("init logger fail", slog.Any("err", err))
		return nil
	}

	Logger().Debug("show config", slog.Any("conf", conf))
	// ErrorRegistry().ShowErrors()

	err = utility.InitO11YTracer(&conf.O11Y, Shutdown(), Version().ServiceName)
	if err != nil {
		Logger().Error("init o11y tracer fail", slog.Any("err", err))
		return nil
	}

	return logWriter
}

//

func initLogger(conf *wlog.Config) (w io.WriteCloser, err error) {
	var logger *wlog.Logger

	if conf.Filename != "" {
		w, err = wlog.NewRotateWriter(conf.Filename, -1)
		if err != nil {
			return
		}
		logger = wlog.NewFileLogger(w, conf)
	} else {
		w = os.Stderr
		logger = wlog.NewConsoleLogger(w, conf)
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
