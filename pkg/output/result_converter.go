package output

import (
	"html"
	"strconv"
	"time"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// statusColorIndex maps an HTTP status code to a color index.
// 0=default, 1=2xx, 2=3xx, 3=4xx, 4=5xx
func statusColorIndex(status int64) int {
	switch {
	case status >= 200 && status < 300:
		return 1
	case status >= 300 && status < 400:
		return 2
	case status >= 400 && status < 500:
		return 3
	case status >= 500 && status < 600:
		return 4
	default:
		return 0
	}
}

var ansiColors = [5]string{ANSI_CLEAR, ANSI_GREEN, ANSI_BLUE, ANSI_YELLOW, ANSI_RED}
var htmlColors = [5]string{"black", "#adea9e", "#bbbbe6", "#d2cb7e", "#de8dc1"}

// StatusColorANSI returns the ANSI escape code for the given HTTP status code.
func StatusColorANSI(status int64) string {
	return ansiColors[statusColorIndex(status)]
}

// StatusColorHTML returns the HTML hex color for the given HTTP status code.
func StatusColorHTML(status int64) string {
	return htmlColors[statusColorIndex(status)]
}

// JsonResult is the structure used for JSON file output.
type JsonResult struct {
	Input            map[string]string   `json:"input"`
	Position         int                 `json:"position"`
	StatusCode       int64               `json:"status"`
	ContentLength    int64               `json:"length"`
	ContentWords     int64               `json:"words"`
	ContentLines     int64               `json:"lines"`
	ContentType      string              `json:"content-type"`
	RedirectLocation string              `json:"redirectlocation"`
	ScraperData      map[string][]string `json:"scraper"`
	Duration         time.Duration       `json:"duration"`
	ResultFile       string              `json:"resultfile"`
	Url              string              `json:"url"`
	Host             string              `json:"host"`
}

// ToJsonResult converts an ffuf.Result into a JsonResult.
func ToJsonResult(r ffuf.Result) JsonResult {
	strinput := make(map[string]string)
	for k, v := range r.Input {
		strinput[k] = string(v)
	}
	return JsonResult{
		Input:            strinput,
		Position:         r.Position,
		StatusCode:       r.StatusCode,
		ContentLength:    r.ContentLength,
		ContentWords:     r.ContentWords,
		ContentLines:     r.ContentLines,
		ContentType:      r.ContentType,
		RedirectLocation: r.RedirectLocation,
		ScraperData:      r.ScraperData,
		Duration:         r.Duration,
		ResultFile:       r.ResultFile,
		Url:              r.Url,
		Host:             r.Host,
	}
}

// htmlResult is the structure used for HTML and Markdown template output.
type htmlResult struct {
	Input            map[string]string
	Position         int
	StatusCode       int64
	ContentLength    int64
	ContentWords     int64
	ContentLines     int64
	ContentType      string
	RedirectLocation string
	ScraperData      string
	Duration         time.Duration
	ResultFile       string
	Url              string
	Host             string
	HTMLColor        string
	FfufHash         string
}

// htmlFileOutput is the top-level structure passed to HTML and Markdown templates.
type htmlFileOutput struct {
	CommandLine string
	Time        string
	Keys        []string
	Results     []htmlResult
}

// ToHtmlResult converts an ffuf.Result into an htmlResult.
// If escapeHTML is true, scraper data keys and values are HTML-escaped.
func ToHtmlResult(r ffuf.Result, escapeHTML bool) htmlResult {
	ffufhash := ""
	strinput := make(map[string]string)
	for k, v := range r.Input {
		if k == "FFUFHASH" {
			ffufhash = string(v)
		} else {
			strinput[k] = string(v)
		}
	}
	strscraper := ""
	for k, v := range r.ScraperData {
		if len(v) > 0 {
			key := k
			if escapeHTML {
				key = html.EscapeString(k)
			}
			strscraper = strscraper + "<p><b>" + key + ":</b><br />"
			firstval := true
			for _, val := range v {
				if !firstval {
					strscraper += "<br />"
				}
				if escapeHTML {
					strscraper += html.EscapeString(val)
				} else {
					strscraper += val
				}
				firstval = false
			}
			strscraper += "</p>"
		}
	}
	return htmlResult{
		Input:            strinput,
		Position:         r.Position,
		StatusCode:       r.StatusCode,
		ContentLength:    r.ContentLength,
		ContentWords:     r.ContentWords,
		ContentLines:     r.ContentLines,
		ContentType:      r.ContentType,
		RedirectLocation: r.RedirectLocation,
		ScraperData:      strscraper,
		Duration:         r.Duration,
		ResultFile:       r.ResultFile,
		Url:              r.Url,
		Host:             r.Host,
		HTMLColor:        StatusColorHTML(r.StatusCode),
		FfufHash:         ffufhash,
	}
}

// ToCsvRow converts an ffuf.Result into a CSV row of strings.
func ToCsvRow(r ffuf.Result) []string {
	res := make([]string, 0)
	ffufhash := ""
	for k, v := range r.Input {
		if k == "FFUFHASH" {
			ffufhash = string(v)
		} else {
			res = append(res, string(v))
		}
	}
	res = append(res, r.Url)
	res = append(res, r.RedirectLocation)
	res = append(res, strconv.Itoa(r.Position))
	res = append(res, strconv.FormatInt(r.StatusCode, 10))
	res = append(res, strconv.FormatInt(r.ContentLength, 10))
	res = append(res, strconv.FormatInt(r.ContentWords, 10))
	res = append(res, strconv.FormatInt(r.ContentLines, 10))
	res = append(res, r.ContentType)
	res = append(res, r.Duration.String())
	res = append(res, r.ResultFile)
	res = append(res, ffufhash)
	return res
}
