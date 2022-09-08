package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTUser_Parse(t *testing.T) {
	type args struct {
		jwtToken  string
		jwtSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    *JWTUser
		wantErr bool
	}{
		{
			name: `empty token`,
			args: args{
				jwtSecret: `some_secret`,
				jwtToken:  ``,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: `invalid token`,
			args: args{
				jwtSecret: `some_secret`,
				jwtToken:  `invalid.token`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: `all ok`,
			args: args{
				jwtToken:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`,
				jwtSecret: `test`,
			},
			want: &JWTUser{
				JWT: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.u2tRKS4GWGieHric1tRvOFVbpEVY-lb9_cijO5_Pwt0`,
				ID:  1,
			},
			wantErr: false,
		},
		{
			name: `invalid userID in token`,
			args: args{
				jwtToken:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidGVzdCJ9.U1-C54dqqj02LpDRN9VqCt5-El5hTADVtgtGp4Y9RcI`,
				jwtSecret: `test`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ju := &JWTUser{}
			got, err := ju.Parse(tt.args.jwtToken, tt.args.jwtSecret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
