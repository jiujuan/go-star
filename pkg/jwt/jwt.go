
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jiujuan/go-star/pkg/config"
)

type Manager struct {
	secret []byte
	expire time.Duration
}

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func New(cfg *config.Config) *Manager {
	d, _ := time.ParseDuration(cfg.JWT.Expire)
	return &Manager{
		secret: []byte(cfg.JWT.Secret),
		expire: d,
	}
}

func (m *Manager) Generate(uid string) (string, error) {
	claims := Claims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*Claims); ok && t.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

var Module = fx.Provide(New)