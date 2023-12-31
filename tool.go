package vix

import (
	"fmt"
	"net/http"
	"os"
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

func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func (c *Context) setHeaderJSON(data []byte) {
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
}

func (c *Context) setHeaderSTRING(data string) {
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
}

func (c *Context) setHeaderBYTE(data []byte) {
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
}
