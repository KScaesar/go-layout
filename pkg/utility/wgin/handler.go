package wgin

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
	"github.com/gin-gonic/gin"
)

func GinRoutes(router *gin.Engine, hack utility.Hack, logger *wlog.Logger) gin.HandlerFunc {
	routes := make([]string, 0)
	for _, route := range router.Routes() {
		routes = append(routes, fmt.Sprintf("%-8v %v", route.Method, route.Path))
	}
	resp := strings.Join(routes, "\n")

	return func(c *gin.Context) {
		logger.CtxGetLogger(c.Request.Context()).Info("good",
			slog.Time("print_time", time.Now()),
			slog.Any("print_fn", GinO11YMetric),
		)

		if hack.Challenge(c.Query("hack_api")) {
			c.String(200, resp)
			return
		}
		c.String(200, "hello")
	}
}

func GinSetLoggerLevel(hack utility.Hack, logger *wlog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if hack.Challenge(c.Query("hack_level")) {

			switch c.Query("level") {
			case "info":
				logger.SetLevel(slog.LevelInfo)
				logger.SetStdDefaultLevel()

			case "debug":
				logger.SetLevel(slog.LevelDebug)
				logger.SetStdDefaultLevel()
			}

			c.JSON(http.StatusOK, gin.H{"level": logger.Level()})
			return
		}
	}
}
