package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/assylzhan-a/company-task/internal/domain/entity"
	r "github.com/assylzhan-a/company-task/internal/ports/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewUserRepository(db *pgxpool.Pool) r.UserRepository {
	return &userRepository{
		db:      db,
		timeout: 30 * time.Second,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var exists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", user.Username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return entity.ErrUsernameTaken
	}

	_, err = r.db.Exec(ctx,
		"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		user.ID, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	user := &entity.User{}
	err := r.db.QueryRow(ctx,
		"SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
