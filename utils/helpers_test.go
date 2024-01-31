package utils

import (
	"os"
	"testing"
)

func TestReadEnvOrPanic(t *testing.T) {
	cases := []struct {
		name        string
		envName     string
		envValue    string
		expectPanic bool
	}{
		{
			name:        "Environment variable exists",
			envName:     "EXISTING_VARIABLE",
			envValue:    "value",
			expectPanic: false,
		},
		{
			name:        "Environment variable does not exist",
			envName:     "NON_EXISTING_VARIABLE",
			expectPanic: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envName, tt.envValue)
				defer os.Unsetenv(tt.envName)
			}
			defer func() {
				if r := recover(); r != nil && !tt.expectPanic {
					t.Errorf("%s should not have panicked but did", tt.name)
				} else if r == nil && tt.expectPanic {
					t.Errorf("%s should have panicked but did not", tt.name)
				}
			}()
			ReadEnvOrPanic(tt.envName)
		})
	}
}
