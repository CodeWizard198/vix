package vix

import (
	"fmt"
	"time"
)

type logging struct {
}

func buildLogging() *logging {
	return &logging{}
}

func (m *logging) build() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {

			start := time.Now()

			defer func(start time.Time) {
				cost := time.Since(start)
				co := fmt.Sprintf("%v", cost)
				l := logs{
					IP:     GetIP(ctx.Req),
					Method: ctx.Req.Method,
					Path:   ctx.Req.URL.Path,
					Times:  co,
					Status: ctx.ResponseCode,
				}
				l.logging()
			}(start)

			next(ctx)
		}
	}
}
