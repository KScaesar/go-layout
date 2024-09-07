package inject

import (
	"context"
	"net/http"
	"time"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/adapters/api"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/gin-gonic/gin"
)

func NewHttpMux(conf *configs.Config, svc *Service) *gin.Engine {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	router.Use(
		gin.Recovery(),
		utility.GinHttpObservability(svc.Name),
	)

	v1 := router.Group("/api/v1")

	v1.POST("/users", api.RegisterUser(svc.Transaction, svc.UserService))
	v1.GET("/users", api.QueryMultiUser(svc.UserService))

	router.GET("", utility.GinRoutes(router, ""))
	return router
}

func ServeApiServer(port string, handler http.Handler) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		utility.DefaultShutdown.Notify(err)
	}()
	utility.DefaultShutdown.AddPriorityShutdownAction(0, "api", func() error {
		return server.Shutdown(context.Background())
	})
}
