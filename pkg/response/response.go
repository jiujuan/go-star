package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(http.StatusOK, body{Code: 0, Msg: "success", Data: data})
		return
	}
	// 这里可以扩展统一的错误码映射
	c.JSON(http.StatusBadRequest, body{Code: 1, Msg: err.Error()})
}