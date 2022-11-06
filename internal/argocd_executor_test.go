package argocd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_durationStringToContext(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		ctx, cancel, err := durationStringToContext("")
		require.NoError(t, err)
		t.Cleanup(cancel)
		assert.Equal(t, context.Background(), ctx)
	})

	t.Run("invalid", func(t *testing.T) {
		_, _, err := durationStringToContext("invalid")
		require.Error(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		ctx, cancel, err := durationStringToContext("1s")
		require.NoError(t, err)
		t.Cleanup(cancel)
		assert.NotEqual(t, context.Background(), ctx)
	})
}
