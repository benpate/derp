package derp

import (
	"sync"
	"sync/atomic"
)

// Reporter wraps the "Report" method, which reports a derp error to an external
// source. Reporters are responsible for handling and swallowing any errors they generate.
type Reporter interface {
	Report(error)
}

// ReporterList is a concurrency-safe registry of Reporters, which are called in succession
// whenever the Report() function is called.
//
// Readers never lock: Report loads the current list with a single atomic read, so error
// reporting stays wait-free no matter how often the list changes.  Writers publish a NEW list
// with one atomic store -- a published list is never modified in place -- which makes it safe
// to reconfigure reporters at runtime (for example, when a server reloads its logging
// configuration) while other goroutines are mid-Report.
//
// RULE: To replace the whole list at runtime, use Set -- it swaps once, so there is never a
// moment with no reporters.  A Clear-then-Add sequence is also safe, but it opens a brief
// window where errors reported by other goroutines reach an empty list and vanish.
type ReporterList struct {
	lock      sync.Mutex                 // serializes writers; readers never take it
	reporters atomic.Pointer[[]Reporter] // the current, immutable list
}

// Set replaces the entire list with the provided reporters in a single atomic swap.
func (list *ReporterList) Set(reporters ...Reporter) {

	// Cloned defensively: the caller may keep (and mutate) its own slice, and a published
	// list must never change underneath a concurrent Report.
	value := make([]Reporter, len(reporters))
	copy(value, reporters)

	list.reporters.Store(&value)
}

// Add appends a new reporter to this list.  This lets the developer configure
// and append additional reporters during initialization.
func (list *ReporterList) Add(reporter Reporter) {

	// RULE: Copy-on-write under the writer lock.  Appending to the published slice in place
	// could grow a backing array while a concurrent Report is ranging it.
	list.lock.Lock()
	defer list.lock.Unlock()

	current := list.slice()

	value := make([]Reporter, len(current), len(current)+1)
	copy(value, current)
	value = append(value, reporter)

	list.reporters.Store(&value)
}

// Clear removes all reporters from this list.  It is useful for removing the library
// default JSON reporter from the list, in the event that you don't want to report
// errors to the console.
func (list *ReporterList) Clear() {
	list.reporters.Store(&[]Reporter{})
}

// Len returns the number of reporters currently in this list.
func (list *ReporterList) Len() int {
	return len(list.slice())
}

// slice returns the current, immutable list of reporters.  Callers may range it freely --
// it will never change -- but must NEVER write into it.
func (list *ReporterList) slice() []Reporter {

	if value := list.reporters.Load(); value != nil {
		return *value
	}

	// A zero-value ReporterList is usable, and empty.  This returns an empty (not nil)
	// slice so that the return value is always safe to range, index-check, and re-slice.
	return []Reporter{}
}
