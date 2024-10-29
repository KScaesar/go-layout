package inject

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/adapters"
	"github.com/KScaesar/go-layout/pkg/adapters/api"
	"github.com/KScaesar/go-layout/pkg/utility/wfiber"
)

// NewFiberRouter
// 由於 fiber 本身的限制, error handler 必須分為兩個部分來處理
//
// 1. fiber 本身的錯誤:
// 比如 api route 找不到, 依靠 config.ErrorHandler 進行處理
//
// 2. 商業邏輯 service layer 錯誤:
// 每個 api handler 都必須手動呼叫 adapters.HandleErrorByFiber
// 商業邏輯錯誤無法依靠 config.ErrorHandler 進行集中處理的原因, 有兩個因素互相影響:
//
//	2-1. 所有 mw 執行後 config.ErrorHandler 才會執行. 但 mw 在 handler 之後, 必須取得 http code
//	2-2. middleware 無法取得 handler route, ref: https://github.com/gofiber/fiber/issues/3138
func NewFiberRouter(conf *pkg.Config, db *gorm.DB, svc *Service) *fiber.App {
	router := fiber.New(fiber.Config{
		ErrorHandler:          adapters.HandleErrorByFiber,
		AppName:               pkg.Version().ServiceName,
		DisableStartupMessage: true,
	})

	o11yMetric := wfiber.O11YMetric(pkg.Version().ServiceName)
	o11yLogger1, o11yLogger2 := wfiber.O11YLogger(conf.Http.Debug, conf.O11Y.EnableTrace, pkg.Logger())
	transaction := wfiber.GormTX(db, nil, pkg.Logger())
	router.Use(
		recover.New(recover.Config{EnableStackTrace: true}),
		cors.New(),
		wfiber.O11YTrace(conf.O11Y.EnableTrace),
		o11yLogger1,
	)

	fixFiberIssue3138 := func(handler fiber.Handler) []fiber.Handler {
		return []fiber.Handler{o11yMetric, o11yLogger2, transaction, handler}
	}

	router.Get("/:id", fixFiberIssue3138(api.HelloFiber())...)
	router.Get("/logger/level", fixFiberIssue3138(wfiber.ChangeLoggerLevel(conf.Hack, pkg.Logger()))...)

	return router
}

func ServeFiber(port string, debug bool, router *fiber.App) {
	router.Hooks().OnListen(func(_ fiber.ListenData) error {
		if debug {
			wfiber.ShowRoutes(router)
		}
		pkg.Logger().Info("api start", slog.String("url", "http://0.0.0.0:"+port))
		return nil
	})
	pkg.Shutdown().AddPriorityShutdownAction(0, "api", router.Server().Shutdown)

	go func() {
		err := router.Listen("0.0.0.0:" + port)
		pkg.Shutdown().Notify(err)
	}()
}
