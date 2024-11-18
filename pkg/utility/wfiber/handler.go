package wfiber

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func ChangeLoggerLevel(hack utility.Hack, wlogger *wlog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		v1 := hack.Value()
		v2 := c.Query("hack")
		wlogger.Debug("print hack", slog.Group("hack",
			slog.String("v1", v1),
			slog.String("v2", v2),
		))

		if !hack.Challenge(v2) {
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
