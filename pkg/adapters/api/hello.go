package api

import (
	"log/slog"
	"time"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/gin-gonic/gin"
)

func Hello(hack utility.Hack) func(c *gin.Context) {
	return func(c *gin.Context) {
		pkg.Logger().CtxGetLogger(c.Request.Context()).Info("Hello",
			slog.Time("print_time", time.Now()),
			slog.Any("print_fn", Hello),
		)

		if hack.Challenge(c.Query("hack_api")) {
			c.String(200, "hack")
			return
		}

		c.String(200, "world")
		return
	}
}
