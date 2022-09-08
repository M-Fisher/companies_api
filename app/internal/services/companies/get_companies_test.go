package companies

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/internal/models"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

func Test_service_GetCompanies(t *testing.T) {
	ctxCancelled, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name    string
		ctx     context.Context
		params  models.Company
		want    []*models.Company
		wantErr error
	}{
		{
			name: `Get companies DB error`,
			ctx:  context.Background(),
			params: models.Company{
				Name: "TestError",
			},
			wantErr: fmt.Errorf("failed to get companies: %w", errors.New(`db error`)),
		},
		{
			name: `Get companies context canceled`,
			ctx:  ctxCancelled,
			params: models.Company{
				Name: "Test",
			},
			wantErr: errors.New(`context canceled`),
		},
		{
			name: `Get companies OK`,
			ctx:  context.Background(),
			params: models.Company{
				Name: "TestOk",
			},
			wantErr: nil,
			want: []*models.Company{
				{
					Name: "Comp1",
				},
				{
					Name: "Comp2",
				},
			},
		},
	}

	qmocks := new(postgres.MockCompaniesQueries)
	qmocks.On("GetCompanies", mock.Anything, postgres.Company{Name: "TestError"}).Return(nil, errors.New(`db error`))
	qmocks.On("GetCompanies", mock.Anything, postgres.Company{Name: "TestOk"}).Return([]*postgres.Company{
		{
			Name: "Comp1",
		},
		{
			Name: "Comp2",
		},
	}, nil)

	s := &service{
		db: &postgres.DB{
			Queries: qmocks,
			Log:     zap.NewExample(),
		},
		log: zap.NewExample(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetCompanies(tt.ctx, tt.params)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
