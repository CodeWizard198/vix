package vix

type Middleware func(next HandleFunc) HandleFunc
