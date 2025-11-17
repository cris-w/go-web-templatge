package http

import (
	"net/http"

	"power-supply-sys/pkg/common"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    int(common.ErrCodeSuccess),
		Message: "success",
		Data:    data,
	})
}

// HandleError 统一错误处理
func HandleError(c *gin.Context, err error) {
	if err == nil {
		SuccessResponse(c, nil)
		return
	}

	// 如果是 AppError，使用其中的信息
	if common.IsAppError(err) {
		appErr := err.(*common.AppError)
		c.JSON(appErr.Code.GetHTTPStatus(), Response{
			Code:    int(appErr.Code),
			Message: appErr.Message,
			Data:    appErr.Data,
		})
		return
	}

	// 未知错误，返回内部错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(common.ErrCodeInternalError),
		Message: "服务器内部错误",
	})
}

// HandleSuccess 统一成功响应
func HandleSuccess(c *gin.Context, data any) {
	SuccessResponse(c, data)
}

// PageResponse 分页响应
type PageResponse struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

// SuccessPageResponse 分页成功响应
func SuccessPageResponse(c *gin.Context, list any, total int64, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    int(common.ErrCodeSuccess),
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
func HandlePageSuccess(c *gin.Context, list any, total int64, page, size int) {
	SuccessPageResponse(c, list, total, page, size)
}
