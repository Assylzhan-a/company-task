package usecase

import (
	"github.com/assylzhan-a/company-task/config"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	r "github.com/assylzhan-a/company-task/internal/ports/repository"
	uc "github.com/assylzhan-a/company-task/internal/ports/usecase"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type userUseCase struct {
	userRepo r.UserRepository
}

func NewUserUseCase(userRepo r.UserRepository) uc.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) Register(username, password string) error {
	user, err := entity.NewUser(username, password)
	if err != nil {
		return err
	}

	return u.userRepo.Create(user)
}

func (u *userUseCase) Login(username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	if err := user.ComparePassword(password); err != nil {
		return "", entity.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(config.Load().JWTSecret))
}
