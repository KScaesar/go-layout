package api

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/adapters"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func HelloGin(hack utility.Hack) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := pkg.Logger().CtxGetLogger(c.Request.Context())

		logger.Info("HelloGin",
			slog.Time("print_time", time.Now()),
			slog.Any("print_fn", HelloGin),
		)

		if hack.Challenge(c.Query("hack_api")) {
			c.String(200, "hack")
			return
		}

		c.String(200, "world")
		return
	}
}

func HelloFiber() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		c.UserContext()

		defer func() {
			if err != nil {
				err = adapters.HandleErrorByFiber(c, err)
			}
		}()

		// err = nil
		err = pkg.ErrInvalidUsername
		if err != nil {
			return fmt.Errorf("xxx_service: %w", err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"acs": "hello"})
	}
}
