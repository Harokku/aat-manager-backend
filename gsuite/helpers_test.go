package gsuite

import (
	"testing"
)

func TestCheckA1Validity(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "Valid Input",
			s:    "Test!A1:B2",
			want: true,
		},
		{
			name: "Double letter interval",
			s:    "Test!AA12:BB43",
			want: true,
		},
		{
			name: "Invalid Input No Delimiter",
			s:    "TestA1:B2",
			want: false,
		},
		{
			name: "Invalid Input Special Characters",
			s:    "Test!@#:$%^",
			want: false,
		},
		{
			name: "Invalid Input Lowercase Letters",
			s:    "Test!a1:b2",
			want: false,
		},
		{
			name: "Invalid Input Spaces Included",
			s:    "Test !A1:B2",
			want: false,
		},
		{
			name: "Empty Input",
			s:    "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkA1Validity(tt.s); got != tt.want {
				t.Errorf("checkA1Validity()/n got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColNumToName(t *testing.T) {
	testCases := []struct {
		name     string
		colNum   int
		expected string
	}{
		{
			name:     "Test Single Digit",
			colNum:   1,
			expected: "A",
		},
		{
			name:     "Test Double Digit",
			colNum:   26,
			expected: "Z",
		},
		{
			name:     "Test Triple Digit",
			colNum:   27,
			expected: "AA",
		},
		{
			name:     "Test Triple High Digit",
			colNum:   52,
			expected: "AZ",
		},
		{
			name:     "Test Large Number",
			colNum:   700,
			expected: "ZX",
		},
		{
			name:     "Test Zero Value",
			colNum:   0,
			expected: "",
		},
		{
			name:     "Test Negative Number",
			colNum:   -5,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := colNumToName(tc.colNum)
			if result != tc.expected {
				t.Errorf("colNumToName(%d) = %s; expected %s", tc.colNum, result, tc.expected)
			}
		})
	}
}
