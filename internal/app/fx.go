package app

import (
	"github.com/jiujuan/go-star/internal/handler"
	"github.com/jiujuan/go-star/internal/middleware"
	"github.com/jiujuan/go-star/internal/repository"
	"github.com/jiujuan/go-star/internal/router"
	"github.com/jiujuan/go-star/internal/service"
	"github.com/jiujuan/go-star/pkg/cache"
	"github.com/jiujuan/go-star/pkg/config"
	"github.com/jiujuan/go-star/pkg/db"
	"github.com/jiujuan/go-star/pkg/jwt"
	"github.com/jiujuan/go-star/pkg/logger"
	"github.com/jiujuan/go-star/pkg/redis"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	config.Module,
	logger.Module,
	db.Module,
	redis.Module,
	cache.Module,
	jwt.Module,

	fx.Provide(repository.NewUserRepo),
	fx.Provide(service.NewUserService),
	fx.Provide(handler.NewAuthHandler),
	fx.Provide(router.NewRouter),
)