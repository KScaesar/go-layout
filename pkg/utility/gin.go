package utility

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GinRoutes(router *gin.Engine, hack Hack, logger *WrapLogger) func(*gin.Context) {
	routes := make([]string, 0)
	for _, route := range router.Routes() {
		routes = append(routes, fmt.Sprintf("%-8v %v", route.Method, route.Path))
	}
	resp := strings.Join(routes, "\n")

	return func(c *gin.Context) {
		call := func() {
			logger.CtxGetLogger(c.Request.Context()).Info("good",
				slog.Time("now", time.Now()),
			)
		}
		call()

		if hack.Challenge(c.Query("hack_api")) {
			c.String(200, resp)
			return
		}
		c.String(200, "hello")
	}
}
