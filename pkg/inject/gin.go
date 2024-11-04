package inject

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/adapters/api"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wgin"
)

func NewGinRouter(conf *pkg.Config, db *gorm.DB, svc *Service) *gin.Engine {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.NoRoute(func(c *gin.Context) {
		err := pkg.ErrNotExists.(*utility.CustomError)
		c.JSON(err.HttpStatus(), gin.H{"msg": err.Error()})
	})

	o11yLogger1, o11yLogger2 := wgin.O11YLogger(conf.Http.Debug, conf.O11Y.EnableTrace, pkg.Logger())
	router.Use(
		gin.Recovery(),
		wgin.O11YTrace(conf.O11Y.EnableTrace),
		wgin.O11YMetric(pkg.Version().ServiceName),
		o11yLogger1,
		o11yLogger2,
		wgin.GormTX(db, nil, pkg.Logger()),
	)

	router.GET("/:id", api.HelloGin(conf.Hack))
	router.GET("/logger/level", wgin.ChangeLoggerLevel(conf.Hack, pkg.Logger()))

	v1 := router.Group("/api/v1")

	v1.POST("/users", api.RegisterUser(svc.UserService))
	v1.GET("/users", api.QueryMultiUser(svc.UserService))

	return router
}

func ServeGin(port string, handler http.Handler) {
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: handler,
	}
	id := fmt.Sprintf("gin(%p)", server)
	pkg.Shutdown().AddPriorityShutdownAction(0, id, func() error {
		return server.Shutdown(context.Background())
	})

	go func() {
		pkg.Logger().Info("api start", slog.String("url", "http://0.0.0.0:"+port))
		err := server.ListenAndServe()
		pkg.Shutdown().Notify(err)
	}()
}
