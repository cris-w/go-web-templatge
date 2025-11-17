package middleware

import (
	httputil "power-supply-sys/internal/transport/http"
	"power-supply-sys/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果响应已经发送，不再处理错误
		if c.Writer.Written() {
			return
		}

		// 处理请求中的错误
		if len(c.Errors) > 0 {
			ginErr := c.Errors.Last()
			err := ginErr.Err

			// 记录错误日志
			logger.Error("Request error",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Error(err),
			)

			// 统一错误响应
			httputil.HandleError(c, err)
			c.Abort()
		}
	}
}
