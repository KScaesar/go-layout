package api

import (
	"net/http"

	"github.com/KScaesar/go-layout/pkg/app"
	"github.com/gin-gonic/gin"
)

func RegisterUser(svc app.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req app.RegisterUserRequest

		ctx := c.Request.Context()
		err := svc.RegisterUser(ctx, &req)
		if err != nil {

			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func QueryMultiUser(svc app.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req app.QueryMultiUserRequest

		ctx := c.Request.Context()
		resp, err := svc.QueryMultiUser(ctx, &req)
		if err != nil {

			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
