package companies

import (
	"context"
	"fmt"

	"github.com/M-Fisher/companies_api/app/internal/models"
)

func (s *service) GetCompanies(ctx context.Context, params models.Company) ([]*models.Company, error) {
	select {
	case <-ctx.Done():
		s.log.Debug("Skipping getting companies due to ctx cancelled")
		return nil, ctx.Err()
	default:
		companies, err := s.db.Queries.GetCompanies(ctx, makeDBCompanyFromRequest(&params))
		if err != nil {
			return nil, fmt.Errorf("failed to get companies: %w", err)
		}
		res := make([]*models.Company, len(companies))
		for i, c := range companies {
			res[i] = makeCompanyFromDBResponse(c)
		}
		return res, nil
	}
}
