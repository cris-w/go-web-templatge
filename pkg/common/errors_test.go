package common

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCode_GetMessage(t *testing.T) {
	tests := []struct {
		name string
		code ErrorCode
		want string
	}{
		{
			name: "成功消息",
			code: ErrCodeSuccess,
			want: "成功",
		},
		{
			name: "参数错误",
			code: ErrCodeInvalidParam,
			want: "参数错误",
		},
		{
			name: "未授权",
			code: ErrCodeUnauthorized,
			want: "未授权，请先登录",
		},
		{
			name: "资源不存在",
			code: ErrCodeNotFound,
			want: "资源不存在",
		},
		{
			name: "内部错误",
			code: ErrCodeInternalError,
			want: "服务器内部错误",
		},
		{
			name: "未知错误码",
			code: ErrorCode(9999),
			want: "未知错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code.GetMessage()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestErrorCode_GetHTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		code ErrorCode
		want int
	}{
		{
			name: "未授权错误返回401",
			code: ErrCodeUnauthorized,
			want: http.StatusUnauthorized,
		},
		{
			name: "无效Token返回401",
			code: ErrCodeInvalidToken,
			want: http.StatusUnauthorized,
		},
		{
			name: "Token过期返回401",
			code: ErrCodeTokenExpired,
			want: http.StatusUnauthorized,
		},
		{
			name: "禁止访问返回403",
			code: ErrCodeForbidden,
			want: http.StatusForbidden,
		},
		{
			name: "资源不存在返回404",
			code: ErrCodeNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "资源已存在返回409",
			code: ErrCodeAlreadyExists,
			want: http.StatusConflict,
		},
		{
			name: "参数错误返回400",
			code: ErrCodeInvalidParam,
			want: http.StatusBadRequest,
		},
		{
			name: "内部错误返回500",
			code: ErrCodeInternalError,
			want: http.StatusInternalServerError,
		},
		{
			name: "数据库错误返回500",
			code: ErrCodeDatabaseError,
			want: http.StatusInternalServerError,
		},
		{
			name: "未知错误返回500",
			code: ErrorCode(9999),
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code.GetHTTPStatus()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name    string
		appErr  *AppError
		wantErr string
	}{
		{
			name: "只有code和message",
			appErr: &AppError{
				Code:    ErrCodeInvalidParam,
				Message: "参数错误",
			},
			wantErr: "code: 1001, message: 参数错误",
		},
		{
			name: "包含底层错误",
			appErr: &AppError{
				Code:    ErrCodeInternalError,
				Message: "内部错误",
				Err:     errors.New("database connection failed"),
			},
			wantErr: "code: 5000, message: 内部错误, error: database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.appErr.Error()
			assert.Equal(t, tt.wantErr, got)
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	appErr := &AppError{
		Code:    ErrCodeInternalError,
		Message: "内部错误",
		Err:     underlyingErr,
	}

	unwrapped := appErr.Unwrap()
	assert.Equal(t, underlyingErr, unwrapped)

	// 测试没有底层错误的情况
	appErrNoUnderlying := &AppError{
		Code:    ErrCodeInvalidParam,
		Message: "参数错误",
	}
	assert.Nil(t, appErrNoUnderlying.Unwrap())
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name    string
		code    ErrorCode
		message string
		want    *AppError
	}{
		{
			name:    "自定义消息",
			code:    ErrCodeInvalidParam,
			message: "用户名格式错误",
			want: &AppError{
				Code:    ErrCodeInvalidParam,
				Message: "用户名格式错误",
			},
		},
		{
			name:    "使用默认消息",
			code:    ErrCodeInvalidParam,
			message: "",
			want: &AppError{
				Code:    ErrCodeInvalidParam,
				Message: "参数错误",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.code, tt.message)
			assert.Equal(t, tt.want.Code, got.Code)
			assert.Equal(t, tt.want.Message, got.Message)
			assert.Nil(t, got.Err)
		})
	}
}

func TestNewErrorWithErr(t *testing.T) {
	underlyingErr := errors.New("connection timeout")

	tests := []struct {
		name    string
		code    ErrorCode
		message string
		err     error
	}{
		{
			name:    "自定义消息和错误",
			code:    ErrCodeInternalError,
			message: "数据库连接超时",
			err:     underlyingErr,
		},
		{
			name:    "默认消息和错误",
			code:    ErrCodeInternalError,
			message: "",
			err:     underlyingErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrorWithErr(tt.code, tt.message, tt.err)
			assert.Equal(t, tt.code, got.Code)
			if tt.message != "" {
				assert.Equal(t, tt.message, got.Message)
			} else {
				assert.Equal(t, tt.code.GetMessage(), got.Message)
			}
			assert.Equal(t, tt.err, got.Err)
		})
	}
}

func TestWrapError(t *testing.T) {
	underlyingErr := errors.New("database error")

	appErr := WrapError(ErrCodeDatabaseError, underlyingErr)

	assert.Equal(t, ErrCodeDatabaseError, appErr.Code)
	assert.Equal(t, ErrCodeDatabaseError.GetMessage(), appErr.Message)
	assert.Equal(t, underlyingErr, appErr.Err)
}

func TestErrInvalidParam(t *testing.T) {
	message := "用户名不能为空"
	err := ErrInvalidParam(message)

	assert.Equal(t, ErrCodeInvalidParam, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Nil(t, err.Err)
}

func TestErrUnauthorized(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "自定义消息",
			message: "登录已过期",
			want:    "登录已过期",
		},
		{
			name:    "默认消息",
			message: "",
			want:    "未授权，请先登录",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrUnauthorized(tt.message)
			assert.Equal(t, ErrCodeUnauthorized, err.Code)
			assert.Equal(t, tt.want, err.Message)
		})
	}
}

