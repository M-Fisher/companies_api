package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isCountryAllowed(t *testing.T) {
	type args struct {
		action  Action
		country string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: `Create from CY - ok`,
			args: args{
				action:  ActionCompanyCreate,
				country: "CY",
			},
			want: true,
		},
		{
			name: `Delete from CY - ok`,
			args: args{
				action:  ActionCompanyDelete,
				country: "CY",
			},
			want: true,
		},
		{
			name: `Creating from other country not allowed`,
			args: args{
				action:  ActionCompanyCreate,
				country: "US",
			},
			want: false,
		},
		{
			name: `Deleting from other country not allowed`,
			args: args{
				action:  ActionCompanyDelete,
				country: "US",
			},
			want: false,
		},
		{
			name: `Unknown action not allowed`,
			args: args{
				action:  -1,
				country: "US",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCountryAllowed(tt.args.action, tt.args.country)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_service_IsActionAllowed(t *testing.T) {
	type args struct {
		action Action
		ip     string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: `Creating from CY allowed`,
			args: args{
				action: ActionCompanyCreate,
				ip:     `192.200.200.124`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: `Deleting from CY allowed`,
			args: args{
				action: ActionCompanyDelete,
				ip:     `192.200.200.124`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: `Creating from US not allowed`,
			args: args{
				action: ActionCompanyCreate,
				ip:     `192.200.150.124`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: `Deleting from US not allowed`,
			args: args{
				action: ActionCompanyDelete,
				ip:     `192.200.150.124`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: `Deleting from unknown region not allowed`,
			args: args{
				action: ActionCompanyDelete,
				ip:     `192.2`,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: `Unknown action not allowed`,
			args: args{
				action: -1,
				ip:     `200.200.150.124`,
			},
			want:    false,
			wantErr: false,
		},
	}

	clmock := new(MockIPDataProvider)
	clmock.On("GetRequestLocation", `192.200.200.124`).Return(`CY`, nil)
	clmock.On("GetRequestLocation", `192.200.150.124`).Return(`US`, nil)
	clmock.On("GetRequestLocation", `192.2`).Return(``, errors.New(`invalid ip`))
	s := &service{
		client: clmock,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.IsActionAllowed(tt.args.action, tt.args.ip)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
