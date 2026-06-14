package output

import (
	"encoding/csv"
	"os"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

func writeCSV(filename string, config *ffuf.Config, res []ffuf.Result, encode bool) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := GetKeywords(config)
	header = append(header, staticheaders...)

	if err := w.Write(header); err != nil {
		return err
	}
	for _, r := range res {
		if encode {
			inputs := make(map[string][]byte, len(r.Input))
			for k, v := range r.Input {
				inputs[k] = []byte(base64encode(v))
			}
			r.Input = inputs
		}

		err := w.Write(ToCsvRow(r))
		if err != nil {
			return err
		}
	}
	return nil
}
