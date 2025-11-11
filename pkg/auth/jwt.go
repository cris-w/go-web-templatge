package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 自定义声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager JWT 管理器
type JWTManager struct {
	secret      string
	expireHours int
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(secret string, expireHours int) *JWTManager {
	return &JWTManager{
		secret:      secret,
		expireHours: expireHours,
	}
}

// GenerateToken 生成 JWT token
func (m *JWTManager) GenerateToken(userID uint, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(m.expireHours) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT token
func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 token
func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新的 token
	return m.GenerateToken(claims.UserID, claims.Username)
}
