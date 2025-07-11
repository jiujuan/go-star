package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/jiujuan/go-star/pkg/config"
)

var Rdb *redis.Client

func New(cfg *config.Config) *redis.Client {
	opts := &redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: 5,
	}
	Rdb = redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := Rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	return Rdb
}

var Module = fx.Provide(New)