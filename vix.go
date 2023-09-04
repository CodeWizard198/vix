package vix

import (
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
}

func NewVIX(options ...HTTPServerOption) *HTTPServer {
	server := &HTTPServer{route: newRouter()}
	l := buildLogging()
	logMiddle := l.build()
	options = append(options, ServerWithMiddleware(logMiddle))
	for _, opt := range options {
		opt(server)
	}
	return server
}

func ServerWithMiddleware(middles ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.middles = middles
	}
}

//func(h *HTTPServer) AddServerMiddleware(middleware Middleware){
//	h.middles = append(h.middles, middleware)
//}

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
	// 查找路由树 执行命中的业务逻辑
	root(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	success, match := h.route.checkRouter(ctx.Req.Method, ctx.Req.URL.Path)
	if !success || match == nil {
		ctx.STRING(http.StatusNotFound, "")
		return
	}
	ctx.PathParam = match.param
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
