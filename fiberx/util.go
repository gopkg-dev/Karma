package fiberx

import (
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
