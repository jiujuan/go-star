package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiujuan/go-star/internal/service"
	"github.com/jiujuan/go-star/pkg/jwt"
	"github.com/jiujuan/go-star/pkg/response"
	"github.com/jiujuan/go-star/pkg/validator"
)

type AuthHandler struct {
	svc *service.UserService
}

func NewAuthHandler(svc *service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	type req struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Password string `json:"password" binding:"required,min=6,max=32"`
	}

	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		response.JSON(c, nil, validator.Struct(r))
		return
	}

	user, err := h.svc.Register(c.Request.Context(), r.Username, r.Password)
	response.JSON(c, gin.H{"id": user.ID}, err)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	type req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		response.JSON(c, nil, validator.Struct(r))
		return
	}

	token, err := h.svc.Login(c.Request.Context(), r.Username, r.Password)
	response.JSON(c, gin.H{"token": token}, err)
}

// Me 获取当前登录用户信息
func (h *AuthHandler) Me(c *gin.Context) {
	uid, _ := c.Get("current_user_id")
	user, err := h.svc.GetByID(c.Request.Context(), uid.(string))
	response.JSON(c, gin.H{"user": user}, err)
}