package repository

import (
	"context"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	"github.com/google/uuid"
)

type CompanyRepository interface {
	CreateWithOutboxEvent(ctx context.Context, company *entity.Company, event *entity.OutboxEvent) error
	UpdateWithOutboxEvent(ctx context.Context, company *entity.Company, event *entity.OutboxEvent) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error)
	GetByName(ctx context.Context, name string) (*entity.Company, error)
	GetOutboxEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error)
	DeleteOutboxEvent(ctx context.Context, id uuid.UUID) error
}
