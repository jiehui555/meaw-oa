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

type loginResponse struct {
	Token string `json:"token"`
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

	token, err := common.GenerateToken(user.ID)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		return common.Fail(c, fiber.StatusInternalServerError, "internal error")
	}

	slog.Info("user logged in", "name", user.Name, "id", user.ID)

	return common.Success(c, loginResponse{Token: token})
}
