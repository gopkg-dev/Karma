package fiberx

import (
	"reflect"

	"github.com/gofiber/fiber/v2"

	"github.com/gopkg-dev/karma/errors"
	"github.com/gopkg-dev/karma/gormx"
	"github.com/gopkg-dev/karma/validator"
)

// Map is a shortcut for map[string]interface{}, useful for JSON returns
type Map map[string]interface{}

// ParseBody ...
func ParseBody(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return errors.BadRequest("Failed to parse body: %s", err.Error())
	}
	return nil
}

// ParseBodyAndValidate ...
func ParseBodyAndValidate(c *fiber.Ctx, out interface{}) error {
	if err := ParseBody(c, out); err != nil {
		return err
	}
	return validator.Validate(out)
}

// ParseQuery ...
func ParseQuery(c *fiber.Ctx, out interface{}) error {
	if err := c.QueryParser(out); err != nil {
		return errors.BadRequest("Failed to parse query: %s", err.Error())
	}
	return nil
}

// ParseQueryAndValidate ...
func ParseQueryAndValidate(c *fiber.Ctx, out interface{}) error {
	if err := ParseQuery(c, out); err != nil {
		return err
	}
	return validator.Validate(out)
}

// ParseParams ...
func ParseParams(c *fiber.Ctx, out interface{}) error {
	if err := c.ParamsParser(out); err != nil {
		return errors.BadRequest("Failed to parse params: %s", err.Error())
	}
	return nil
}

// ParseParamsAndValidate ...
func ParseParamsAndValidate(c *fiber.Ctx, out interface{}) error {
	if err := ParseParams(c, out); err != nil {
		return err
	}
	return validator.Validate(out)
}

// Response is a API response
type Response struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Total   int64         `json:"total,omitempty"`
	Error   *errors.Error `json:"error,omitempty"`
}

func ResSuccess(c *fiber.Ctx, v interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    v,
	})
}

func ResOK(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
	})
}

func ResPage(c *fiber.Ctx, v interface{}, pr *gormx.PaginationResult) error {
	var total int64
	if pr != nil {
		total = pr.Total
	}
	reflectValue := reflect.Indirect(reflect.ValueOf(v))
	if reflectValue.Kind() == reflect.Ptr && reflectValue.IsZero() {
		v = make([]interface{}, 0)
	}
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    v,
		Total:   total,
	})
}
