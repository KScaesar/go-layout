package inject

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/adapters/api"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHttpMux(conf *configs.Config, db *gorm.DB, svc *Service) *gin.Engine {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	router.
		Use(
			gin.Recovery(),
			utility.GinO11YTrace(conf.O11Y.EnableTrace),
			utility.GinO11YMetric(pkg.DefaultVersion().ServiceName, conf.O11Y.EnableTrace),
		).
		Use(utility.GinO11YLogger(conf.Http.Debug, conf.O11Y.EnableTrace, pkg.DefaultLogger())...).
		Use(
			utility.GinGormTransaction(db, []string{}),
		)

	router.GET("/logger/level", utility.GinSetLoggerLevel(conf.Hack, pkg.DefaultLogger()))

	v1 := router.Group("/api/v1")

	v1.POST("/users", api.RegisterUser(svc.UserService))
	v1.GET("/users", api.QueryMultiUser(svc.UserService))

	router.GET("", utility.GinRoutes(router, conf.Hack, pkg.DefaultLogger()))
	return router
}

func ServeApiServer(port string, handler http.Handler) {
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	protocol := func() string { return "http" }

	go func() {
		pkg.DefaultLogger().Info("api start", slog.String("url", protocol()+"://0.0.0.0:"+port))
		err := server.ListenAndServe()
		pkg.DefaultShutdown().Notify(err)
	}()
	pkg.DefaultShutdown().AddPriorityShutdownAction(0, "api", func() error {
		return server.Shutdown(context.Background())
	})
}
