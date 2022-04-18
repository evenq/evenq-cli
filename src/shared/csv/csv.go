package csv

import (
	"compress/gzip"
	"encoding/csv"
	"os"
	"strings"
)

func ReadHeaders(path string) ([]string, error) {
	var headers []string
	var r *csv.Reader

	file, err := os.Open(path)
	if err != nil {
		return headers, err
	}

	defer file.Close()

	// support gzipped files
	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return headers, err
		}
		defer gr.Close()

		r = csv.NewReader(gr)
	} else {
		r = csv.NewReader(file)
	}

	record, err := r.Read()

	return record, err
}
