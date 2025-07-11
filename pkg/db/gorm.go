package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(cfg Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
}