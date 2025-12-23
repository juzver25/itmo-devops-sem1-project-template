package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
)

func getCSVFromZipBody(body []byte) (io.ReadCloser, error) {
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("zip: %w", err)
	}

	for _, f := range zr.File {
		if path.Base(f.Name) == "data.csv" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("zip open data.csv: %w", err)
			}
			return rc, nil
		}
	}

	for _, f := range zr.File {
		name := path.Base(f.Name)
		if strings.HasSuffix(strings.ToLower(name), ".csv") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("zip open csv: %w", err)
			}
			return rc, nil
		}
	}

	return nil, fmt.Errorf("no csv file found in zip")
}
