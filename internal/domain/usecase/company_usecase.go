package usecase

import (
	"context"
	"encoding/json"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	r "github.com/assylzhan-a/company-task/internal/ports/repository"
	uc "github.com/assylzhan-a/company-task/internal/ports/usecase"
	"github.com/assylzhan-a/company-task/pkg/logger"
	"github.com/google/uuid"
	"time"
)

type companyUseCase struct {
	repo   r.CompanyRepository
	logger *logger.Logger
}

func NewCompanyUseCase(repo r.CompanyRepository, logger *logger.Logger) uc.CompanyUseCase {
	return &companyUseCase{repo: repo, logger: logger}
}

func (uc *companyUseCase) Create(ctx context.Context, company *entity.Company) error {
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	payload, err := json.Marshal(company)
	if err != nil {
		return err
	}

	event := &entity.OutboxEvent{
		ID:        uuid.New(),
		EventType: "company_created",
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.CreateWithOutboxEvent(ctx, company, event); err != nil {
		uc.logger.Error("Failed to create company with outbox event", "error", err, "companyID", company.ID)
		return err
	}

	return nil
}

func (uc *companyUseCase) Patch(ctx context.Context, id uuid.UUID, patch *entity.PatchCompany) error {
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

	company.UpdatedAt = time.Now()

	payload, err := json.Marshal(company)
	if err != nil {
		return err
	}

	event := &entity.OutboxEvent{
		ID:        uuid.New(),
		EventType: "company_updated",
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.UpdateWithOutboxEvent(ctx, company, event); err != nil {
		uc.logger.Error("Failed to update company with outbox event", "error", err)
		return err
	}

	return nil
}

func (uc *companyUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *companyUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	return uc.repo.GetByID(ctx, id)
}
