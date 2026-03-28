package handler

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"

	"github.com/jiehui555/meaw-oa/internal/common"
	"github.com/jiehui555/meaw-oa/internal/model"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	DB *gorm.DB
}

// NewCaptchaHandler 创建验证码处理器实例
func NewCaptchaHandler(db *gorm.DB) *CaptchaHandler {
	return &CaptchaHandler{DB: db}
}

// GetC 生成并返回验证码图片
func (h *CaptchaHandler) GetC(c fiber.Ctx) error {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		slog.Error("生成验证码失败", "error", err)
		return common.Fail(c, fiber.StatusInternalServerError, "生成验证码失败")
	}

	answer := cp.Store.Get(id, false)

	captcha := model.Captcha{
		CaptchaID: id,
		Answer:    answer,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := h.DB.Create(&captcha).Error; err != nil {
		slog.Error("保存验证码失败", "error", err)
		return common.Fail(c, fiber.StatusInternalServerError, "保存验证码失败")
	}

	return common.Success(c, fiber.Map{
		"captcha_id":  id,
		"captcha_img": b64s,
	})
}
