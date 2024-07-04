package fiberx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gopkg-dev/karma/errors"
)

func DefaultNotFoundHandler(c *fiber.Ctx) error {
	return errors.NotFound("Route '[%s] %s' does not exist in this API!", c.Method(), c.OriginalURL())
}

func DefaultLimitReachedHandler(_ *fiber.Ctx) error {
	return errors.TooManyRequests("TooManyRequests", "Too Many Request")
}

func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(errors.Error{
			Code:    fiberErr.Code,
			Reason:  errors.UnknownReason,
			Message: fiberErr.Message,
		})
	}
	e := errors.FromError(err)
	return c.Status(e.Code).JSON(e)
}
