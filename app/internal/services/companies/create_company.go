package companies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/services/events"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

func (s *service) CreateCompany(ctx context.Context, company models.Company) (uint64, error) {
	select {
	case <-ctx.Done():
		s.log.Debug("Skipping creating company due to ctx cancelled")
		return 0, ctx.Err()
	default:
		var (
			compID uint64
			err    error
		)
		err = s.db.ExecTx(ctx, func(q *postgres.Queries) error {
			compID, err = q.CreateCompany(ctx, makeDBCompanyFromRequest(&company))
			if err != nil {
				return fmt.Errorf("failed to create company: %w", err)
			}

			dbCompany, err := q.GetCompanyByID(ctx, compID)
			if err != nil {
				return fmt.Errorf("failed to get created company: %w", err)
			}
			resJSON, err := json.Marshal(makeCompanyFromDBResponse(dbCompany))
			if err != nil {
				return fmt.Errorf("failed to marshal company for event: %w", err)
			}

			err = s.event.SendEvent(ctx, events.EventCompanyCreated, resJSON)
			if err != nil {
				return fmt.Errorf("failed to send company create event: %w", err)
			}

			return err
		})

		return compID, err
	}
}
