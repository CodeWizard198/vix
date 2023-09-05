package opentelemetry

import (
	"github.com/CodeWizard198/vix"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const tracingName = "github/CodeWizard198/vix/middleware/open-telemetry"

// MiddlewareOpenTelemetry 链路追踪
type MiddlewareOpenTelemetry struct {
	Trace trace.Tracer
}

func (m MiddlewareOpenTelemetry) BuildOpenTelemetry() vix.Middleware {
	if m.Trace == nil {
		m.Trace = otel.GetTracerProvider().Tracer(tracingName)
	}
	return func(next vix.HandleFunc) vix.HandleFunc {
		return func(ctx *vix.Context) {
			// 与上下游的调用联系起来
			extract := otel.GetTextMapPropagator().Extract(ctx.Req.Context(), propagation.HeaderCarrier(ctx.Req.Header))

			eCTX, span := m.Trace.Start(extract, "")
			// 最后要把span关闭
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.ip", vix.GetIP(ctx.Req)))

			// 将span的ctx往下传
			ctx.Req = ctx.Req.WithContext(eCTX)

			// 继续调用下一个
			next(ctx)

			span.SetName(ctx.Req.URL.Path)
		}
	}
}
