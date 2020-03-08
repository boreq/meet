package hydro_test

import (
	"github.com/boreq/hydro/adapters/hydro"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUUIDGenerator(t *testing.T) {
	g := hydro.NewUUIDGenerator()

	u1, err := g.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, u1)

	u2, err := g.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, u2)

	require.NotEqual(t, u1, u2)
}
