package repository

import (
	"context"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	r "github.com/assylzhan-a/company-task/internal/ports/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) r.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		user.ID, user.Username, user.Password)
	return err
}

func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow(context.Background(),
		"SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
