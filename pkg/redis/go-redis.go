package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/jiujuan/go-star/pkg/config"
)

// Client 在原有 *redis.Client 上再包一层，方便后期扩展（如链路追踪、指标）
type Client struct {
	*redis.Client
}

// 全局单例
var Rdb *Client

// ---------- 初始化 ----------
func New(cfg *config.Config) *Client {
	opts := &redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	r := redis.NewClient(opts)
	// 探活
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := r.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("redis ping failed: %w", err))
	}
	Rdb = &Client{r}
	return Rdb
}

// ---------- Key 通用 ----------
func (r *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

func (r *Client) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Del(ctx, keys...).Result()
}

func (r *Client) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.Client.Expire(ctx, key, ttl).Result()
}

func (r *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

// ---------- String ----------
func (r *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *Client) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Client) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// ---------- Hash ----------
func (r *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.HSet(ctx, key, values...).Err()
}

func (r *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

func (r *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

func (r *Client) HDel(ctx context.Context, key string, fields ...string) error {
	return r.Client.HDel(ctx, key, fields...).Err()
}

// ---------- List ----------
func (r *Client) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.Client.LPush(ctx, key, values...).Result()
}

func (r *Client) RPop(ctx context.Context, key string) (string, error) {
	return r.Client.RPop(ctx, key).Result()
}

func (r *Client) LLen(ctx context.Context, key string) (int64, error) {
	return r.Client.LLen(ctx, key).Result()
}

// ---------- Set ----------
func (r *Client) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.Client.SAdd(ctx, key, members...).Result()
}

func (r *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}

func (r *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.Client.SIsMember(ctx, key, member).Result()
}

// ---------- Sorted Set ----------
func (r *Client) ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return r.Client.ZAdd(ctx, key, members...).Result()
}

func (r *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.Client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ---------- Pipeline ----------
func (r *Client) Pipeline(ctx context.Context, fn func(p redis.Pipeliner) error) error {
	pipe := r.Client.Pipeline()
	if err := fn(pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	return err
}

// ---------- Transaction ----------
func (r *Client) Tx(ctx context.Context, fn func(tx redis.Pipeliner) error) error {
	pipe := r.Client.TxPipeline()
	if err := fn(pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	return err
}

// ---------- Lua ----------
func (r *Client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.Client.Eval(ctx, script, keys, args...).Result()
}

// ---------- 分布式锁 ----------
type Lock struct {
	key    string
	value  string
	ttl    time.Duration
	client *Client
}

// NewLock 创建可重入/不可重入锁的包装
func (r *Client) NewLock(key string, ttl time.Duration) *Lock {
	return &Lock{
		key:    key,
		value:  "lock:" + uuid.New().String(),
		ttl:    ttl,
		client: r,
	}
}

// Acquire 获取锁（基于 SET NX PX）
func (l *Lock) Acquire(ctx context.Context) (bool, error) {
	return l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
}

// Release 释放锁（基于 Lua 脚本保证原子性）
func (l *Lock) Release(ctx context.Context) error {
	script := `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end`
	_, err := l.client.Eval(ctx, script, []string{l.key}, l.value)
	return err
}

// ---------- 错误简化 ----------
var ErrNil = redis.Nil

// ---------- Fx 模块 ----------
var Module = fx.Provide(New)