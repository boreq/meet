package adapters_test

import (
	"testing"

	"github.com/boreq/meet/internal/eventsourcing/adapters"
)

func RunTestMemory(t *testing.T, test Test) {
	adapter := adapters.NewMemoryPersistenceAdapter()
	test(t, adapter)
}
