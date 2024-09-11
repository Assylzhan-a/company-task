package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}

func NewUser(username, password string) (*User, error) {
	if err := validateCredentials(username, password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
	}, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func validateCredentials(username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if username == "" {
		return ErrEmptyUsername
	}

	if password == "" {
		return ErrEmptyPassword
	}

	return nil
}
