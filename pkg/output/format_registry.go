package output

import (
	"fmt"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// FormatWriter is the function signature for output format writers.
type FormatWriter func(filename string, config *ffuf.Config, results []ffuf.Result) error

type formatEntry struct {
	writer FormatWriter
	ext    string
}

var formatRegistry = map[string]formatEntry{}

// RegisterFormat registers an output format writer with its name and file extension.
// New output formats can be added by calling this function without modifying SaveFile.
func RegisterFormat(name, ext string, writer FormatWriter) {
	formatRegistry[name] = formatEntry{writer: writer, ext: ext}
}

// GetFormatWriter returns the writer function for a named format.
func GetFormatWriter(name string) (FormatWriter, bool) {
	entry, ok := formatRegistry[name]
	if !ok {
		return nil, false
	}
	return entry.writer, true
}

// GetFormatExt returns the file extension for a named format.
func GetFormatExt(name string) (string, bool) {
	entry, ok := formatRegistry[name]
	if !ok {
		return "", false
	}
	return entry.ext, true
}

// GetAllFormatNames returns the names of all registered formats (excluding "all").
func GetAllFormatNames() []string {
	names := make([]string, 0, len(formatRegistry))
	for name := range formatRegistry {
		names = append(names, name)
	}
	return names
}

func init() {
	RegisterFormat("json", ".json", writeJSON)
	RegisterFormat("ejson", ".ejson", writeEJSON)
	RegisterFormat("html", ".html", func(filename string, config *ffuf.Config, results []ffuf.Result) error {
		return writeHTML(filename, config, results)
	})
	RegisterFormat("md", ".md", writeMarkdown)
	RegisterFormat("csv", ".csv", func(filename string, config *ffuf.Config, results []ffuf.Result) error {
		return writeCSV(filename, config, results, false)
	})
	RegisterFormat("ecsv", ".ecsv", func(filename string, config *ffuf.Config, results []ffuf.Result) error {
		return writeCSV(filename, config, results, true)
	})
}

// WriteAllFormats writes results in all registered formats, using the base filename
// with each format's extension appended.
func WriteAllFormats(baseFilename string, config *ffuf.Config, results []ffuf.Result, errorFn func(string)) error {
	for name, entry := range formatRegistry {
		outputFile := baseFilename + entry.ext
		err := entry.writer(outputFile, config, results)
		if err != nil {
			errorFn(fmt.Sprintf("%s: %s", name, err.Error()))
		}
	}
	return nil
}
