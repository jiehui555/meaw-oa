package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/jiehui555/meaw-oa/internal/common"
	"github.com/jiehui555/meaw-oa/internal/model"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	var req loginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return common.FailWithCode(c, 400, "请求体无效")
	}

	if req.Name == "" || req.Password == "" {
		return common.FailWithCode(c, 400, "用户名和密码为必填项")
	}

	var user model.User
	if err := h.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		slog.Warn("登录失败：用户不存在", "name", req.Name)
		return common.FailWithCode(c, 401, "用户名或密码无效")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		slog.Warn("登录失败：密码错误", "name", req.Name)
		return common.FailWithCode(c, 401, "用户名或密码无效")
	}

	tokens, err := generateTokens(user.ID)
	if err != nil {
		slog.Error("生成令牌失败", "error", err)
		return common.Fail(c, fiber.StatusInternalServerError, "服务器内部错误")
	}

	slog.Info("用户登录成功", "name", user.Name, "id", user.ID)

	return common.Success(c, tokens)
}

func (h *UserHandler) Refresh(c fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return common.FailWithCode(c, 400, "请求体无效")
	}

	if req.RefreshToken == "" {
		return common.FailWithCode(c, 400, "refresh_token 为必填项")
	}

	claims, err := common.ParseToken(req.RefreshToken)
	if err != nil {
		return common.FailWithCode(c, 401, "refresh token 无效")
	}

	if claims.TokenType != "refresh" {
		return common.FailWithCode(c, 401, "令牌类型无效")
	}

	tokens, err := generateTokens(claims.UserID)
	if err != nil {
		slog.Error("生成令牌失败", "error", err)
		return common.Fail(c, fiber.StatusInternalServerError, "服务器内部错误")
	}

	return common.Success(c, tokens)
}

func generateTokens(userID uint) (*tokenResponse, error) {
	accessToken, err := common.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := common.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
