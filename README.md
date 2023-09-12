<<<<<<< HEAD
=======
# vix
Go lightweight web framework

>>>>>>> 7d059dbe4d2d3a712ba7bf4ccf80bc7b11f07903
# VIX
一个基于go http的轻量级无侵入式的Go Web框架
- 支持自定义中间件
- 集成链路追踪，适配市面上主流的链路追踪服务
- 集成普罗米修斯
- 支持在Panic中重启服务
- 支持输出每次请求的请求日志
- 封装了request的常用操作，用起来更方便
- 封装了文件操作，默认接受文件类型为断点续传
- 封装了cookie的支持
- 后续将会推出对于session的支持vix-session、gorm的支持vix-gorm，以及缓存框架的支持vix-cache

# 获取VIX
``go get github.com/CodeWizard198/vix``
## 创建一个简单的HTTP服务器
```Go
<<<<<<< HEAD
v := vix.NewVIX()
=======
v := NewVIX()
>>>>>>> 7d059dbe4d2d3a712ba7bf4ccf80bc7b11f07903

v.GET("/list", func (ctx *vix.Context) {
    v.STRING(http.StatusOK, "Hello VIX")	
})

err := v.Run(":8000")

if err != nil {
    panic(err)	
}
```

## 创建带有插件的HTTP服务器
内部定义有日志打印、链路追踪、Panic恢复、普罗米修斯的中间件

```go
// 构建日志中间件 
logMiddle := logging.BuildLogging().Build()
		
// 构建Panic恢复中间件	
recoverMiddle := recovering.BuildRecover(http.StatusInternalServerError, []byte("you're panic"))
	
// 定义发生Panic要做的事
recoverMiddle.Operate = func(ctx *vix.Context) {
    fmt.Println("panic... url:", ctx.Req.URL.Path)
}
	
// 将中间件构建成责任链 初始化服务器
v := vix.NewVIX(vix.ServerWithMiddleware(logMiddle, recoverMiddle.Build()))

v.GET("/list", func(ctx *vix.Context) {
    ctx.STRING(http.StatusOK, "Hello list, you're successes")
})

v.GET("/shop", func(ctx *vix.Context) {
    panic("you're fail, panic...")
})

// 启动服务器
err := v.Run(":8080")
if err != nil {
    fmt.Println("server fail...")
}
```

## 自定义中间件
```go
logMiddle := logging.BuildLogging()

middle := func() vix.Middleware {
	return func(next vix.HandleFunc) vix.HandleFunc {
            return func(ctx *vix.Context) {
                fmt.Println(ctx.Req.Method)
                next(ctx)
            }
	}
}

v := vix.NewVIX(vix.ServerWithMiddleware(logMiddle.Build(), middle()))

v.GET("/help", func(ctx *vix.Context) {
    ctx.STRING(http.StatusOK, "Hello VIX")
})

err := v.Run(":8000")

if err != nil {
    panic(err)
}

```

## 关于
创作者：一个在读大三学生

该框架还不成熟，无法和上面上主流的框架相提并论

这是在学习go阶段写出来的一个小框架，写的不好勿喷

欢迎大家对框架升级建议，和问题的指出

<<<<<<< HEAD
联系我：2231764545@qq.com
=======
联系我：2231764545@qq.com
>>>>>>> 7d059dbe4d2d3a712ba7bf4ccf80bc7b11f07903
