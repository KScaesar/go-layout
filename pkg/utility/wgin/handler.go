package wgin

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func ChangeLoggerLevel(hack utility.Hack, wlogger *wlog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !hack.Challenge(c.Query("hack")) {
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
