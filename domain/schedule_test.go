package domain_test

import (
	"github.com/boreq/hydro/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchedule(t *testing.T) {
	testCases := []struct {
		Name        string
		Periods     []domain.Period
		ExpectError bool
	}{
		{
			Name:        "empty_periods",
			Periods:     []domain.Period{},
			ExpectError: false,
		},
		{
			Name: "non_overlaping_periods",
			Periods: []domain.Period{
				domain.MustNewPeriod(
					domain.MustNewTime(11, 00),
					domain.MustNewTime(12, 00),
				),
				domain.MustNewPeriod(
					domain.MustNewTime(13, 00),
					domain.MustNewTime(14, 00),
				),
			},
			ExpectError: false,
		},
		{
			Name: "overlaping_periods",
			Periods: []domain.Period{
				domain.MustNewPeriod(
					domain.MustNewTime(11, 00),
					domain.MustNewTime(13, 00),
				),
				domain.MustNewPeriod(
					domain.MustNewTime(12, 00),
					domain.MustNewTime(14, 00),
				),
			},
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := domain.NewSchedule(testCase.Periods)
			if testCase.ExpectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPeriod(t *testing.T) {
	testCases := []struct {
		Name        string
		Start       domain.Time
		End         domain.Time
		ExpectError bool
	}{
		{
			Name:        "valid_period",
			Start:       domain.MustNewTime(11, 00),
			End:         domain.MustNewTime(12, 00),
			ExpectError: false,
		},
		{
			Name:        "start_after_end",
			Start:       domain.MustNewTime(12, 00),
			End:         domain.MustNewTime(11, 00),
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			p, err := domain.NewPeriod(testCase.Start, testCase.End)
			if testCase.ExpectError {
				require.Error(t, err)
				require.True(t, p.IsZero())
			} else {
				require.NoError(t, err)
				require.False(t, p.IsZero())
			}
		})
	}
}

func TestPeriod_Overlaps(t *testing.T) {
	testCases := []struct {
		Name     string
		First    domain.Period
		Second   domain.Period
		Overlaps bool
	}{
		{
			Name: "no_overlap",
			First: domain.MustNewPeriod(
				domain.MustNewTime(11, 00),
				domain.MustNewTime(12, 00),
			),
			Second: domain.MustNewPeriod(
				domain.MustNewTime(13, 00),
				domain.MustNewTime(14, 00),
			),
			Overlaps: false,
		},
		{
			Name: "no_overlap_touching",
			First: domain.MustNewPeriod(
				domain.MustNewTime(11, 00),
				domain.MustNewTime(11, 59),
			),
			Second: domain.MustNewPeriod(
				domain.MustNewTime(12, 00),
				domain.MustNewTime(13, 00),
			),
			Overlaps: false,
		},
		{
			Name: "overlaps",
			First: domain.MustNewPeriod(
				domain.MustNewTime(11, 00),
				domain.MustNewTime(12, 00),
			),
			Second: domain.MustNewPeriod(
				domain.MustNewTime(12, 00),
				domain.MustNewTime(13, 00),
			),
			Overlaps: true,
		},
		{
			Name: "contained",
			First: domain.MustNewPeriod(
				domain.MustNewTime(11, 00),
				domain.MustNewTime(14, 00),
			),
			Second: domain.MustNewPeriod(
				domain.MustNewTime(12, 00),
				domain.MustNewTime(13, 00),
			),
			Overlaps: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Overlaps, testCase.First.Overlaps(testCase.Second))
			assert.Equal(t, testCase.Overlaps, testCase.Second.Overlaps(testCase.First))
		})
	}
}

func TestPeriod_String(t *testing.T) {
	p := domain.MustNewPeriod(
		domain.MustNewTime(1, 1),
		domain.MustNewTime(12, 12),
	)

	require.Equal(t, "01:01 - 12:12", p.String())
}

func TestPeriod_IsZero(t *testing.T) {
	require.True(t, domain.Period{}.IsZero())
}

func TestTime(t *testing.T) {
	testCases := []struct {
		Name        string
		Hour        int
		Minute      int
		ExpectError bool
	}{
		{
			Name:        "negative hour",
			Hour:        -1,
			Minute:      0,
			ExpectError: true,
		},
		{
			Name:        "negative minute",
			Hour:        0,
			Minute:      -1,
			ExpectError: true,
		},
		{
			Name:        "beginning_of_the_day",
			Hour:        00,
			Minute:      00,
			ExpectError: false,
		},
		{
			Name:        "end_of_the_day",
			Hour:        23,
			Minute:      59,
			ExpectError: false,
		},
		{
			Name:        "too_large_hour",
			Hour:        24,
			Minute:      00,
			ExpectError: true,
		},
		{
			Name:        "too_large_minute",
			Hour:        00,
			Minute:      60,
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := domain.NewTime(testCase.Hour, testCase.Minute)
			if testCase.ExpectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTime_Relativity(t *testing.T) {
	testCases := []struct {
		Name   string
		Before domain.Time
		After  domain.Time
	}{
		{
			Name:   "same_hour",
			Before: domain.MustNewTime(12, 15),
			After:  domain.MustNewTime(12, 16),
		},
		{
			Name:   "different_hour",
			Before: domain.MustNewTime(11, 15),
			After:  domain.MustNewTime(12, 15),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			require.True(t, testCase.Before.Before(testCase.After))
			require.False(t, testCase.Before.After(testCase.After))

			require.True(t, testCase.After.After(testCase.Before))
			require.False(t, testCase.After.Before(testCase.Before))
		})
	}
}

func TestTime_String(t *testing.T) {
	tm := domain.MustNewTime(0, 12)
	require.Equal(t, "00:12", tm.String())
}
