package output

import (
	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// ResultCollector manages the lifecycle and collection of results.
// It is embedded in Stdoutput to separate collection concerns from terminal output formatting.
type ResultCollector struct {
	Results        []ffuf.Result
	CurrentResults []ffuf.Result
}

// NewResultCollector creates a new empty ResultCollector.
func NewResultCollector() *ResultCollector {
	return &ResultCollector{
		Results:        make([]ffuf.Result, 0),
		CurrentResults: make([]ffuf.Result, 0),
	}
}

// Reset resets the current results slice.
func (rc *ResultCollector) Reset() {
	rc.CurrentResults = make([]ffuf.Result, 0)
}

// Cycle moves CurrentResults to Results and resets the current results slice.
func (rc *ResultCollector) Cycle() {
	rc.Results = append(rc.Results, rc.CurrentResults...)
	rc.Reset()
}

// GetCurrentResults returns the current results slice.
func (rc *ResultCollector) GetCurrentResults() []ffuf.Result {
	return rc.CurrentResults
}

// SetCurrentResults sets the current results slice.
func (rc *ResultCollector) SetCurrentResults(results []ffuf.Result) {
	rc.CurrentResults = results
}

// AddResult appends a single result to the current results.
func (rc *ResultCollector) AddResult(r ffuf.Result) {
	rc.CurrentResults = append(rc.CurrentResults, r)
}

// AllResults returns the combined Results and CurrentResults.
func (rc *ResultCollector) AllResults() []ffuf.Result {
	return append(rc.Results, rc.CurrentResults...)
}
