package adapters_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/boreq/hydro/internal/eventsourcing"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/require"
)

type Test func(t *testing.T, adapter eventsourcing.PersistenceAdapter)

type TestRunner func(t *testing.T, test Test)

func TestAdapters(t *testing.T) {
	adapters := []struct {
		Name       string
		TestRunner TestRunner
	}{
		{
			Name:       "bolt",
			TestRunner: RunTestBolt,
		},
	}

	tests := []struct {
		Name string
		Test Test
	}{
		{
			Name: "save_empty_events",
			Test: testSaveEmptyEvents,
		},
	}

	for _, adapter := range adapters {
		t.Run(adapter.Name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.Name, func(t *testing.T) {
					adapter.TestRunner(t, test.Test)
				})
			}
		})
	}
}

func testSaveEmptyEvents(t *testing.T, adapter eventsourcing.PersistenceAdapter) {
	uuid := someAggregateUUID()

	err := adapter.SaveEvents(uuid, nil)
	require.Equal(t, eventsourcing.EmptyEventsErr, err)
}

func someAggregateUUID() eventsourcing.AggregateUUID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return eventsourcing.AggregateUUID(ulid.String())
}
