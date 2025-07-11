package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jiujuan/go-star/pkg/jwt"
)

const CurrentUserID = "current_user_id"

// JWT 鉴权中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == "" || tokenStr == auth {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Authorization header required"})
			return
		}

		claims, err := jwt.M.Parse(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			return
		}

		c.Set(CurrentUserID, claims.UserID)
		c.Next()
	}
}