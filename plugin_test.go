package derp

import (
	"testing"

	"github.com/benpate/derp/plugins"
	"github.com/stretchr/testify/assert"
)

// countingPlugin records how many times Report is called, so tests can
// verify that registered plugins are actually invoked.
type countingPlugin struct {
	count int
}

func (plugin *countingPlugin) Report(error) {
	plugin.count++
}

func TestPlugins(t *testing.T) {

	// Plugins are initialized containing a single item: ConsolePlugin{}
	assert.Equal(t, 1, len(Plugins))

	// Test making the list empty
	Plugins.Clear()
	assert.Equal(t, 0, len(Plugins))

	// Test adding items to the list
	Plugins.Add(plugins.JSON{})
	Plugins.Add(plugins.JSON{})
	Plugins.Add(plugins.JSON{})
	assert.Equal(t, 3, len(Plugins))
}

// TestReporterList_Local verifies that Add/Clear mutate the receiver itself,
// not the package-global Plugins (the bug this method signature fixed).
func TestReporterList_Local(t *testing.T) {

	globalLen := len(Plugins)

	list := ReporterList{}
	assert.Equal(t, 0, len(list))

	// Add must grow the local list
	list.Add(plugins.JSON{})
	list.Add(plugins.JSON{})
	assert.Equal(t, 2, len(list))

	// ...and must NOT touch the global Plugins
	assert.Equal(t, globalLen, len(Plugins))

	// Clear must empty the local list, leaving the global untouched
	list.Clear()
	assert.Equal(t, 0, len(list))
	assert.Equal(t, globalLen, len(Plugins))
}

// TestReporterList_Report verifies that every registered reporter is invoked
// once per call to Report, and that a nil error is never reported.
func TestReporterList_Report(t *testing.T) {

	first := &countingPlugin{}
	second := &countingPlugin{}

	Plugins.Clear()
	Plugins.Add(first)
	Plugins.Add(second)

	Report(NotFound("location", "message"))
	Report(NotFound("location", "message"))

	assert.Equal(t, 2, first.count)
	assert.Equal(t, 2, second.count)

	// A nil error must NOT be reported to any plugin
	Report(nil)
	assert.Equal(t, 2, first.count)
	assert.Equal(t, 2, second.count)
}
