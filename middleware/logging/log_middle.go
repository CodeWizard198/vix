package logging

import (
	"fmt"
	"github.com/CodeWizard198/vix"
	"net/http"
	"time"
)

type MiddlewareLogging struct {
}

func BuildLogging() *MiddlewareLogging {
	return &MiddlewareLogging{}
}

func (m *MiddlewareLogging) Build() vix.Middleware {
	return func(next vix.HandleFunc) vix.HandleFunc {
		return func(ctx *vix.Context) {

			start := time.Now()

			defer func(start time.Time) {
				cost := time.Since(start)
				co := fmt.Sprintf("%v", cost)
				if ctx.ResponseStatusCode == 0 {
					ctx.ResponseStatusCode = http.StatusInternalServerError
				}
				l := logs{
					IP:     vix.GetIP(ctx.Req),
					Method: ctx.Req.Method,
					Path:   ctx.Req.URL.Path,
					Times:  co,
					Status: ctx.ResponseStatusCode,
				}
				l.logging()
			}(start)

			next(ctx)
		}
	}
}
