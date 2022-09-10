package companies

import (
	"context"

	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/services/auth"
	"github.com/M-Fisher/companies_api/app/internal/services/events"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

type CompaniesService interface {
	CreateCompany(ctx context.Context, company models.Company) (uint64, error)
	DeleteCompany(ctx context.Context, companyID uint64) error
	GetCompanies(ctx context.Context, params models.Company) ([]*models.Company, error)
	GetCompany(ctx context.Context, compID uint64) (*models.Company, error)
	UpdateCompany(ctx context.Context, compID uint64, company models.Company) (*models.Company, error)
}

type service struct {
	db          *postgres.DB
	authService auth.AuthService
	event       events.EventsService
	log         *zap.Logger
}

func NewService(db *postgres.DB, authService auth.AuthService, eventService events.EventsService, log *zap.Logger) *service {
	return &service{
		db:          db,
		authService: authService,
		event:       eventService,
		log:         log,
	}
}
