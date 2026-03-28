package common

import "github.com/gofiber/fiber/v3"

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 返回成功响应
func Success(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    0,
		Message: "成功",
		Data:    data,
	})
}

// Fail 返回失败响应（使用自定义 HTTP 状态码）
func Fail(c fiber.Ctx, httpStatus int, message string) error {
	return c.Status(httpStatus).JSON(Response{
		Code:    httpStatus,
		Message: message,
	})
}

// FailWithCode 返回失败响应（使用业务错误码，HTTP 状态码固定为 200）
func FailWithCode(c fiber.Ctx, code int, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    code,
		Message: message,
	})
}
