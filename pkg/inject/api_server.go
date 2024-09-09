package inject

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KScaesar/go-layout/configs"
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
			utility.GinO11YMetric(svc.Name, conf.O11Y.EnableTrace),
		).
		Use(utility.GinO11YLogger(conf.Http.Debug, conf.O11Y.EnableTrace)...).
		Use(
			utility.GinGormTransaction(db, []string{}),
		)

	v1 := router.Group("/api/v1")

	v1.POST("/users", api.RegisterUser(svc.UserService))
	v1.GET("/users", api.QueryMultiUser(svc.UserService))

	router.GET("", utility.GinRoutes(router, conf.Hack))
	return router
}

func ServeApiServer(port string, handler http.Handler) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	protocol := func() string { return "http" }

	go func() {
		utility.DefaultLogger().Info("api start", slog.String("url", protocol()+"://localhost:"+port))
		err := server.ListenAndServe()
		utility.DefaultShutdown().Notify(err)
	}()
	utility.DefaultShutdown().AddPriorityShutdownAction(0, "api", func() error {
		return server.Shutdown(context.Background())
	})
}
