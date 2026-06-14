package output

import (
	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// StatusColor returns the appropriate color code/string for a given HTTP status code.
// For terminal output, pass ANSI color constants; for HTML output, pass HTML hex color strings.
// This unifies the color mapping logic previously duplicated in Stdoutput.colorize and colorizeResults.
func StatusColor(status int64, color2xx, color3xx, color4xx, color5xx, defaultColor string) string {
	if status >= 200 && status <= 299 {
		return color2xx
	}
	if status >= 300 && status <= 399 {
		return color3xx
	}
	if status >= 400 && status <= 499 {
		return color4xx
	}
	if status >= 500 && status <= 599 {
		return color5xx
	}
	return defaultColor
}

// StatusHTMLColor returns the HTML hex color for a given HTTP status code.
func StatusHTMLColor(status int64) string {
	return StatusColor(status, "#adea9e", "#bbbbe6", "#d2cb7e", "#de8dc1", "black")
}

// StatusANSIColor returns the ANSI terminal color code for a given HTTP status code.
func StatusANSIColor(status int64) string {
	return StatusColor(status, ANSI_GREEN, ANSI_BLUE, ANSI_YELLOW, ANSI_RED, ANSI_CLEAR)
}

// ColorizeResults returns a new slice of results with HTMLColor set based on status code.
func ColorizeResults(results []ffuf.Result) []ffuf.Result {
	newResults := make([]ffuf.Result, 0, len(results))
	for _, r := range results {
		result := r
		result.HTMLColor = StatusHTMLColor(result.StatusCode)
		newResults = append(newResults, result)
	}
	return newResults
}
