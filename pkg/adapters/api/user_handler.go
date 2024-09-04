package api

import (
	"context"
	"net/http"

	"github.com/KScaesar/go-layout/pkg/app"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/gin-gonic/gin"
)

func RegisterUser(tx utility.Transaction, svc app.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req app.RegisterUserRequest

		err := tx(c.Request.Context(), func(ctxTX context.Context) error {
			return svc.RegisterUser(ctxTX, &req)
		})
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
