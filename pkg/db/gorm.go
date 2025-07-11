package db

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/jiujuan/go-star/pkg/config"
)

// DB 封装 *gorm.DB，方便后续扩展
type DB struct {
	*gorm.DB
}

// New 根据配置初始化 GORM，支持读写分离、连接池、慢查询日志
func New(cfg *config.Config) (*DB, error) {
	// 统一日志级别
	var logLevel logger.LogLevel
	switch cfg.MySQL.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}

	// GORM 配置
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("gorm open error: %w", err)
	}

	// 获取底层 sql.DB 以设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdle)
	if lt, err := time.ParseDuration(cfg.MySQL.MaxLifetime); err == nil {
		sqlDB.SetConnMaxLifetime(lt)
	}

	// 可选：读写分离（主从）示例，若不需要可删除
	// _ = db.Use(dbresolver.Register(dbresolver.Config{
	// 	Sources:  []gorm.Dialector{mysql.Open(cfg.MySQL.DSN)},
	// 	Replicas: []gorm.Dialector{mysql.Open("slave dsn")},
	// }))

	return &DB{db}, nil
}

/* --------------------------------------------------------------------
   通用 CRUD 封装
-------------------------------------------------------------------- */

// Create 插入单条记录
func (db *DB) Create(ctx context.Context, value interface{}) error {
	return db.WithContext(ctx).Create(value).Error
}

// FirstByID 根据主键查询单条
func (db *DB) FirstByID(ctx context.Context, dest interface{}, id interface{}) error {
	return db.WithContext(ctx).First(dest, id).Error
}

// FirstWhere 根据条件查询一条
func (db *DB) FirstWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.WithContext(ctx).Where(query, args...).First(dest).Error
}

// FindWhere 根据条件查询多条
func (db *DB) FindWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.WithContext(ctx).Where(query, args...).Find(dest).Error
}

// Updates 按主键更新指定字段（map / struct）
func (db *DB) Updates(ctx context.Context, model interface{}, values interface{}) error {
	return db.WithContext(ctx).Model(model).Updates(values).Error
}

// DeleteByID 按主键删除
func (db *DB) DeleteByID(ctx context.Context, model interface{}, id interface{}) error {
	return db.WithContext(ctx).Delete(model, id).Error
}

// Count 按条件计数
func (db *DB) Count(ctx context.Context, model interface{}, query string, args ...interface{}) (int64, error) {
	var total int64
	err := db.WithContext(ctx).Model(model).Where(query, args...).Count(&total).Error
	return total, err
}

// Paginate 统一分页查询
type Page struct {
	Page  int `json:"page"`  // 第几页，从 1 开始
	Size  int `json:"size"`  // 每页条数
	Total int64
}

func (p *Page) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 || p.Size > 1000 {
		p.Size = 10
	}
	return (p.Page - 1) * p.Size
}

func (p *Page) Limit() int {
	if p.Size <= 0 || p.Size > 1000 {
		p.Size = 10
	}
	return p.Size
}

func (db *DB) Paginate(ctx context.Context, dest interface{}, page *Page, query string, args ...interface{}) error {
	dbCtx := db.WithContext(ctx).Model(dest).Where(query, args...)
	// 先查总数
	if err := dbCtx.Count(&page.Total).Error; err != nil {
		return err
	}
	// 再查分页数据
	return dbCtx.Offset(page.Offset()).Limit(page.Limit()).Find(dest).Error
}

// Transaction 执行事务
func (db *DB) Transaction(ctx context.Context, fn func(tx *DB) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&DB{tx})
	})
}


//  Fx 模块
var Module = fx.Provide(New)