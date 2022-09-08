package events

import (
	"context"

	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/config"
)

type EventsService interface {
	SendEvent(ctx context.Context, event EventName, data []byte) error
	Stop() error
}

type Producer interface {
	Close() error
	WriteMessage(ctx context.Context, eventName string, data []byte) error
}

type service struct {
	producer Producer
	log      *zap.Logger
}

func NewService(conf *config.Kafka, producerClient Producer, log *zap.Logger) *service {
	return &service{
		producer: producerClient,
		log:      log,
	}
}

func (s *service) Stop() error {
	s.log.Info("Closing kafka connection")
	return s.producer.Close()
}
