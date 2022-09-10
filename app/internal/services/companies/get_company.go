package companies

import (
	"context"
	"fmt"

	"github.com/M-Fisher/companies_api/app/internal/models"
)

func (s *service) GetCompany(ctx context.Context, copmID uint64) (*models.Company, error) {
	select {
	case <-ctx.Done():
		s.log.Debug("Skipping getting company due to ctx cancelled")
		return nil, ctx.Err()
	default:
		company, err := s.db.Queries.GetCompanyByID(ctx, copmID)
		if err != nil {
			return nil, fmt.Errorf("failed to get company: %w", err)
		}

		return makeCompanyFromDBResponse(company), nil
	}
}
