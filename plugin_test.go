package derp

import (
	"testing"

	"github.com/benpate/derp/plugins"
	"github.com/stretchr/testify/assert"
)

func TestPlugins(t *testing.T) {

	// Plugins are initialized containing a single item: ConsolePlugin{}
	assert.Equal(t, 1, len(Plugins))

	// Test making the list empty
	Plugins.Clear()
	assert.Equal(t, 0, len(Plugins))

	// Test adding items to the list
	Plugins.Add(plugins.Console{})
	Plugins.Add(plugins.Console{})
	Plugins.Add(plugins.Console{})
	assert.Equal(t, 3, len(Plugins))

}
