package output

import "github.com/ffuf/ffuf/v2/pkg/ffuf"

// ResultCollector manages the collection and lifecycle of scan results.
type ResultCollector struct {
	Results        []ffuf.Result
	CurrentResults []ffuf.Result
}

// NewResultCollector creates a new ResultCollector with empty result slices.
func NewResultCollector() *ResultCollector {
	return &ResultCollector{
		Results:        make([]ffuf.Result, 0),
		CurrentResults: make([]ffuf.Result, 0),
	}
}

// Reset clears the current results slice.
func (rc *ResultCollector) Reset() {
	rc.CurrentResults = make([]ffuf.Result, 0)
}

// Cycle moves CurrentResults into Results and resets CurrentResults.
func (rc *ResultCollector) Cycle() {
	rc.Results = append(rc.Results, rc.CurrentResults...)
	rc.Reset()
}

// GetCurrentResults returns the current results slice.
func (rc *ResultCollector) GetCurrentResults() []ffuf.Result {
	return rc.CurrentResults
}

// SetCurrentResults replaces the current results slice.
func (rc *ResultCollector) SetCurrentResults(results []ffuf.Result) {
	rc.CurrentResults = results
}
