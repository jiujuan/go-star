package repository

import (
	"context"

	"github.com/jiujuan/go-star/internal/model"
	"github.com/jiujuan/go-star/pkg/db"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(d *db.DB) *UserRepo {
	return &UserRepo{db: d}
}

// Create 插入一条用户记录
func (r *UserRepo) Create(ctx context.Context, u *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindByUsername 根据用户名查询
func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &u, err
}