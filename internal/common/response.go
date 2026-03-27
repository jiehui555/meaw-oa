package common

import "github.com/gofiber/fiber/v3"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: "成功",
		Data:    data,
	})
}

func Fail(c fiber.Ctx, httpStatus int, message string) error {
	return c.Status(httpStatus).JSON(Response{
		Code:    httpStatus,
		Message: message,
	})
}

func FailWithCode(c fiber.Ctx, code int, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    code,
		Message: message,
	})
}
