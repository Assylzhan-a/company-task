package usecase

import (
	"context"
	"github.com/assylzhan-a/company-task/internal/company/domain"
	"github.com/google/uuid"
	"time"
)

type companyUseCase struct {
	repo domain.CompanyRepository
}

func NewCompanyUseCase(repo domain.CompanyRepository) domain.CompanyUseCase {
	return &companyUseCase{repo: repo}
}

func (uc *companyUseCase) Create(ctx context.Context, company *domain.Company) error {
	company.ID = uuid.New()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	return uc.repo.Create(ctx, company)
}

func (uc *companyUseCase) Patch(ctx context.Context, id uuid.UUID, patch *domain.PatchCompany) error {
	company, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if patch.Name != nil {
		company.Name = *patch.Name
	}
	if patch.Description != nil {
		company.Description = patch.Description
	}
	if patch.AmountOfEmployees != nil {
		company.AmountOfEmployees = *patch.AmountOfEmployees
	}
	if patch.Registered != nil {
		company.Registered = *patch.Registered
	}
	if patch.Type != nil {
		company.Type = *patch.Type
	}

	return uc.repo.Update(ctx, company)
}

func (uc *companyUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *companyUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return uc.repo.GetByID(ctx, id)
}
