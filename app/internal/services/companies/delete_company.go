package companies

import (
	"context"
	"fmt"

	"github.com/M-Fisher/companies_api/app/internal/services/events"
)

func (s *service) DeleteCompany(ctx context.Context, companyID uint64) error {
	select {
	case <-ctx.Done():
		s.log.Debug("Skipping deleting company due to ctx cancelled")
		return ctx.Err()
	default:
		err := s.db.Queries.DeleteCompany(ctx, companyID)
		if err != nil {
			return fmt.Errorf("failed to delete company: %w", err)
		}
		err = s.event.SendEvent(ctx, events.EventCompanyDeleted, []byte(fmt.Sprintf(`{"id": %d}`, companyID)))
		if err != nil {
			return fmt.Errorf("failed to send company delete event: %w", err)
		}

		return nil
	}
}
