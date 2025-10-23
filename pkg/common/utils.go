package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseUintParam 从 URL 参数中解析 uint 类型 ID
func ParseUintParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, ErrInvalidParam("无效的" + paramName)
	}
	return uint(id), nil
}

// GetPageInfo 获取分页信息，返回 page 和 pageSize
func GetPageInfo(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
