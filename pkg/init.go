package pkg

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func init() {
	defaultShutdown.Store(utility.NewShutdown(context.Background(), 0, Logger().Logger))
}

// Init initializes the necessary default global variables
func Init(conf *configs.Config) {
	logger := wlog.NewLogger(os.Stdout, &conf.Logger, wlog.DefaultFormat()...)
	logger.Logger = logger.With(slog.String("svc", Version().ServiceName))
	Logger().PointToNew(logger)
	Logger().SetStdDefaultLevel()
	Logger().SetStdDefaultLogger()
}

// 大部分的情況不允許 pkg目錄 以外的程式碼
// 去改變 default global variable (variable is pointer)
// 透過函數存取物件情況, 不在以上討論的範圍
//
// 想清楚使用情境 (use case) 到底要修改 pointer 本身 or 修改 pointer to 物件
// 比如以下兩個例子, 根據情境, 體現不同物件的設計方式, 所採用的做法不同
//
// 若想要透過 config 檔案控制 shutdown wait seconds
// 必須使用 SetShutdown 替換 defaultShutdown 指標
//
// 若想要透過 config 檔案控制 logger 行為
// 不應該替換 defaultLogger 指標, 而是變改指向的物件 PointToNew
// 這樣才可以讓 shutdown 使用的 logger 也改變行為

var defaultLogger = wlog.NewLoggerWhenNormalRun(false)

func Logger() *wlog.Logger {
	return defaultLogger
}

//

var defaultShutdown atomic.Pointer[utility.Shutdown]

func Shutdown() *utility.Shutdown {
	return defaultShutdown.Load()
}

func SetShutdown(s *utility.Shutdown) {
	defaultShutdown.Store(s)
}

//

var defaultErrorRegistry = utility.NewErrorRegistry()

func ErrorRegistry() *utility.ErrorRegistry {
	return defaultErrorRegistry
}
