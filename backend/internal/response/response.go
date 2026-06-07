package response

import "github.com/gofiber/fiber/v2"

// Envelope is the standard response wrapper.
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{Success: true, Data: data})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Envelope{Success: true, Data: data})
}

func Fail(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(Envelope{
		Success: false,
		Error:   &APIError{Code: code, Message: message},
	})
}

func BadRequest(c *fiber.Ctx, msg string) error {
	return Fail(c, fiber.StatusBadRequest, "bad_request", msg)
}
func Unauthorized(c *fiber.Ctx, msg string) error {
	return Fail(c, fiber.StatusUnauthorized, "unauthorized", msg)
}
func Forbidden(c *fiber.Ctx, msg string) error {
	return Fail(c, fiber.StatusForbidden, "forbidden", msg)
}
func NotFound(c *fiber.Ctx, msg string) error {
	return Fail(c, fiber.StatusNotFound, "not_found", msg)
}
func TooManyRequests(c *fiber.Ctx, msg string) error {
	return Fail(c, fiber.StatusTooManyRequests, "rate_limited", msg)
}
func Internal(c *fiber.Ctx) error {
	return Fail(c, fiber.StatusInternalServerError, "internal_error", "Terjadi kesalahan internal")
}
