package service

import (
	"context"
	"time"

	"github.com/jiujuan/go-star/internal/model"
	"github.com/jiujuan/go-star/internal/repository"
	"github.com/jiujuan/go-star/pkg/cache"
	"github.com/jiujuan/go-star/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  *repository.UserRepo
	cache *cache.Cache
	jwt   *jwt.Manager
}

func NewUserService(repo *repository.UserRepo, cache *cache.Cache, j *jwt.Manager) *UserService {
	return &UserService{repo: repo, cache: cache, jwt: j}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &model.User{Username: username, Password: string(hash)}
	return s.repo.Create(ctx, u)
}

func (s *UserService) Login(ctx context.Context, username, password string) (string, error) {
	token, err := s.jwt.Generate(username)
	return token, err
}

func (s *UserService) CreateUser(ctx context.Context, u *model.User) error {
	return s.db.Transaction(ctx, func(tx *db.DB) error {
		// 复杂业务：先检查重名，再写入
		var exist model.User
		if err := tx.FirstWhere(ctx, &exist, "username = ?", u.Username); err == nil {
			return fmt.Errorf("username already exists")
		}
		return tx.Create(ctx, u)
	})
}