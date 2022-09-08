package companies

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/M-Fisher/companies_api/app/internal/services/auth"
	"github.com/M-Fisher/companies_api/app/internal/services/events"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

func Test_service_DeleteCompany(t *testing.T) {
	ctxCancelled, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name    string
		ctx     context.Context
		compID  uint64
		wantErr error
	}{
		{
			name:    `Delete company DB error`,
			ctx:     context.Background(),
			compID:  1,
			wantErr: fmt.Errorf("failed to delete company: %w", errors.New(`db error`)),
		},
		{
			name:    `Delete company context canceled`,
			ctx:     ctxCancelled,
			compID:  2,
			wantErr: errors.New(`context canceled`),
		},
		{
			name:    `Delete company OK`,
			ctx:     context.Background(),
			compID:  3,
			wantErr: nil,
		},
		{
			name:    `Delete company - event sending failed`,
			ctx:     context.Background(),
			compID:  4,
			wantErr: fmt.Errorf("failed to send company delete event: %w", errors.New(`failed to send event`)),
		},
	}

	qmocks := new(postgres.MockCompaniesQueries)
	qmocks.On("DeleteCompany", mock.Anything, uint64(1)).Return(errors.New(`db error`))
	qmocks.On("DeleteCompany", mock.Anything, uint64(3)).Return(nil)
	qmocks.On("DeleteCompany", mock.Anything, uint64(4)).Return(nil)
	authmocks := new(auth.MockAuthService)
	evmocks := new(events.MockEventsService)
	evmocks.On("SendEvent", mock.Anything, events.EventCompanyDeleted, []byte(`{"id": 3}`)).Return(nil)
	evmocks.On("SendEvent", mock.Anything, events.EventCompanyDeleted, []byte(`{"id": 4}`)).Return(errors.New("failed to send event"))

	s := &service{
		db: &postgres.DB{
			Queries: qmocks,
			Log:     zap.NewExample(),
		},
		authService: authmocks,
		event:       evmocks,
		log:         zap.NewExample(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteCompany(tt.ctx, tt.compID)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
