package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT 密钥
var jwtSecret = []byte("meaw-oa-secret")

// Claims JWT 声明结构
type Claims struct {
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"` // 令牌类型：access 或 refresh
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成访问令牌
// 有效期为 2 小时
func GenerateAccessToken(userID uint) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken 生成刷新令牌
// 有效期为 30 天
func GenerateRefreshToken(userID uint) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT 令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}
