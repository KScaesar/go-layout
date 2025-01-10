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
	_Shutdown.Store(utility.NewShutdown(context.Background(), -1, Logger().Slog()))
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

	logWriter, err := initLogger(conf.Filepath.Logger, &conf.Logger)
	if err != nil {
		Logger().Slog().Error("init logger fail", slog.Any("err", err))
		return nil
	}

	Logger().Slog().Debug("show config",
		slog.Any("conf", wlog.JsonValue(true, conf)),
		slog.String("node_id", conf.NodeId()),
	)
	if conf.ShowErrCode {
		ErrorRegistry().ShowErrors()
	}

	err = utility.InitO11YTracer(&conf.O11Y, Shutdown(), Version().ServiceName)
	if err != nil {
		Logger().Slog().Error("init o11y tracer fail", slog.Any("err", err))
		return nil
	}

	return logWriter
}

//

func initLogger(filename string, conf *wlog.Config) (w io.WriteCloser, err error) {
	wlogger, w, err := wlog.LoggerFactory(filename, conf)
	if err != nil {
		return nil, err
	}

	wlogger.WithAttribute(func(logger *slog.Logger) *slog.Logger {
		return logger.With(slog.String("svc", Version().ServiceName))
	})
	Logger().PointToNew(wlogger)
	Logger().SetStdDefaultLevel()
	Logger().SetStdDefaultLogger()
	return
}

var _Logger = wlog.NewStderrLoggerWhenNormal(false)

func Logger() *wlog.Logger {
	return _Logger
}

//

var _Shutdown atomic.Pointer[utility.Shutdown]

func Shutdown() *utility.Shutdown {
	return _Shutdown.Load()
}

//

var _ErrorRegistry = utility.NewErrorRegistry()

func ErrorRegistry() *utility.ErrorRegistry {
	return _ErrorRegistry
}
