package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jiujuan/go-star/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局单例
var L *logrus.Logger

// ctxKey 用于在 context 中存放字段
type ctxKey struct{}

// Init 根据配置初始化 logrus
func Init(cfg *config.Config) {
	L = logrus.New()
	// 1. 日志级别
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	L.SetLevel(level)

	// 2. 格式
	switch cfg.Log.Format {
	case "text":
		L.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	default:
		L.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// 3. 输出目标组合
	var outs []io.Writer
	if cfg.Log.Console {
		outs = append(outs, os.Stdout)
	}
	if cfg.Log.File.Enable {
		// 使用 lumberjack 实现自动切割 + 清理
		// 创建目录
		_ = os.MkdirAll(filepath.Dir(cfg.Log.File.Path), 0o755)
		outs = append(outs, &lumberjack.Logger{
			Filename:   cfg.Log.File.Path,
			MaxSize:    parseSize(cfg.Log.File.MaxSize), // 100MB
			MaxAge:     parseDurationDays(cfg.Log.File.MaxAge),
			MaxBackups: cfg.Log.File.MaxBackups,
			Compress:   cfg.Log.File.Compress,
		})
	}

	switch len(outs) {
	case 0:
		L.SetOutput(io.Discard)
	case 1:
		L.SetOutput(outs[0])
	default:
		L.SetOutput(io.MultiWriter(outs...))
	}
}

// ---------------- 快捷函数 ----------------
func Debugf(format string, args ...interface{}) {
	L.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	L.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	L.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	L.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	L.Fatalf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	L.Panicf(format, args...)
}

// ---------------- 结构化字段 ----------------
func WithField(key string, value interface{}) *logrus.Entry {
	return L.WithField(key, value)
}
func WithFields(fields logrus.Fields) *logrus.Entry {
	return L.WithFields(fields)
}

// ---------------- Context 链路字段注入 ----------------
// FromContext 取出保存在 ctx 中的 fields
func FromContext(ctx context.Context) *logrus.Entry {
	if fields, ok := ctx.Value(ctxKey{}).(logrus.Fields); ok {
		return L.WithFields(fields)
	}
	return logrus.NewEntry(L)
}

// WithContext 把 fields 写入 ctx（中间件或 handler 使用）
func WithContext(ctx context.Context, fields logrus.Fields) context.Context {
	return context.WithValue(ctx, ctxKey{}, fields)
}

// ---------------- 内部工具 ----------------
// parseSize "100MB" -> 100
func parseSize(s string) int {
	var v int
	var unit string
	_, _ = fmt.Sscanf(s, "%d%s", &v, &unit)
	switch unit {
	case "KB", "kb":
		return v
	case "MB", "mb":
		return v
	case "GB", "gb":
		return v * 1024
	default:
		return v // 默认 MB
	}
}

// parseDurationDays "30d" -> 30
func parseDurationDays(s string) int {
	var v int
	var unit string
	_, _ = fmt.Sscanf(s, "%d%s", &v, &unit)
	if unit == "d" || unit == "D" {
		return v
	}
	return 30
}

// Fx 模块
var Module = fx.Invoke(func(cfg *config.Config) { Init(cfg) })
