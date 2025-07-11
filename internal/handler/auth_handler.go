package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiujuan/go-star/internal/service"
	"github.com/jiujuan/go-star/pkg/response"
)

type AuthHandler struct {
	svc *service.UserService
}

func NewAuthHandler(svc *service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	type req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		response.JSON(c, nil, err)
		return
	}
	user, err := h.svc.Register(c.Request.Context(), r.Username, r.Password)
	response.JSON(c, gin.H{"id": user.ID}, err)
}

func (h *AuthHandler) Login(c *gin.Context)  { /* 同上略 */ }
func (h *AuthHandler) Me(c *gin.Context)     { /* 同上略 */ }