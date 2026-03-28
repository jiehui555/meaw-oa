package model

import "time"

// Captcha 验证码模型
type Captcha struct {
	ID        uint      `gorm:"primaryKey"`
	CaptchaID string    `gorm:"size:64;uniqueIndex;not null"` // 验证码唯一标识
	Answer    string    `gorm:"size:10;not null"`             // 验证码答案
	ExpiresAt time.Time `gorm:"not null;index"`               // 过期时间
	CreatedAt time.Time
}
