package vix

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	router := newRouter()
	routes := []string{
		"/home/me",
		"/shop/list",
		"/home/me/info",
		"/",
		"/login/*",
		"/list/@(^[0-9]+$)/ls",
	}
	routes_test := []string{
		"/home/a",
		"/shop/list",
		"/home/me",
		"/home/me/info/e",
		"/car",
		"/",
		"/shop/list/a",
		"/login/along/20030222",
		"/list/1./ls",
	}
	for _, route := range routes {
		router.addRouter(http.MethodGet, route, func(ctx *Context) {
			_, _ = ctx.Resp.Write([]byte("Hello VIX"))
		})
	}
	for _, route := range routes_test {
		checkRouter, _ := router.checkRouter(http.MethodGet, route)
		fmt.Println(checkRouter)
	}
}
