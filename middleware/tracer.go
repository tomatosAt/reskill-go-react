package middleware

import (
	"os"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

func GoFiberTracer(skipper *util.HttpSkipper) fiber.Handler {
	h, _ := os.Hostname()
	return otelfiber.Middleware(
		otelfiber.WithCustomAttributes(func(ctx *fiber.Ctx) []attribute.KeyValue {
			return []attribute.KeyValue{
				attribute.String("http.server.name", h),
				attribute.String("http.request.id", util.GetHttpRequestId(ctx.Context())),
				attribute.String("http.request.content-type", ctx.Get(fiber.HeaderContentType)),
			}
		}),
		otelfiber.WithNext(func(ctx *fiber.Ctx) bool {
			return skipper.Has(ctx.Method(), ctx.Path())
		}),
	)
}
