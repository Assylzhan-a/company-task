package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type CompanyType string

var ValidCompanyTypes = []CompanyType{
	"Corporations",
	"NonProfit",
	"Cooperative",
	"Sole Proprietorship",
}

type Company struct {
	ID                uuid.UUID   `json:"id" validate:"required"`
	Name              string      `json:"name" validate:"required,max=15"`
	Description       *string     `json:"description,omitempty" validate:"omitempty,max=3000"`
	AmountOfEmployees int         `json:"amount_of_employees" validate:"required,min=1"`
	Registered        bool        `json:"registered" validate:"required"`
	Type              CompanyType `json:"type" validate:"required,companyType"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type PatchCompany struct {
	Name              *string      `json:"name,omitempty" validate:"omitempty,max=15"`
	Description       *string      `json:"description,omitempty" validate:"omitempty,max=3000"`
	AmountOfEmployees *int         `json:"amount_of_employees,omitempty" validate:"omitempty,min=1"`
	Registered        *bool        `json:"registered,omitempty"`
	Type              *CompanyType `json:"type,omitempty" validate:"omitempty,companyType"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("companyType", validateCompanyType)
}

func validateCompanyType(fl validator.FieldLevel) bool {
	value := CompanyType(fl.Field().String())
	for _, t := range ValidCompanyTypes {
		if value == t {
			return true
		}
	}
	return false
}

func (c *Company) Validate() error {
	return validate.Struct(c)
}

func (pc *PatchCompany) Validate() error {
	return validate.Struct(pc)
}

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Company, error)
	GetByName(ctx context.Context, name string) (*Company, error)
}

type CompanyUseCase interface {
	Create(ctx context.Context, company *Company) error
	Patch(ctx context.Context, id uuid.UUID, patch *PatchCompany) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Company, error)
}
