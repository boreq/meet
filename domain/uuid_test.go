package domain_test

import (
	"testing"

	"github.com/boreq/hydro/domain"
	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	testCases := []struct {
		Name        string
		Constructor uuidConstructor
	}{
		{
			Name: "rack_uuid",
			Constructor: func(s string) (uuid, error) {
				return domain.NewRackUUID(s)
			},
		},
		{
			Name: "device_uuid",
			Constructor: func(s string) (uuid, error) {
				return domain.NewDeviceUUID(s)
			},
		},
	}

	uuidTestCases := []struct {
		Name          string
		UUID          string
		ErrorExpected bool
	}{
		{
			Name:          "empty_string",
			UUID:          "",
			ErrorExpected: true,
		},
		{
			Name:          "valid_uuid",
			UUID:          "uuid",
			ErrorExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			for _, uuidTestCase := range uuidTestCases {
				t.Run(uuidTestCase.Name, func(t *testing.T) {
					uuid, err := testCase.Constructor(uuidTestCase.UUID)
					if uuidTestCase.ErrorExpected {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
						require.False(t, uuid.IsZero())
					}
				})
			}
		})
	}
}

type uuid interface {
	IsZero() bool
}

type uuidConstructor func(string) (uuid, error)
