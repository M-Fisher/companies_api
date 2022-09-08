package events

import (
	"context"

	"go.uber.org/zap"
)

type EventName string

const (
	EventCompanyCreated EventName = `company_create`
	EventCompanyUpdated EventName = `company_update`
	EventCompanyDeleted EventName = `company_delete`
)

func (s *service) SendEvent(ctx context.Context, event EventName, data []byte) error {
	s.log.Debug("Sending event", zap.String("event_name", string(event)), zap.String("data", string(data)))
	return s.producer.WriteMessage(ctx, string(event), data)
}
