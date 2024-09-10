// internal/user/domain/user.go

package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}

type UserRepository interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
}

type UserUseCase interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
