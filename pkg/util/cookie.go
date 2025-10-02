package util

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetCookie(ctx *fiber.Ctx, keys string, values string, expiration *time.Time) error {
	if expiration == nil {
		expireDate := time.Now().Add(24 * time.Hour)
		expiration = &expireDate // Cookie expiration time
	}
	cookie := new(fiber.Cookie)
	cookie.Name = keys
	cookie.Value = values
	cookie.Expires = *expiration
	cookie.HTTPOnly = true
	// Set the cookie
	ctx.Cookie(cookie)
	return nil
}

func SetCookieHTTPOnlyFalse(ctx *fiber.Ctx, keys string, values string, expiration *time.Time) error {
	if expiration == nil {
		expireDate := time.Now().Add(24 * time.Hour)
		expiration = &expireDate // Cookie expiration time
	}
	cookie := new(fiber.Cookie)
	cookie.Name = keys
	cookie.Value = values
	cookie.Expires = *expiration
	cookie.HTTPOnly = false
	// Set the cookie
	ctx.Cookie(cookie)
	return nil
}

func GetCookie(c *fiber.Ctx, keys string) (*string, error) {
	// Get the cookie
	cookie := c.Cookies(keys)
	if cookie == "" {
		return nil, errors.New("can not get data")
	}
	return &cookie, nil
}

func SetExpiredCookie(ctx *fiber.Ctx, keys string) error {
	// Create a cookie with a past expiration date
	expiredCookie := &fiber.Cookie{
		Name:     keys,
		Value:    "",
		Expires:  time.Unix(0, 0), // Set to a past date
		MaxAge:   -1,              // Ensure it expires immediately
		HTTPOnly: true,
	}

	// Set the cookie in the response
	ctx.Cookie(expiredCookie)
	return nil
}

func DeleteCookie(ctx *fiber.Ctx, name string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   true,
	})
}
