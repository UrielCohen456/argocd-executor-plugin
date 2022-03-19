package argocd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgocdPlugin(t *testing.T) {
	t.Run("Should return handler func", func (t *testing.T) {
		got := ArgocdPlugin(context.Background())
		assert.NotNil(t, got)
	})	
}