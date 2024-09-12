package inject

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/adapters/api"
	"github.com/KScaesar/go-layout/pkg/utility/wgin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewGinRouter(conf *configs.Config, db *gorm.DB, svc *Service) *gin.Engine {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	o11yLogger1, o11yLogger2 := wgin.O11YLogger(conf.Http.Debug, conf.O11Y.EnableTrace, pkg.Logger())
	router.Use(
		gin.Recovery(),
		wgin.O11YTrace(conf.O11Y.EnableTrace),
		wgin.O11YMetric(pkg.Version().ServiceName, conf.O11Y.EnableTrace),
		o11yLogger1,
		o11yLogger2,
		wgin.GormTX(db, nil, pkg.Logger()),
	)

	router.GET("/", api.Hello(conf.Hack))
	router.GET("/logger/level", wgin.ChangeLoggerLevel(conf.Hack, pkg.Logger()))

	v1 := router.Group("/api/v1")

	v1.POST("/users", api.RegisterUser(svc.UserService))
	v1.GET("/users", api.QueryMultiUser(svc.UserService))

	return router
}

func ServeGin(port string, handler http.Handler) {
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		pkg.Logger().Info("api start", slog.String("url", "http://0.0.0.0:"+port))
		err := server.ListenAndServe()
		pkg.Shutdown().Notify(err)
	}()
	pkg.Shutdown().AddPriorityShutdownAction(0, "api", func() error {
		return server.Shutdown(context.Background())
	})
}
