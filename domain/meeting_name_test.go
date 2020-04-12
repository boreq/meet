package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMeetingName(t *testing.T) {
	testCases := []struct {
		Name          string
		Value         string
		ErrorExpected bool
	}{
		{
			Name:          "empty_string",
			Value:         "",
			ErrorExpected: true,
		},
		{
			Name:          "valid_lower_case",
			Value:         "meetingname",
			ErrorExpected: false,
		},
		{
			Name:          "valid_mixed_case",
			Value:         "MeetingName",
			ErrorExpected: false,
		},
		{
			Name:          "delimited_with_pause",
			Value:         "meeting-name",
			ErrorExpected: false,
		},
		{
			Name:          "delimited_with_underscore",
			Value:         "meeting_name",
			ErrorExpected: false,
		},
		{
			Name:          "invalid_special",
			Value:         "meeting$name",
			ErrorExpected: true,
		},
		{
			Name:          "invalid_space",
			Value:         "meeting name",
			ErrorExpected: true,
		},
		{
			Name:          "invalid_too_long",
			Value:         strings.Repeat("a", 200),
			ErrorExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			meetingName, err := NewMeetingName(testCase.Value)
			if testCase.ErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.False(t, meetingName.IsZero())
			}
		})
	}
}

func TestMeetingNameIsCaseInsensitive(t *testing.T) {
	require.True(t, MustNewMeetingName("meetingname") == MustNewMeetingName("MeetingName"))
}
