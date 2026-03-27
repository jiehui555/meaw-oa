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
		return common.FailWithCode(c, 400, "invalid request body")
	}

	if req.Name == "" || req.Password == "" {
		return common.FailWithCode(c, 400, "name and password are required")
	}

	var user model.User
	if err := h.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		slog.Warn("login failed: user not found", "name", req.Name)
		return common.FailWithCode(c, 401, "invalid name or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		slog.Warn("login failed: wrong password", "name", req.Name)
		return common.FailWithCode(c, 401, "invalid name or password")
	}

	tokens, err := generateTokens(user.ID)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return common.Fail(c, fiber.StatusInternalServerError, "internal error")
	}

	slog.Info("user logged in", "name", user.Name, "id", user.ID)

	return common.Success(c, tokens)
}

func (h *UserHandler) Refresh(c fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return common.FailWithCode(c, 400, "invalid request body")
	}

	if req.RefreshToken == "" {
		return common.FailWithCode(c, 400, "refresh_token is required")
	}

	claims, err := common.ParseToken(req.RefreshToken)
	if err != nil {
		return common.FailWithCode(c, 401, "invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return common.FailWithCode(c, 401, "invalid token type")
	}

	tokens, err := generateTokens(claims.UserID)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return common.Fail(c, fiber.StatusInternalServerError, "internal error")
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
