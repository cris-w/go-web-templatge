package common

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

const (
	// 成功
	ErrCodeSuccess ErrorCode = 0

	// 客户端错误 1xxx
	ErrCodeInvalidParam   ErrorCode = 1001 // 参数错误
	ErrCodeUnauthorized   ErrorCode = 1002 // 未授权
	ErrCodeForbidden      ErrorCode = 1003 // 禁止访问
	ErrCodeNotFound       ErrorCode = 1004 // 资源不存在
	ErrCodeAlreadyExists  ErrorCode = 1005 // 资源已存在
	ErrCodeInvalidToken   ErrorCode = 1006 // Token无效
	ErrCodeTokenExpired   ErrorCode = 1007 // Token过期
	ErrCodeInvalidRequest ErrorCode = 1008 // 请求格式错误

	// 服务端错误 5xxx
	ErrCodeInternalError ErrorCode = 5000 // 内部错误
	ErrCodeDatabaseError ErrorCode = 5001 // 数据库错误
	ErrCodeCacheError    ErrorCode = 5002 // 缓存错误
	ErrCodeServiceError  ErrorCode = 5003 // 服务错误
)

// errorMessages 错误码对应的默认消息
var errorMessages = map[ErrorCode]string{
	ErrCodeSuccess:        "成功",
	ErrCodeInvalidParam:   "参数错误",
	ErrCodeUnauthorized:   "未授权，请先登录",
	ErrCodeForbidden:      "禁止访问",
	ErrCodeNotFound:       "资源不存在",
	ErrCodeAlreadyExists:  "资源已存在",
	ErrCodeInvalidToken:   "Token无效",
	ErrCodeTokenExpired:   "Token已过期",
	ErrCodeInvalidRequest: "请求格式错误",
	ErrCodeInternalError:  "服务器内部错误",
	ErrCodeDatabaseError:  "数据库操作失败",
	ErrCodeCacheError:     "缓存操作失败",
	ErrCodeServiceError:   "服务调用失败",
}

// GetMessage 获取错误码对应的消息
func (e ErrorCode) GetMessage() string {
	if msg, ok := errorMessages[e]; ok {
		return msg
	}
	return "未知错误"
}

// GetHTTPStatus 获取错误码对应的 HTTP 状态码
func (e ErrorCode) GetHTTPStatus() int {
	if e >= 1000 && e < 2000 {
		// 客户端错误
		switch e {
		case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeTokenExpired:
			return http.StatusUnauthorized
		case ErrCodeForbidden:
			return http.StatusForbidden
		case ErrCodeNotFound:
			return http.StatusNotFound
		case ErrCodeAlreadyExists:
			return http.StatusConflict
		default:
			return http.StatusBadRequest
		}
	}
	// 服务端错误
	return http.StatusInternalServerError
}

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode // 业务错误码
	Message string    // 用户可见的错误消息
	Err     error     // 底层错误（用于日志，不返回给客户端）
	Data    any       // 附加数据
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// Unwrap 实现错误链
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewError 创建新的应用错误
func NewError(code ErrorCode, message string) *AppError {
	if message == "" {
		message = code.GetMessage()
	}
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithErr 创建带底层错误的应用错误
func NewErrorWithErr(code ErrorCode, message string, err error) *AppError {
	if message == "" {
		message = code.GetMessage()
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapError 包装错误
func WrapError(code ErrorCode, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: code.GetMessage(),
		Err:     err,
	}
}

// 常用错误构造函数

// ErrInvalidParam 参数错误
func ErrInvalidParam(message string) *AppError {
	return NewError(ErrCodeInvalidParam, message)
}

// ErrUnauthorized 未授权错误
func ErrUnauthorized(message string) *AppError {
	if message == "" {
		message = "未授权，请先登录"
	}
	return NewError(ErrCodeUnauthorized, message)
}

// ErrForbidden 禁止访问错误
func ErrForbidden(message string) *AppError {
	if message == "" {
		message = "权限不足"
	}
	return NewError(ErrCodeForbidden, message)
}

// ErrNotFound 资源不存在错误
func ErrNotFound(resource string) *AppError {
	return NewError(ErrCodeNotFound, fmt.Sprintf("%s不存在", resource))
}

// ErrAlreadyExists 资源已存在错误
func ErrAlreadyExists(resource string) *AppError {
	return NewError(ErrCodeAlreadyExists, fmt.Sprintf("%s已存在", resource))
}

// ErrInvalidToken Token无效错误
func ErrInvalidToken() *AppError {
	return NewError(ErrCodeInvalidToken, "Token无效")
}

// ErrTokenExpired Token过期错误
func ErrTokenExpired() *AppError {
	return NewError(ErrCodeTokenExpired, "Token已过期")
}

// ErrInternal 内部错误
func ErrInternal(err error) *AppError {
	return NewErrorWithErr(ErrCodeInternalError, "服务器内部错误", err)
}

// ErrDatabase 数据库错误
func ErrDatabase(err error) *AppError {
	return NewErrorWithErr(ErrCodeDatabaseError, "数据库操作失败", err)
}

// IsAppError 判断是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取应用错误
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return ErrInternal(err)
}
