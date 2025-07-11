package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jiujuan/go-star/internal/handler"
	"github.com/jiujuan/go-star/internal/middleware"
)

type Router struct {
	Auth *handler.AuthHandler
}

func NewRouter(auth *handler.AuthHandler) *Router {
	return &Router{Auth: auth}
}

func (r *Router) Register(app *gin.Engine) {
	app.Use(middleware.RequestID(), middleware.Recover(), middleware.CORS())

	api := app.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.Auth.Register)
			auth.POST("/login", r.Auth.Login)
		}
		user := api.Group("/users")
		user.Use(middleware.JWT())
		{
			user.GET("/me", r.Auth.Me)
		}
	}
}