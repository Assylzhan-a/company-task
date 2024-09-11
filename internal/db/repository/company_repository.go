package repository

import (
	"context"
	"time"

	"github.com/assylzhan-a/company-task/internal/domain/entity"
	r "github.com/assylzhan-a/company-task/internal/ports/repository"
	"github.com/assylzhan-a/company-task/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type companyRepo struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

func NewCompanyRepository(pool *pgxpool.Pool) r.CompanyRepository {
	return &companyRepo{
		pool:    pool,
		timeout: 30 * time.Second,
	}
}

func (r *companyRepo) CreateWithOutboxEvent(ctx context.Context, company *entity.Company, event *entity.OutboxEvent) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError("Failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Insert company
	_, err = tx.Exec(ctx, `
		INSERT INTO companies (id, name, description, amount_of_employees, registered, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, company.ID, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type, company.CreatedAt, company.UpdatedAt)
	if err != nil {
		return errors.NewInternalServerError("Failed to create company")
	}

	// Insert outbox event
	_, err = tx.Exec(ctx, `
		INSERT INTO outbox_events (id, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4)
	`, event.ID, event.EventType, event.Payload, event.CreatedAt)
	if err != nil {
		return errors.NewInternalServerError("Failed to create outbox event")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.NewInternalServerError("Failed to commit transaction")
	}

	return nil
}

func (r *companyRepo) UpdateWithOutboxEvent(ctx context.Context, company *entity.Company, event *entity.OutboxEvent) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError("Failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Update company
	_, err = tx.Exec(ctx, `
		UPDATE companies
		SET name = $2, description = $3, amount_of_employees = $4, registered = $5, type = $6, updated_at = $7
		WHERE id = $1
	`, company.ID, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type, company.UpdatedAt)
	if err != nil {
		return errors.NewInternalServerError("Failed to update company")
	}

	// Insert outbox event
	_, err = tx.Exec(ctx, `
		INSERT INTO outbox_events (id, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4)
	`, event.ID, event.EventType, event.Payload, event.CreatedAt)
	if err != nil {
		return errors.NewInternalServerError("Failed to create outbox event")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.NewInternalServerError("Failed to commit transaction")
	}

	return nil
}

func (r *companyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

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

func (r *companyRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, name, description, amount_of_employees, registered, type, created_at, updated_at FROM companies WHERE id = $1`
	var company entity.Company
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&company.ID, &company.Name, &company.Description, &company.AmountOfEmployees,
		&company.Registered, &company.Type, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NewNotFoundError("Company not found")
		}
		return nil, errors.NewInternalServerError("Failed to get company")
	}
	return &company, nil
}

func (r *companyRepo) GetOutboxEvents(ctx context.Context, limit int) ([]*entity.OutboxEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
		SELECT id, event_type, payload, created_at
		FROM outbox_events
		ORDER BY created_at
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get outbox events")
	}
	defer rows.Close()

	var events []*entity.OutboxEvent
	for rows.Next() {
		var event entity.OutboxEvent
		if err := rows.Scan(&event.ID, &event.EventType, &event.Payload, &event.CreatedAt); err != nil {
			return nil, errors.NewInternalServerError("Failed to scan outbox event")
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *companyRepo) DeleteOutboxEvent(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.pool.Exec(ctx, "DELETE FROM outbox_events WHERE id = $1", id)
	if err != nil {
		return errors.NewInternalServerError("Failed to delete outbox event")
	}
	return nil
}
