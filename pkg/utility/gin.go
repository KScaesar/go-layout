package utility

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func GinRoutes(router *gin.Engine, debugKey string) func(*gin.Context) {
	routes := make([]string, 0)
	for _, route := range router.Routes() {
		routes = append(routes, fmt.Sprintf("%-8v %v", route.Method, route.Path))
	}
	resp := strings.Join(routes, "\n")
	return func(c *gin.Context) {
		if c.Query("debug") == debugKey {
			c.String(200, resp)
		}
	}
}
