package fiberx

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetToken Get access token from header or query parameter
func GetToken(c *fiber.Ctx) string {
	var token string
	auth := c.Get(fiber.HeaderAuthorization)
	prefix := "Bearer "

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = auth
	}

	if token == "" {
		token = c.Query("accessToken")
	}

	return token
}

// GetClientIP Get client IP address
func GetClientIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Real-IP"); net.ParseIP(ip) != nil {
		return ip
	}
	if xff := c.Get(fiber.HeaderXForwardedFor); xff != "" {
		for _, ip := range strings.Split(xff, ", ") {
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	return c.IP()
}
