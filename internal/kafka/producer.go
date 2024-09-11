package kafka

import (
	"context"
	"github.com/assylzhan-a/company-task/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(ctx context.Context, topic string, key, value []byte) error
	Close() error
}

type CompanyProducer struct {
	writer *kafka.Writer
	logger *logger.Logger
}

func NewProducer(brokers []string, logger *logger.Logger) Producer {
	return &CompanyProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
	}
}

func (p *CompanyProducer) Produce(ctx context.Context, topic string, key, value []byte) error {
	message := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	}

	err := p.writer.WriteMessages(ctx, message)
	if err != nil {
		p.logger.Error("Failed to write message to Kafka", "error", err)
		return err
	}

	return nil
}

func (p *CompanyProducer) Close() error {
	return p.writer.Close()
}
