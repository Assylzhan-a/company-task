package usecase

import (
	"context"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	"github.com/google/uuid"
)

type CompanyUseCase interface {
	Create(ctx context.Context, company *entity.Company) error
	Patch(ctx context.Context, id uuid.UUID, patch *entity.PatchCompany) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error)
}
