package inject

import (
	"log/slog"
	"time"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/adapters"
	"github.com/KScaesar/go-layout/pkg/utility/wfiber"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func NewFiberRouter(conf *configs.Config, db *gorm.DB, svc *Service) *fiber.App {
	router := fiber.New(fiber.Config{
		ErrorHandler:          fiber.DefaultErrorHandler,
		AppName:               pkg.Version().ServiceName,
		DisableStartupMessage: true,
		WriteTimeout:          1 * time.Minute,
	})

	o11yLogger1, o11yLogger2 := wfiber.O11YLogger(conf.Http.Debug, conf.O11Y.EnableTrace, pkg.Logger())
	router.Use(
		recover.New(),
		cors.New(),
		wfiber.O11YTrace(conf.O11Y.EnableTrace),
		wfiber.O11YMetric(pkg.Version().ServiceName),
		o11yLogger1,
		o11yLogger2,
	)

	// 為了利用 fiber.DefaultErrorHandler, 讓 o11y mw 保證可以讀取到 http status code, 所以分為 router, root
	root := fiber.New(fiber.Config{ErrorHandler: adapters.FiberErrorHandler})
	router.Mount("", root)

	root.Use(
		wfiber.GormTX(db, nil, pkg.Logger()),
	)

	root.Get("/logger/level", wfiber.ChangeLoggerLevel(conf.Hack, pkg.Logger()))

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

	go func() {
		err := router.Listen("0.0.0.0:" + port)
		pkg.Shutdown().Notify(err)
	}()
	pkg.Shutdown().AddPriorityShutdownAction(0, "api", router.Server().Shutdown)
}
