package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KScaesar/go-layout/pkg/app"
)

func RegisterUser(svc app.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req app.RegisterUserRequest

		ctx := c.Request.Context()
		err := svc.RegisterUser(ctx, &req)
		if err != nil {
			c.Error(err)
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
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
