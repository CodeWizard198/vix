package vix

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func checkPath(method, path string) {
	compile, _ := regexp.Compile(`(//)`)
	if compile.MatchString(path) {
		panic(fmt.Sprintf("注册路由[%s-%s]失败,error:路由路径不能出现连续的[\"/\"]", method, path))
	}

	if path[0] != '/' {
		panic(fmt.Sprintf("注册路由[%s-%s]失败,error:路由路径应以[\"/\"]开头", method, path))
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic(fmt.Sprintf("注册路由[%s-%s]失败,error:路由路径不能以[\"/\"]结尾", method, path))
	}
}

func GetIP(req *http.Request) string {
	var remoteAddr string
	remoteAddr = req.RemoteAddr
	if remoteAddr != "" {
		return remoteAddr
	}

	remoteAddr = req.Header.Get("ipv4")
	if remoteAddr != "" {
		return remoteAddr
	}

	remoteAddr = req.Header.Get("XForwardedFor")
	if remoteAddr != "" {
		return remoteAddr
	}

	remoteAddr = req.Header.Get("X-Forwarded-For")
	if remoteAddr != "" {
		return remoteAddr
	}

	remoteAddr = req.Header.Get("X-Real-Ip")
	if remoteAddr != "" {
		return remoteAddr
	}
	return "127.0.0.1"
}

func (c *Context) setHeaderJSON(code int, data []byte) {
	c.Resp.WriteHeader(code)
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
}

func (c *Context) setHeaderSTRING(code int, data string) {
	c.Resp.WriteHeader(code)
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
}

func (c *Context) setHeaderBYTE(code int, data []byte) {
	c.Resp.WriteHeader(code)
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
}

func (c *Context) errResponse() {
	c.Resp.WriteHeader(http.StatusInternalServerError)
	_, _ = c.Resp.Write([]byte("响应失败"))
}
