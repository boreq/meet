package domain

import (
	"fmt"
	"sort"

	"github.com/boreq/errors"
)

type Schedule struct {
	periods []Period
}

func NewSchedule(periods []Period) (Schedule, error) {
	s := &Schedule{}

	for _, period := range periods {
		if err := s.addPeriod(period); err != nil {
			return Schedule{}, errors.Wrapf(err, "period '%s' could not be added", period)
		}
	}

	return *s, nil
}

func MustNewSchedule(periods []Period) Schedule {
	schedule, err := NewSchedule(periods)
	if err != nil {
		panic(err)
	}

	return schedule
}

func (s Schedule) Equal(o Schedule) bool {
	a := s.Periods()
	b := o.Periods()

	if len(a) != len(b) {
		return false
	}

	sort.Slice(a, func(i, j int) bool {
		return a[i].start.Before(a[j].start)
	})

	sort.Slice(b, func(i, j int) bool {
		return b[i].start.Before(b[j].start)
	})

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func (s Schedule) Periods() []Period {
	periods := make([]Period, len(s.periods))
	copy(periods, s.periods)
	return periods
}

func (s *Schedule) addPeriod(period Period) error {
	for _, existingPeriod := range s.periods {
		if existingPeriod.Overlaps(period) {
			return fmt.Errorf("period overlaps the existing period '%s'", existingPeriod)
		}
	}

	s.periods = append(s.periods, period)
	return nil
}

type Period struct {
	start Time
	end   Time
}

func NewPeriod(start Time, end Time) (Period, error) {
	if !start.Before(end) {
		return Period{}, errors.New("start must be before end")
	}

	return Period{
		start: start,
		end:   end,
	}, nil
}

func MustNewPeriod(start Time, end Time) Period {
	p, err := NewPeriod(start, end)
	if err != nil {
		panic(err)
	}

	return p
}

func (p Period) Start() Time {
	return p.start
}

func (p Period) End() Time {
	return p.end
}

func (p Period) Overlaps(other Period) bool {
	// p contains other
	if p.start.Before(other.start) && p.end.After(other.end) {
		return true
	}

	// other contains p
	if other.start.Before(p.start) && other.end.After(p.end) {
		return true
	}

	// p overlaps the left side of other
	if p.start.Before(other.start) && !p.end.Before(other.start) {
		return true
	}

	// p overlaps the right side of other
	if p.end.After(other.end) && !p.start.After(other.end) {
		return true
	}

	return false
}

func (p Period) String() string {
	return fmt.Sprintf("%s - %s", p.start, p.end)
}

func (p Period) IsZero() bool {
	return p == Period{}
}

type Time struct {
	hour   int
	minute int
}

func NewTime(hour int, minute int) (Time, error) {
	t := Time{
		hour:   hour,
		minute: minute,
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return Time{}, fmt.Errorf("incorrect time '%s'", t)
	}

	return t, nil
}

func MustNewTime(hour int, minute int) Time {
	t, err := NewTime(hour, minute)
	if err != nil {
		panic(err)
	}

	return t
}

func (t Time) Hour() int {
	return t.hour
}

func (t Time) Minute() int {
	return t.minute
}

func (t Time) Before(o Time) bool {
	if t.hour == o.hour {
		return t.minute < o.minute

	}
	return t.hour < o.hour
}

func (t Time) After(o Time) bool {
	if t.hour == o.hour {
		return t.minute > o.minute

	}
	return t.hour > o.hour
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.hour, t.minute)
}
