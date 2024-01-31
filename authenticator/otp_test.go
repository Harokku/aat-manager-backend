package authenticator

import (
	"aat-manager/db"
	"errors"
	"net/mail"
	"testing"
)

func TestCheckOtpAndDelete(t *testing.T) {
	cases := []struct {
		name    string
		email   mail.Address
		otp     int
		setup   func(d *db.InMemoryDb)
		want    bool
		wantErr error
	}{
		{
			name:    "valid OTP",
			email:   mail.Address{Address: "user1@test.com"},
			otp:     1234,
			setup:   func(d *db.InMemoryDb) { d.Set("user1", "1234") },
			want:    true,
			wantErr: nil,
		},
		{
			name:    "malformed email",
			email:   mail.Address{Address: "usertest.com"},
			otp:     1234,
			setup:   func(d *db.InMemoryDb) {},
			want:    false,
			wantErr: ErrMalformedMail,
		},
		{
			name:    "non-existence user",
			email:   mail.Address{Address: "user2@test.com"},
			otp:     1234,
			setup:   func(d *db.InMemoryDb) {},
			want:    false,
			wantErr: ErrUserNotFound,
		},
		{
			name:    "invalid OTP",
			email:   mail.Address{Address: "user3@test.com"},
			otp:     1234,
			setup:   func(d *db.InMemoryDb) { d.Set("user3", "5678") },
			want:    false,
			wantErr: nil,
		},
		{
			name:    "non-numeric stored OTP",
			email:   mail.Address{Address: "user4@test.com"},
			otp:     1234,
			setup:   func(d *db.InMemoryDb) { d.Set("user4", "abcd") },
			want:    false,
			wantErr: ErrNonNumericValue,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			memoryDb := db.NewDB()
			tt.setup(memoryDb)

			got, err := CheckOtpAndDelete(tt.email, tt.otp, memoryDb)

			if got != tt.want {
				t.Errorf("CheckOtpAndDelete() = %v, want %v", got, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CheckOtpAndDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
