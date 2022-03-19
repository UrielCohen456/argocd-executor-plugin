package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgocdPlugin(t *testing.T) {
	t.Run("Should return handler func", func (t *testing.T) {
		got := ArgocdPlugin()
		assert.NotNil(t, got)
	})	
}