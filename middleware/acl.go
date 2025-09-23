package middleware

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/pkg/cache"
)

func NewACLMiddleware(skipper *SkipperPath, store *cache.Redis) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if skipper != nil && skipper.Test(ctx) {
			return ctx.Next()
		}
		if _, ok := aclIsBlockingState(ctx.IP(), store); ok {
			return rejectBlockingRequest(ctx, "Access denied")
		}
		return ctx.Next()
	}
}

func aclIsBlockingState(ip string, store *cache.Redis) (int64, bool) {
	if c, err := store.Client.Get(fmt.Sprintf(config.StoreBlockingState, ip)).Int64(); err == nil {
		_ = store.Client.Incr(fmt.Sprintf(config.StoreBlockingState, ip))
		return c + 1, true
	}
	return 0, false
}

func rejectBlockingRequest(ctx *fiber.Ctx, reason string) error {
	return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
		"error": reason,
	})
}
