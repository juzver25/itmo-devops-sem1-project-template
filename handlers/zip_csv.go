package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

func getCSVFromZipBody(body []byte) (io.ReadCloser, error) {
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("zip: %w", err)
	}

	for _, f := range zr.File {
		if f.Name == "data.csv" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("zip open data.csv: %w", err)
			}
			return rc, nil
		}
	}

	return nil, fmt.Errorf("data.csv not found in zip")
}
