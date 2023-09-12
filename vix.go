package vix

import (
	"fmt"
	"net"
	"net/http"
)

// HandleFunc 处理方法
type HandleFunc func(ctx *Context)

// HTTPServerOption option模式
type HTTPServerOption func(server *HTTPServer)

var _ Server = &HTTPServer{}

// Server 抽象的Server服务器
type Server interface {
	http.Handler
	// Run 服务启动入口
	Run(addr string) error
	// AddRouter 添加路由信息
	AddRouter(method, path string, handleFunc HandleFunc)
}

type HTTPServer struct {
	route   *router
	middles []Middleware
	logF    func(message string, args ...any)
}

func NewVIX(options ...HTTPServerOption) *HTTPServer {
	server := &HTTPServer{
		route: newRouter(),
		logF: func(message string, args ...any) {
			fmt.Printf(message, args...)
		},
	}
	for _, opt := range options {
		opt(server)
	}
	return server
}

// ServerWithMiddleware 创建HTTPServerOption，创建中间件列表
func ServerWithMiddleware(middles ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.middles = middles
	}
}

// ServeHTTP 处理HTTP请求的入口
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Resp: w,
		Req:  r,
	}
	// 从后往前构建责任链
	root := h.serve
	for i := len(h.middles) - 1; i >= 0; i-- {
		root = h.middles[i](root)
	}

	var flush Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flushResponseData(ctx)
		}
	}

	// 使用flush中间件
	// 此时的flush中间件会位于中间件列表的最后
	root = flush(root)

	// 查找路由树 执行命中的业务逻辑
	root(ctx)
}

// 刷新响应信息缓存
func (h *HTTPServer) flushResponseData(ctx *Context) {
	if ctx.ResponseStatusCode == 0 {
		ctx.ResponseStatusCode = http.StatusNotFound
	}
	ctx.Resp.WriteHeader(ctx.ResponseStatusCode)
	length, err := ctx.Resp.Write(ctx.ResponseData)
	if err != nil || len(ctx.ResponseData) != length {
		fmt.Println("写入错误")
		h.logF("写入数据错误%v\n", err)
	}
}

func (h *HTTPServer) serve(ctx *Context) {
	success, match := h.route.checkRouter(ctx.Req.Method, ctx.Req.URL.Path)
	if !success || match == nil {
		ctx.ResponseStatusCode = http.StatusInternalServerError
		return
	}
	_ = ctx.Req.ParseForm()
	ctx.Form = ctx.Req.Form
	ctx.PathParam = match.param
	ctx.MatchRouter = match.pod.route
	match.pod.handler(ctx)
}

func (h *HTTPServer) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 处理一些必要的逻辑
	return http.Serve(listener, h)
}

func (h *HTTPServer) AddRouter(method, path string, handleFunc HandleFunc) {
	h.route.addRouter(method, path, handleFunc)
}

// GET Get请求
func (h *HTTPServer) GET(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodGet, path, handleFunc)
}

// POST Post请求
func (h *HTTPServer) POST(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodPost, path, handleFunc)
}

// PUT Put请求
func (h *HTTPServer) PUT(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodPut, path, handleFunc)
}

// DELETE Delete请求
func (h *HTTPServer) DELETE(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodDelete, path, handleFunc)
}

// HEAD Head请求
func (h *HTTPServer) HEAD(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodHead, path, handleFunc)
}

// PATCH Patch请求
func (h *HTTPServer) PATCH(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodPatch, path, handleFunc)
}

// CONNECT Connect请求
func (h *HTTPServer) CONNECT(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodConnect, path, handleFunc)
}

// OPTIONS Options请求
func (h *HTTPServer) OPTIONS(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodOptions, path, handleFunc)
}

// TRACE Trace请求
func (h *HTTPServer) TRACE(path string, handleFunc HandleFunc) {
	h.AddRouter(http.MethodTrace, path, handleFunc)
}
