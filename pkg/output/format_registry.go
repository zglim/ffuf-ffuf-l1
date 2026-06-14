package output

import "github.com/ffuf/ffuf/v2/pkg/ffuf"

// FormatWriter is a function that writes results to a file in a specific format.
type FormatWriter func(filename string, config *ffuf.Config, res []ffuf.Result) error

// FormatRegistration holds the extension and writer function for an output format.
type FormatRegistration struct {
	Extension string
	Writer    FormatWriter
}

var formatRegistry = map[string]FormatRegistration{}

// RegisterFormat registers a new output format with its file extension and writer function.
func RegisterFormat(name, ext string, writer FormatWriter) {
	formatRegistry[name] = FormatRegistration{Extension: ext, Writer: writer}
}

func init() {
	RegisterFormat("json", ".json", writeJSON)
	RegisterFormat("ejson", ".ejson", writeEJSON)
	RegisterFormat("html", ".html", writeHTML)
	RegisterFormat("md", ".md", writeMarkdown)
	RegisterFormat("csv", ".csv", func(filename string, config *ffuf.Config, res []ffuf.Result) error {
		return writeCSV(filename, config, res, false)
	})
	RegisterFormat("ecsv", ".ecsv", func(filename string, config *ffuf.Config, res []ffuf.Result) error {
		return writeCSV(filename, config, res, true)
	})
}
