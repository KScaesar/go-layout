package utility

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GinRoutes(router *gin.Engine, hack Hack) func(*gin.Context) {
	routes := make([]string, 0)
	for _, route := range router.Routes() {
		routes = append(routes, fmt.Sprintf("%-8v %v", route.Method, route.Path))
	}
	resp := strings.Join(routes, "\n")

	return func(c *gin.Context) {
		CtxGetLogger(c.Request.Context()).Info("good",
			slog.Time("now", time.Now()),
		)

		if hack.Challenge(c.Query("hack_api")) {
			c.String(200, resp)
			return
		}
		c.String(200, "hello")
	}
}
