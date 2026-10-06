package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasher(t *testing.T) {

	testCase := []struct {
		name     string
		input    string
		expected string
		result   bool
	}{

		{
			name:     "normal case",
			input:    "1234",
			expected: "1234",
			result:   true,
		},
		{
			name:     "error case",
			input:    "12346",
			expected: "12345",
			result:   false,
		},
	}

	for _, tc := range testCase {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			// t.Parallel()
			hasher := NewHasher()
			expectedResult, err := hasher.HashPassword(tc.input)
			if err != nil {
				t.Errorf("Expected %s, got %s", tc.expected, expectedResult)
				return
			}
			actual := hasher.CheckPasswordHash(tc.expected, expectedResult)
			assert.Equal(t, tc.result, actual)
		})
	}

}
