package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         string
		ErrorExpected bool
	}{
		{
			Name:          "empty_string",
			Value:         "",
			ErrorExpected: false,
		},
		{
			Name:          "real_name",
			Value:         "John Smith",
			ErrorExpected: false,
		},
		{
			Name:          "nickname",
			Value:         "user_123",
			ErrorExpected: false,
		},
		{
			Name:          "invalid_too_long",
			Value:         strings.Repeat("a", 200),
			ErrorExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := NewParticipantName(testCase.Value)
			if testCase.ErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
