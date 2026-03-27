package model

import "time"

type Captcha struct {
	ID        uint      `gorm:"primaryKey"`
	CaptchaID string    `gorm:"size:64;uniqueIndex;not null"`
	Answer    string    `gorm:"size:10;not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}
