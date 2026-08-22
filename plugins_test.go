package derp

import (
	"sync"
	"sync/atomic"
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
	assert.Equal(t, 1, Plugins.Len())

	// Test making the list empty
	Plugins.Clear()
	assert.Equal(t, 0, Plugins.Len())

	// Test adding items to the list
	Plugins.Add(plugins.JSON{})
	Plugins.Add(plugins.JSON{})
	Plugins.Add(plugins.JSON{})
	assert.Equal(t, 3, Plugins.Len())
}

// TestReporterList_Local verifies that Add/Clear mutate the receiver itself,
// not the package-global Plugins (the bug this method signature fixed).
func TestReporterList_Local(t *testing.T) {

	globalLen := Plugins.Len()

	// A zero-value ReporterList is usable
	var list ReporterList
	assert.Equal(t, 0, list.Len())

	// Add must grow the local list
	list.Add(plugins.JSON{})
	list.Add(plugins.JSON{})
	assert.Equal(t, 2, list.Len())

	// ...and must NOT touch the global Plugins
	assert.Equal(t, globalLen, Plugins.Len())

	// Clear must empty the local list, leaving the global untouched
	list.Clear()
	assert.Equal(t, 0, list.Len())
	assert.Equal(t, globalLen, Plugins.Len())
}

// TestReporterList_Set verifies that Set replaces the whole list in one step, and that later
// changes to the caller's slice do not leak into the published list.
func TestReporterList_Set(t *testing.T) {

	var list ReporterList

	first := &countingPlugin{}
	second := &countingPlugin{}

	// Set replaces whatever was there before
	list.Add(plugins.JSON{})
	list.Set(first, second)
	assert.Equal(t, 2, list.Len())

	// The published list is a defensive copy: mutating the caller's slice changes nothing
	reporters := []Reporter{first}
	list.Set(reporters...)
	reporters[0] = second
	assert.Equal(t, 1, list.Len())
	assert.Same(t, first, list.slice()[0])
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

// TestReporterList_ConcurrentSetAndReport is the reason this type holds an atomic pointer: a
// server that reloads its logging configuration calls Set from a reload goroutine while every
// other goroutine in the process may be inside Report.  Run with -race.  Against the old bare
// slice this was undefined behavior -- a torn slice header mid-range; against the atomic swap
// every Report sees one complete generation.
func TestReporterList_ConcurrentSetAndReport(t *testing.T) {

	var list ReporterList

	first := &atomicCountingPlugin{}
	second := &atomicCountingPlugin{}
	list.Set(first)

	var waitGroup sync.WaitGroup

	// The reloader: swap the whole list, repeatedly
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		for index := 0; index < 1000; index++ {
			if index%2 == 0 {
				list.Set(first, second)
			} else {
				list.Set(second)
			}
		}
	}()

	// The reporters: what every other goroutine in a process does
	err := NotFound("location", "message")

	for worker := 0; worker < 8; worker++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for index := 0; index < 1000; index++ {
				for _, reporter := range list.slice() {
					reporter.Report(err)
				}
			}
		}()
	}

	waitGroup.Wait()

	// Every Report reached SOME complete generation
	assert.Positive(t, first.count.Load()+second.count.Load())
}

// TestReporterList_ConcurrentAdds pins the copy-on-write in Add: concurrent Adds must not lose
// entries to a torn append.  Run with -race.
func TestReporterList_ConcurrentAdds(t *testing.T) {

	var list ReporterList

	var waitGroup sync.WaitGroup

	for worker := 0; worker < 8; worker++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for index := 0; index < 25; index++ {
				list.Add(&atomicCountingPlugin{})
			}
		}()
	}

	waitGroup.Wait()

	assert.Equal(t, 200, list.Len(), "no Add may be lost")
}

// atomicCountingPlugin counts reports without data races, so it can be shared across the
// goroutines of the concurrency tests above.
type atomicCountingPlugin struct {
	count atomic.Int64
}

// Report counts one reported error.
func (plugin *atomicCountingPlugin) Report(error) {
	plugin.count.Add(1)
}
