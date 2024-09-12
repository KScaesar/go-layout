package wfiber

import (
	"log/slog"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
	"github.com/gofiber/fiber/v2"
)

func ChangeLoggerLevel(hack utility.Hack, wlogger *wlog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !hack.Challenge(c.Query("hack_level")) {
			return nil
		}

		var update bool

		switch c.Query("level") {
		case "info":
			wlogger.SetLevel(slog.LevelInfo)
			wlogger.SetStdDefaultLevel()
			update = true

		case "debug":
			wlogger.SetLevel(slog.LevelDebug)
			wlogger.SetStdDefaultLevel()
			update = true
		}

		logger := wlogger.CtxGetLogger(c.UserContext())
		lvl := wlogger.Level().String()
		if update {
			logger.Info("update logger level", slog.String("level", lvl))
		} else {
			logger.Info("get logger level", slog.String("level", lvl))
		}

		return c.JSON(fiber.Map{"level": lvl})
	}
}
