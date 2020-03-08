package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	testCases := []struct {
		Name          string
		Address       string
		ErrorExpected bool
	}{
		{
			Name:          "empty_string",
			Address:       "",
			ErrorExpected: true,
		},
		{
			Name:          "valid_address",
			Address:       "127.0.0.1:1234",
			ErrorExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			address, err := NewAddress(testCase.Address)
			if testCase.ErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.False(t, address.IsZero())
			}
		})
	}

}
