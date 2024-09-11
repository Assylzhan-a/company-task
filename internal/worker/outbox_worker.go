package worker

import (
	"context"
	"github.com/assylzhan-a/company-task/internal/kafka"
	r "github.com/assylzhan-a/company-task/internal/ports/repository"
	"github.com/assylzhan-a/company-task/pkg/logger"
	"time"
)

type OutboxWorker struct {
	repo     r.CompanyRepository
	producer kafka.Producer
	logger   *logger.Logger
}

func NewOutboxWorker(repo r.CompanyRepository, producer kafka.Producer, logger *logger.Logger) *OutboxWorker {
	return &OutboxWorker{
		repo:     repo,
		producer: producer,
		logger:   logger,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.ProcessOutboxEvents(ctx); err != nil {
				w.logger.Error("Failed to process outbox events", "error", err)
			}
		}
	}
}

func (w *OutboxWorker) ProcessOutboxEvents(ctx context.Context) error {
	events, err := w.repo.GetOutboxEvents(ctx, 100)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := w.producer.Produce(ctx, event.EventType, nil, event.Payload); err != nil {
			w.logger.Error("Failed to produce Kafka message", "error", err, "event_id", event.ID)
			continue
		}

		if err := w.repo.DeleteOutboxEvent(ctx, event.ID); err != nil {
			w.logger.Error("Failed to delete outbox event", "error", err, "event_id", event.ID)
			continue
		}
	}

	return nil
}
