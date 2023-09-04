package vix

// Middleware 洋葱模式
// 只有一个中心，逐层向外拓展
type Middleware func(next HandleFunc) HandleFunc
