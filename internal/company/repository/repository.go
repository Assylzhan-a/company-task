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

func (r *postgresRepository) CreateWithOutboxEvent(ctx context.Context, company *domain.Company, event *domain.OutboxEvent) error {
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

func (r *postgresRepository) UpdateWithOutboxEvent(ctx context.Context, company *domain.Company, event *domain.OutboxEvent) error {
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

func (r *postgresRepository) GetOutboxEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
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

	var events []*domain.OutboxEvent
	for rows.Next() {
		var event domain.OutboxEvent
		if err := rows.Scan(&event.ID, &event.EventType, &event.Payload, &event.CreatedAt); err != nil {
			return nil, errors.NewInternalServerError("Failed to scan outbox event")
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *postgresRepository) DeleteOutboxEvent(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM outbox_events WHERE id = $1", id)
	if err != nil {
		return errors.NewInternalServerError("Failed to delete outbox event")
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
