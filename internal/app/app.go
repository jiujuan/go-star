package app

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Plugin interface {
	Name() string
	Init(*App) error
}

type App struct {
	Router *gin.Engine
	LC     fx.Lifecycle
}

func Bootstrap(opts ...fx.Option) {
	fx.New(
		fx.Provide(NewApp),
		fx.Options(opts...),
		fx.Invoke(registerHooks),
	).Run()
}

func NewApp() *App {
	gin.SetMode(gin.ReleaseMode)
	return &App{Router: gin.New()}
}

func registerHooks(lc fx.Lifecycle, app *App) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Router.Run(":8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil // graceful shutdown
		},
	})
}