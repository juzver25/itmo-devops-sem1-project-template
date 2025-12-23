package handlers

import (
	"archive/zip"
	"bytes"
	"io"
)

func buildZipWithDataCSV(dataCSV []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	fw, err := zw.Create("data.csv")
	if err != nil {
		_ = zw.Close()
		return nil, err
	}

	if _, err := io.Copy(fw, bytes.NewReader(dataCSV)); err != nil {
		_ = zw.Close()
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
