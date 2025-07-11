package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recover 捕获 panic，打印堆栈，返回 500
func Recover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithFields(logrus.Fields{
			"error": recovered,
			"stack": string(debug.Stack()),
		}).Error("panic recovered")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "internal server error"})
	})
}