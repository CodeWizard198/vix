package vix

import (
	"fmt"
	"net/http"
	"testing"
)

type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

func TestServer(t *testing.T) {
	vix := NewVIX()
	vix.GET("/user", func(ctx *Context) {
		ctx.STRING(http.StatusOK, "Hello VIX")
	})
	vix.GET("/home/*", func(ctx *Context) {
		ctx.STRING(http.StatusOK, "通配符匹配成功")
	})
	vix.GET("/login/:user/:pass", func(ctx *Context) {
		_, _ = ctx.Resp.Write([]byte(fmt.Sprintf("username: %s, password: %s", ctx.PathParam["user"], ctx.PathParam["pass"])))
	})
	vix.GET("/shop/@(^[0-9]+$)", func(ctx *Context) {
		_, _ = ctx.Resp.Write([]byte("正则匹配成功"))
	})
	vix.POST("/register", func(ctx *Context) {
		form := &RegisterForm{}
		err := ctx.BindJSONbyOpt(form, false, true)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(*form)
	})
	vix.GET("/shop", func(ctx *Context) {
		num, err := ctx.GetParamValue("id").AsInt()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("number: ", num)
	})
	vix.GET("/pod", func(ctx *Context) {
		paramMap, err := ctx.GetMoreParamValues("id", "number", "sha1")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for k, v := range paramMap {
			fmt.Println("key:", k, "value:", v.Value)
		}
		ctx.STRING(200, "获取paramMap成功")
	})
	err := vix.Run(":8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
