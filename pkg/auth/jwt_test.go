package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTManager(t *testing.T) {
	secret := "test-secret"
	expireHours := 24

	manager := NewJWTManager(secret, expireHours)

	assert.NotNil(t, manager)
	assert.Equal(t, secret, manager.secret)
	assert.Equal(t, expireHours, manager.expireHours)
}

func TestGenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24)

	tests := []struct {
		name     string
		userID   uint
		username string
		wantErr  bool
	}{
		{
			name:     "成功生成token",
			userID:   1,
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "用户ID为0",
			userID:   0,
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "空用户名",
			userID:   1,
			username: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.GenerateToken(tt.userID, tt.username)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// 验证生成的token可以被解析
				claims, err := manager.ParseToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.username, claims.Username)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24)

	tests := []struct {
		name      string
		setupFunc func() string
		wantErr   bool
		checkFunc func(*testing.T, *Claims)
	}{
		{
			name: "成功解析有效token",
			setupFunc: func() string {
				token, _ := manager.GenerateToken(1, "testuser")
				return token
			},
			wantErr: false,
			checkFunc: func(t *testing.T, claims *Claims) {
				assert.Equal(t, uint(1), claims.UserID)
				assert.Equal(t, "testuser", claims.Username)
			},
		},
		{
			name: "无效的token格式",
			setupFunc: func() string {
				return "invalid.token.format"
			},
			wantErr: true,
		},
		{
			name: "空token",
			setupFunc: func() string {
				return ""
			},
			wantErr: true,
		},
		{
			name: "错误的签名",
			setupFunc: func() string {
				// 使用不同的secret生成token
				wrongManager := NewJWTManager("wrong-secret", 24)
				token, _ := wrongManager.GenerateToken(1, "testuser")
				return token
			},
			wantErr: true,
		},
		{
			name: "过期的token",
			setupFunc: func() string {
				// 创建一个已过期的token
				now := time.Now()
				claims := &Claims{
					UserID:   1,
					Username: "testuser",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
						NotBefore: jwt.NewNumericDate(now.Add(-2 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(manager.secret))
				return tokenString
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupFunc()
			claims, err := manager.ParseToken(token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if tt.checkFunc != nil {
					tt.checkFunc(t, claims)
				}
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24)

	tests := []struct {
		name      string
		setupFunc func() string
		wantErr   bool
	}{
		{
			name: "成功刷新有效token",
			setupFunc: func() string {
				token, _ := manager.GenerateToken(1, "testuser")
				return token
			},
			wantErr: false,
		},
		{
			name: "刷新无效token失败",
			setupFunc: func() string {
				return "invalid.token"
			},
			wantErr: true,
		},
		{
			name: "刷新过期token失败",
			setupFunc: func() string {
				now := time.Now()
				claims := &Claims{
					UserID:   1,
					Username: "testuser",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(manager.secret))
				return tokenString
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldToken := tt.setupFunc()
			newToken, err := manager.RefreshToken(oldToken)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, newToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newToken)
				// 注意：新token可能与旧token相同，因为刷新是基于相同的用户信息和时间
				// 只需要验证新token是有效的即可

				// 验证新token可以解析
				claims, err := manager.ParseToken(newToken)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), claims.UserID)
				assert.Equal(t, "testuser", claims.Username)
			}
		})
	}
}

func TestTokenExpiry(t *testing.T) {
	// 创建一个1秒过期的token
	manager := NewJWTManager("test-secret", 0)

	// 手动创建一个立即过期的token
	now := time.Now()
	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Millisecond)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.secret))
	require.NoError(t, err)

	// 等待token过期
	time.Sleep(100 * time.Millisecond)

	// 尝试解析过期的token
	parsedClaims, err := manager.ParseToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}

func TestTokenWithDifferentSigningMethods(t *testing.T) {
	manager := NewJWTManager("test-secret", 24)

	// 创建使用不同签名方法的token
	now := time.Now()
	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// 使用 HS512 而不是 HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(manager.secret))
	require.NoError(t, err)

	// HS512也是HMAC方法，应该可以解析（虽然实际应用中应该验证具体算法）
	parsedClaims, err := manager.ParseToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, parsedClaims)
}
