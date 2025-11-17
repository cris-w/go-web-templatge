package middleware

import (
	"strings"

	"power-supply-sys/pkg/auth"
	"power-supply-sys/pkg/common"
	"power-supply-sys/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// ContextKeyUserID context 中存储用户 ID 的 key
	ContextKeyUserID = "user_id"
	// ContextKeyUsername context 中存储用户名的 key
	ContextKeyUsername = "username"
)

// JWTAuth JWT 认证中间件
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header", zap.String("path", c.Request.URL.Path))
			c.Error(common.ErrUnauthorized(""))
			c.Abort()
			return
		}

		// 验证 token 格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format", zap.String("header", authHeader))
			c.Error(common.ErrInvalidToken())
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 token
		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			logger.Warn("Failed to parse token", zap.Error(err))
			c.Error(common.ErrInvalidToken())
			c.Abort()
			return
		}

		// 将用户信息存入 context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)

		logger.Debug("User authenticated",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// GetUserID 从 context 中获取用户 ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUsername 从 context 中获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// MustGetUserID 从 context 中获取用户 ID，如果不存在则 panic
func MustGetUserID(c *gin.Context) uint {
	userID, ok := GetUserID(c)
	if !ok {
		panic("user_id not found in context")
	}
	return userID
}

// MustGetUsername 从 context 中获取用户名，如果不存在则 panic
func MustGetUsername(c *gin.Context) string {
	username, ok := GetUsername(c)
	if !ok {
		panic("username not found in context")
	}
	return username
}
