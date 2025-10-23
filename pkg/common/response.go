package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(ErrCodeSuccess),
		Message: "success",
		Data:    data,
	})
}

// ErrorResponse 错误响应（兼容旧代码）
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// HandleError 统一错误处理（推荐使用）
func HandleError(c *gin.Context, err error) {
	if err == nil {
		SuccessResponse(c, nil)
		return
	}

	// 如果是 AppError，使用其中的信息
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.Code.GetHTTPStatus(), Response{
			Code:    int(appErr.Code),
			Message: appErr.Message,
			Data:    appErr.Data,
		})
		return
	}

	// 未知错误，返回内部错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(ErrCodeInternalError),
		Message: "服务器内部错误",
	})
}

// HandleSuccess 统一成功响应
func HandleSuccess(c *gin.Context, data interface{}) {
	SuccessResponse(c, data)
}

// PageResponse 分页响应
type PageResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// SuccessPageResponse 分页成功响应
func SuccessPageResponse(c *gin.Context, list interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    int(ErrCodeSuccess),
		Message: "success",
		Data: PageResponse{
			List:  list,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}

// HandlePageSuccess 统一分页成功响应
func HandlePageSuccess(c *gin.Context, list interface{}, total int64, page, size int) {
	SuccessPageResponse(c, list, total, page, size)
}
