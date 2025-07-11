package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	*gocache.Cache
}

func New() *Cache {
	return &Cache{gocache.New(5*time.Minute, 10*time.Minute)}
}

var Module = fx.Provide(New)