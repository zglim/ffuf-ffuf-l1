package output

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

type ejsonFileOutput struct {
	CommandLine string        `json:"commandline"`
	Time        string        `json:"time"`
	Results     []ffuf.Result `json:"results"`
	Config      *ffuf.Config  `json:"config"`
}

type jsonFileOutput struct {
	CommandLine string       `json:"commandline"`
	Time        string       `json:"time"`
	Results     []JsonResult `json:"results"`
	Config      *ffuf.Config `json:"config"`
}

func writeEJSON(filename string, config *ffuf.Config, res []ffuf.Result) error {
	t := time.Now()
	outJSON := ejsonFileOutput{
		CommandLine: config.CommandLine,
		Time:        t.Format(time.RFC3339),
		Results:     res,
	}

	outBytes, err := json.Marshal(outJSON)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, outBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func writeJSON(filename string, config *ffuf.Config, res []ffuf.Result) error {
	t := time.Now()
	outJSON := jsonFileOutput{
		CommandLine: config.CommandLine,
		Time:        t.Format(time.RFC3339),
		Results:     ToJsonResults(res),
		Config:      config,
	}
	outBytes, err := json.Marshal(outJSON)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, outBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
