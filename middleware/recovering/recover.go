package recovering

import "github.com/CodeWizard198/vix"

// MiddlewareRecover 从panic中恢复服务
type MiddlewareRecover struct {
	StatusCode   int
	ResponseData []byte
	Operate      func(ctx *vix.Context)
}

func BuildRecover(status int, data []byte) *MiddlewareRecover {
	return &MiddlewareRecover{
		StatusCode:   status,
		ResponseData: data,
	}
}

func (m *MiddlewareRecover) Build() vix.Middleware {
	return func(next vix.HandleFunc) vix.HandleFunc {
		return func(ctx *vix.Context) {

			defer func() {
				if pan := recover(); pan != nil {
					ctx.ResponseStatusCode = m.StatusCode
					ctx.ResponseData = m.ResponseData
					if m.Operate != nil {
						m.Operate(ctx)
					}
				}
			}()

			next(ctx)
		}
	}
}
