package model

import "gorm.io/gorm"

// User 对应数据库表 users
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:32"`
	Password string `gorm:"size:128"` // 已加密
}

// TableName 显式指定表名
func (User) TableName() string {
	return "users"
}