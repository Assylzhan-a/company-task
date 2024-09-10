package repository

import (
	"context"
	"github.com/assylzhan-a/company-task/internal/company/domain"
	"github.com/assylzhan-a/company-task/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) domain.CompanyRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (id, name, description, amount_of_employees, registered, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query, company.ID, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type, company.CreatedAt, company.UpdatedAt)
	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"companies_name_key\" (SQLSTATE 23505)" {
			return errors.NewConflictError("A company with this name already exists")
		}
		return errors.NewInternalServerError("Failed to create company")
	}
	return nil
}

func (r *postgresRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `
		UPDATE companies
		SET name = $2, description = $3, amount_of_employees = $4, registered = $5, type = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, company.ID, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type, company.UpdatedAt)
	if err != nil {
		return errors.NewInternalServerError("Failed to update company")
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM companies WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errors.NewInternalServerError("Failed to delete company")
	}
	if result.RowsAffected() == 0 {
		return errors.NewNotFoundError("Company not found")
	}
	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	query := `SELECT * FROM companies WHERE id = $1`
	var company domain.Company
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&company.ID, &company.Name, &company.Description, &company.AmountOfEmployees,
		&company.Registered, &company.Type, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errors.NewNotFoundError("Company not found")
		}
		return nil, errors.NewInternalServerError("Failed to get company")
	}
	return &company, nil
}

func (r *postgresRepository) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	query := `SELECT * FROM companies WHERE name = $1`
	var company domain.Company
	err := r.pool.QueryRow(ctx, query, name).Scan(
		&company.ID, &company.Name, &company.Description, &company.AmountOfEmployees,
		&company.Registered, &company.Type, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.NewNotFoundError("Company not found")
		}
		return nil, errors.NewInternalServerError("Failed to get company")
	}
	return &company, nil
}
