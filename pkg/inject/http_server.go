package inject

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KScaesar/go-layout/pkg/adapters/api"
	"github.com/gin-gonic/gin"
)

func NewHttpMux(debug bool, svc *Service) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")

	v1.POST("/users", api.RegisterUser(svc.Transaction, svc.UserService))
	v1.GET("/users", api.QueryMultiUser(svc.UserService))

	return router
}

func NewHttpServer(port string, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
