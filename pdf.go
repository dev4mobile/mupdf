//go:build !ocr

package docconv

import (
	"fmt"
	"io"
)

func ConvertPDF(r io.Reader) (string, map[string]string, error) {
	f, err := NewLocalFile(r)
	if err != nil {
		return "", nil, fmt.Errorf("error creating local file: %v", err)
	}
	defer f.Done()

	bodyResult, metaResult, convertErr := ConvertPDFText(f.Name())
	if convertErr != nil {
		return "", nil, convertErr
	}

	return bodyResult.body, metaResult.meta, nil

}
