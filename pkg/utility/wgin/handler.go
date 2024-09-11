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

func ShowRoutes(router *gin.Engine, hack utility.Hack, logger *wlog.Logger) gin.HandlerFunc {
	routes := make([]string, 0)
	for _, route := range router.Routes() {
		routes = append(routes, fmt.Sprintf("%-8v %v", route.Method, route.Path))
	}
	resp := strings.Join(routes, "\n")

	return func(c *gin.Context) {
		logger.CtxGetLogger(c.Request.Context()).Info("good",
			slog.Time("print_time", time.Now()),
			slog.Any("print_fn", O11YMetric),
		)

		if hack.Challenge(c.Query("hack_api")) {
			c.String(200, resp)
			return
		}
		c.String(200, "hello")
	}
}

func ChangeLoggerLevel(hack utility.Hack, wlogger *wlog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !hack.Challenge(c.Query("hack_level")) {
			return
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

		logger := wlogger.CtxGetLogger(c.Request.Context())
		lvl := wlogger.Level().String()
		if update {
			logger.Info("update logger level", slog.String("level", lvl))
		} else {
			logger.Info("get logger level", slog.String("level", lvl))
		}

		c.JSON(http.StatusOK, gin.H{"level": lvl})
	}
}
