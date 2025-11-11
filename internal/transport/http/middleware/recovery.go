package middleware

import (
	"net/http"
	"power-supply-sys/pkg/common"
	"power-supply-sys/pkg/logger"
	httputil "power-supply-sys/internal/transport/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 恢复panic中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录详细的 panic 信息
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回统一的错误响应
				c.JSON(http.StatusInternalServerError, httputil.Response{
					Code:    int(common.ErrCodeInternalError),
					Message: "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

