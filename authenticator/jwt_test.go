package authenticator

import (
	"aat-manager/utils"
	"os"
	"testing"
)

func TestCreateAndSignJWT(t *testing.T) {
	type args struct {
		user    string
		manager bool
	}
	tests := []struct {
		name      string
		args      args
		setupFunc func()
		wantErr   bool
	}{
		{
			name: "Happy Path",
			args: args{user: "TestUser", manager: true},
			setupFunc: func() {
				os.Setenv(utils.JWTSECRET, "test_secret")
				os.Setenv(utils.JWTEXPIREINMONTH, "30")
			},
			wantErr: false,
		},
		{
			name: "Missing JWT Secret",
			args: args{user: "TestUser", manager: true},
			setupFunc: func() {
				os.Setenv(utils.JWTSECRET, "")
				os.Setenv(utils.JWTEXPIREINMONTH, "30")
			},
			wantErr: true,
		},
		{
			name: "Non Numeric Expire Days",
			args: args{user: "TestUser", manager: true},
			setupFunc: func() {
				os.Setenv(utils.JWTSECRET, "test_secret")
				os.Setenv(utils.JWTEXPIREINMONTH, "invalid_days")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup function for environment.
			tt.setupFunc()
			_, err := CreateAndSignJWT(tt.args.user, tt.args.manager)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAndSignJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
