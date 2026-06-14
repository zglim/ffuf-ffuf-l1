package output

import (
	"encoding/base64"
	"html"
	"strconv"
	"time"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// JsonResult is the JSON representation of a result for file output.
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

// ToJsonResult converts a single ffuf.Result to a JsonResult.
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

// ToJsonResults converts a slice of ffuf.Result to a slice of JsonResult.
func ToJsonResults(results []ffuf.Result) []JsonResult {
	jsonResults := make([]JsonResult, 0, len(results))
	for _, r := range results {
		jsonResults = append(jsonResults, ToJsonResult(r))
	}
	return jsonResults
}

// ToCsvRow converts a ffuf.Result to a CSV row (slice of strings).
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

// toCSV is kept for backward compatibility with existing tests.
func toCSV(r ffuf.Result) []string {
	return ToCsvRow(r)
}

// htmlResult represents a result for HTML and Markdown template output.
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

// ToHtmlResult converts a single ffuf.Result to an htmlResult.
// It handles strinput conversion, ffufhash extraction, scraper serialization, and HTML color assignment.
func ToHtmlResult(r ffuf.Result) htmlResult {
	ffufhash := ""
	strinput := make(map[string]string)
	for k, v := range r.Input {
		if k == "FFUFHASH" {
			ffufhash = string(v)
		} else {
			strinput[k] = string(v)
		}
	}
	strscraper := FormatScraperData(r.ScraperData)
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
		HTMLColor:        StatusHTMLColor(r.StatusCode),
		FfufHash:         ffufhash,
	}
}

// ToHtmlResults converts a slice of ffuf.Result to a slice of htmlResult.
func ToHtmlResults(results []ffuf.Result) []htmlResult {
	htmlResults := make([]htmlResult, 0, len(results))
	for _, r := range results {
		htmlResults = append(htmlResults, ToHtmlResult(r))
	}
	return htmlResults
}

// FormatScraperData serializes scraper data map into an HTML-formatted string.
func FormatScraperData(data map[string][]string) string {
	strscraper := ""
	for k, v := range data {
		if len(v) > 0 {
			strscraper = strscraper + "<p><b>" + html.EscapeString(k) + ":</b><br />"
			firstval := true
			for _, val := range v {
				if !firstval {
					strscraper += "<br />"
				}
				strscraper += html.EscapeString(val)
				firstval = false
			}
			strscraper += "</p>"
		}
	}
	return strscraper
}

// htmlFileOutput is the data structure for HTML/Markdown template rendering.
type htmlFileOutput struct {
	CommandLine string
	Time        string
	Keys        []string
	Results     []htmlResult
}

// CSV static headers appended after input provider keywords.
var staticheaders = []string{"url", "redirectlocation", "position", "status_code", "content_length", "content_words", "content_lines", "content_type", "duration", "resultfile", "Ffufhash"}

func base64encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

// GetKeywords extracts keyword names from config input providers.
func GetKeywords(config *ffuf.Config) []string {
	keywords := make([]string, 0)
	for _, inputprovider := range config.InputProviders {
		keywords = append(keywords, inputprovider.Keyword)
	}
	return keywords
}
