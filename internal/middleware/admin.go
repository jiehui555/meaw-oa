package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jiehui555/meaw-oa/internal/common"
	"github.com/jiehui555/meaw-oa/internal/model"
)

func Admin() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return common.FailWithCode(c, 401, "未认证")
		}

		if user.Name != "admin" {
			return common.FailWithCode(c, 403, "需要管理员权限")
		}

		return c.Next()
	}
}
