package repository

import "github.com/assylzhan-a/company-task/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByUsername(username string) (*entity.User, error)
}
