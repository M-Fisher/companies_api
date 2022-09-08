package companies

import (
	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

func makeCompanyFromDBResponse(comp *postgres.Company) *models.Company {
	return &models.Company{
		ID:      uint64(comp.ID),
		Name:    comp.Name,
		Code:    comp.Code,
		Country: comp.Country,
		Website: comp.Website,
		Phone:   comp.Phone,
	}
}

func makeDBCompanyFromRequest(comp *models.Company) postgres.Company {
	return postgres.Company{
		Name:    comp.Name,
		Code:    comp.Code,
		Country: comp.Country,
		Website: comp.Website,
		Phone:   comp.Phone,
	}
}
