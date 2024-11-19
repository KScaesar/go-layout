package wfiber

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ShowRoutes(router *fiber.App) {
	out := make([]string, 0, 10)

	out = append(out, "")
	for _, route := range router.GetRoutes(true) {
		if route.Method == http.MethodHead {
			continue
		}

		out = append(out, fmt.Sprintf(" %-8v %v", route.Method, route.Path))

		handler := route.Handlers[len(route.Handlers)-1]
		rv := reflect.ValueOf(handler)
		out = append(out, fmt.Sprintf("\u001B[90m  └─ %8s\u001B[0m", runtime.FuncForPC(rv.Pointer()).Name()))
	}
	out = append(out, "")

	fmt.Println(strings.Join(out, "\n"))
}
