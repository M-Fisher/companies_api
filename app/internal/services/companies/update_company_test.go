package companies

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/config"
	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/services/auth"
	"github.com/M-Fisher/companies_api/app/internal/services/events"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

func TestUpdateCompanyContextCancelled(t *testing.T) {
	ctxCancelled, cancel := context.WithCancel(context.Background())
	cancel()

	s := &service{
		log: zap.NewExample(),
	}
	got, err := s.UpdateCompany(ctxCancelled, 1, models.Company{})

	assert.Equal(t, errors.New("context canceled"), err)
	assert.Equal(t, (*models.Company)(nil), got)
}

func TestUpdateCompanyOk(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery("^UPDATE companies SET code = \\$1, country = \\$2, name = \\$3, phone = \\$4, website = \\$5 WHERE id = \\$6 RETURNING id$").WithArgs(`TST`, ``, `test`, ``, ``, uint64(1)).
		WillReturnRows(pgxMock.NewRows([]string{`id`}).AddRow(uint64(1)))
	pgxMock.ExpectQuery("SELECT .+ FROM companies WHERE id = .*").
		WillReturnRows(pgxMock.NewRows([]string{`id`, `name`, `code`, `country`, `website`, `phone`}).AddRow(uint64(1), `test`, `TST`, ``, ``, ``))
	pgxMock.ExpectCommit()

	dbMock, err := postgres.NewTestPostgres(pgxMock, &config.DB{}, zap.NewExample())
	if err != nil {
		t.Fatalf("Failed to init test postgres")
	}

	evmocks := new(events.MockEventsService)
	evmocks.On(
		"SendEvent",
		mock.Anything,
		events.EventCompanyUpdated,
		[]byte(`{"id":1,"name":"test","code":"TST","country":"","website":"","phone":""}`),
	).Return(nil)

	authmocks := new(auth.MockAuthService)
	s := &service{
		db:          dbMock,
		authService: authmocks,
		event:       evmocks,
		log:         zap.NewExample(),
	}
	got, err := s.UpdateCompany(context.Background(), 1, models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Nil(t, err)
	assert.Equal(t, &models.Company{
		ID:   1,
		Name: "test",
		Code: "TST",
	}, got)
}

func TestUpdateCompanyInsertError(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery("^UPDATE companies SET code = \\$1, country = \\$2, name = \\$3, phone = \\$4, website = \\$5 WHERE id = \\$6 RETURNING id$").
		WithArgs(`TST`, ``, `test`, ``, ``, uint64(1)).
		WillReturnError(errors.New(`update err`))
	pgxMock.ExpectRollback()

	dbMock, err := postgres.NewTestPostgres(pgxMock, &config.DB{}, zap.NewExample())
	if err != nil {
		t.Fatalf("Failed to init test postgres")
	}

	authmocks := new(auth.MockAuthService)
	s := &service{
		db:          dbMock,
		authService: authmocks,
		log:         zap.NewExample(),
	}
	got, err := s.UpdateCompany(context.Background(), 1, models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Equal(t, fmt.Errorf("failed to update company: %w", fmt.Errorf("query: %w", errors.New("update err"))), err)
	assert.Equal(t, (*models.Company)(nil), got)
}

func TestUpdateCompanySelectError(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery("^UPDATE companies SET code = \\$1, country = \\$2, name = \\$3, phone = \\$4, website = \\$5 WHERE id = \\$6 RETURNING id$").
		WithArgs(`TST`, ``, `test`, ``, ``, uint64(1)).
		WillReturnRows(pgxMock.NewRows([]string{`id`}).AddRow(uint64(1)))
	pgxMock.ExpectQuery("SELECT .+ FROM companies WHERE id = .*").
		WillReturnError(errors.New(`select err`))
	pgxMock.ExpectRollback()

	dbMock, err := postgres.NewTestPostgres(pgxMock, &config.DB{}, zap.NewExample())
	if err != nil {
		t.Fatalf("Failed to init test postgres")
	}

	authmocks := new(auth.MockAuthService)
	s := &service{
		db:          dbMock,
		authService: authmocks,
		log:         zap.NewExample(),
	}
	got, err := s.UpdateCompany(context.Background(), 1, models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Equal(t, fmt.Errorf("failed to get updated company: %w", fmt.Errorf("query: %w", fmt.Errorf("scany: query one result row: %w", errors.New("select err")))), err)
	assert.Equal(t, (*models.Company)(nil), got)
}

func TestUpdateCompanySendEventErr(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery(
		"^UPDATE companies SET code = \\$1, country = \\$2, name = \\$3, phone = \\$4, website = \\$5 WHERE id = \\$6 RETURNING id$",
	).
		WithArgs(`TST`, ``, `test`, ``, ``, uint64(1)).
		WillReturnRows(pgxMock.NewRows([]string{`id`}).AddRow(uint64(1)))
	pgxMock.ExpectQuery("SELECT .+ FROM companies WHERE id = .*").
		WillReturnRows(
			pgxMock.NewRows(
				[]string{`id`, `name`, `code`, `country`, `website`, `phone`},
			).
				AddRow(uint64(1), `test`, `TST`, ``, ``, ``),
		)
	pgxMock.ExpectRollback()

	dbMock, err := postgres.NewTestPostgres(pgxMock, &config.DB{}, zap.NewExample())
	if err != nil {
		t.Fatalf("Failed to init test postgres")
	}

	evmocks := new(events.MockEventsService)
	evmocks.On(
		"SendEvent",
		mock.Anything,
		events.EventCompanyUpdated,
		[]byte(`{"id":1,"name":"test","code":"TST","country":"","website":"","phone":""}`),
	).Return(errors.New(`send error`))

	authmocks := new(auth.MockAuthService)
	s := &service{
		db:          dbMock,
		authService: authmocks,
		event:       evmocks,
		log:         zap.NewExample(),
	}
	got, err := s.UpdateCompany(context.Background(), 1, models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Equal(t, fmt.Errorf("failed to send company update event: %w", errors.New(`send error`)), err)
	assert.Equal(t, (*models.Company)(nil), got)
}
