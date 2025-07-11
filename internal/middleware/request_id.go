package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

// RequestID 为每个请求注入唯一 ID，并写入响应头
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDKey)
		if id == "" {
			id = uuid.New().String()
		}
		c.Set(string(RequestIDKey), id)      // 存 gin.Context，后续可取出
		c.Header(RequestIDKey, id)           // 写响应头
		c.Next()
	}
}