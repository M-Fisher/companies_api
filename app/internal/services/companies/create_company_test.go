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

func TestCreateCompanyContextCancelled(t *testing.T) {
	ctxCancelled, cancel := context.WithCancel(context.Background())
	cancel()

	s := &service{
		log: zap.NewExample(),
	}
	got, err := s.CreateCompany(ctxCancelled, models.Company{})

	assert.Equal(t, errors.New("context canceled"), err)
	assert.Equal(t, uint64(0), got)
}

func TestCreateCompanyOk(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery("INSERT INTO companies (.+) VALUES (.+) RETURNING id").
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
		events.EventCompanyCreated,
		[]byte(`{"id":1,"name":"test","code":"TST","country":"","website":"","phone":""}`),
	).Return(nil)

	authmocks := new(auth.MockAuthService)
	s := &service{
		db:          dbMock,
		authService: authmocks,
		event:       evmocks,
		log:         zap.NewExample(),
	}
	got, err := s.CreateCompany(context.Background(), models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Nil(t, err)
	assert.Equal(t, uint64(1), got)
}

func TestCreateCompanyInsertError(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery("INSERT INTO companies (.+) VALUES (.+) RETURNING id").
		WillReturnError(errors.New(`insert err`))
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
	got, err := s.CreateCompany(context.Background(), models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Equal(t, fmt.Errorf("failed to create company: %w", fmt.Errorf("query: %w", errors.New("insert err"))), err)
	assert.Equal(t, uint64(0), got)
}

func TestCreateCompanySelectError(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery("INSERT INTO companies (.+) VALUES (.+) RETURNING id").
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
	got, err := s.CreateCompany(context.Background(), models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Equal(t, fmt.Errorf("failed to get created company: %w", fmt.Errorf("query: %w", fmt.Errorf("scany: query one result row: %w", errors.New("select err")))), err)
	assert.Equal(t, uint64(1), got)
}

func TestCreateCompanySendEventErr(t *testing.T) {
	pgxMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("Failed to start pgxmock: %v", err)
	}
	pgxMock.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	pgxMock.ExpectQuery("INSERT INTO companies (.+) VALUES (.+) RETURNING id").
		WillReturnRows(pgxMock.NewRows([]string{`id`}).AddRow(uint64(1)))
	pgxMock.ExpectQuery("SELECT .+ FROM companies WHERE id = .*").
		WillReturnRows(pgxMock.NewRows([]string{`id`, `name`, `code`, `country`, `website`, `phone`}).AddRow(uint64(1), `test`, `TST`, ``, ``, ``))
	pgxMock.ExpectRollback()

	dbMock, err := postgres.NewTestPostgres(pgxMock, &config.DB{}, zap.NewExample())
	if err != nil {
		t.Fatalf("Failed to init test postgres")
	}

	evmocks := new(events.MockEventsService)
	evmocks.On(
		"SendEvent",
		mock.Anything,
		events.EventCompanyCreated,
		[]byte(`{"id":1,"name":"test","code":"TST","country":"","website":"","phone":""}`),
	).Return(errors.New(`send error`))

	authmocks := new(auth.MockAuthService)
	s := &service{
		db:          dbMock,
		authService: authmocks,
		event:       evmocks,
		log:         zap.NewExample(),
	}
	got, err := s.CreateCompany(context.Background(), models.Company{
		Name: "test",
		Code: "TST",
	})

	assert.Equal(t, fmt.Errorf("failed to send company create event: %w", errors.New(`send error`)), err)
	assert.Equal(t, uint64(1), got)
}