func TestErrForbidden(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "自定义消息",
			message: "管理员权限不足",
			want:    "管理员权限不足",
		},
		{
			name:    "默认消息",
			message: "",
			want:    "权限不足",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrForbidden(tt.message)
			assert.Equal(t, ErrCodeForbidden, err.Code)
			assert.Equal(t, tt.want, err.Message)
		})
	}
}

func TestErrNotFound(t *testing.T) {
	resource := "用户"
	err := ErrNotFound(resource)

	assert.Equal(t, ErrCodeNotFound, err.Code)
	assert.Equal(t, "用户不存在", err.Message)
}

func TestErrAlreadyExists(t *testing.T) {
	resource := "用户名"
	err := ErrAlreadyExists(resource)

	assert.Equal(t, ErrCodeAlreadyExists, err.Code)
	assert.Equal(t, "用户名已存在", err.Message)
}

func TestErrInvalidToken(t *testing.T) {
	err := ErrInvalidToken()

	assert.Equal(t, ErrCodeInvalidToken, err.Code)
	assert.Equal(t, "Token无效", err.Message)
}

func TestErrTokenExpired(t *testing.T) {
	err := ErrTokenExpired()

	assert.Equal(t, ErrCodeTokenExpired, err.Code)
	assert.Equal(t, "Token已过期", err.Message)
}

func TestErrInternal(t *testing.T) {
	underlyingErr := errors.New("panic recovered")
	err := ErrInternal(underlyingErr)

	assert.Equal(t, ErrCodeInternalError, err.Code)
	assert.Equal(t, "服务器内部错误", err.Message)
	assert.Equal(t, underlyingErr, err.Err)
}

func TestErrDatabase(t *testing.T) {
	underlyingErr := errors.New("connection lost")
	err := ErrDatabase(underlyingErr)

	assert.Equal(t, ErrCodeDatabaseError, err.Code)
	assert.Equal(t, "数据库操作失败", err.Message)
	assert.Equal(t, underlyingErr, err.Err)
}

func TestIsAppError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "AppError类型",
			err:  ErrInvalidParam("测试"),
			want: true,
		},
		{
			name: "普通error类型",
			err:  errors.New("normal error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAppError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
