package companies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/services/events"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

func (s *service) UpdateCompany(ctx context.Context, compID uint64, company models.Company) (*models.Company, error) {
	select {
	case <-ctx.Done():
		s.log.Debug("Skipping getting companies due to ctx cancelled")
		return nil, ctx.Err()
	default:
		var res *models.Company
		err := s.db.ExecTx(ctx, func(q *postgres.Queries) error {
			compID, err := q.UpdateCompany(ctx, compID, makeDBCompanyFromRequest(&company))
			if err != nil {
				return fmt.Errorf("failed to update company: %w", err)
			}

			dbCompany, err := q.GetCompanyByID(ctx, compID)
			if err != nil {
				return fmt.Errorf("failed to get updated company: %w", err)
			}
			res = makeCompanyFromDBResponse(dbCompany)
			resJSON, err := json.Marshal(res)
			if err != nil {
				return fmt.Errorf("failed to marshal company for event: %w", err)
			}

			err = s.event.SendEvent(ctx, events.EventCompanyUpdated, resJSON)
			if err != nil {
				res = nil
				return fmt.Errorf("failed to send company update event: %w", err)
			}

			return nil
		})

		return res, err
	}
}
