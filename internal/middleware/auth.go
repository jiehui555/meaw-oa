package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"github.com/jiehui555/meaw-oa/internal/common"
	"github.com/jiehui555/meaw-oa/internal/model"
)

func Auth(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return common.FailWithCode(c, 401, "缺少授权头")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return common.FailWithCode(c, 401, "授权头格式无效")
		}

		claims, err := common.ParseToken(parts[1])
		if err != nil {
			return common.FailWithCode(c, 401, "令牌无效")
		}

		if claims.TokenType != "access" {
			return common.FailWithCode(c, 401, "令牌类型无效")
		}

		var user model.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			return common.FailWithCode(c, 401, "用户不存在")
		}

		c.Locals("user", user)
		return c.Next()
	}
}
