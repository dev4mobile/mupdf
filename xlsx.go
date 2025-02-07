package docconv

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// ConvertXlsx converts an MS Word xlsx file to text.
func ConvertXlsx(r io.Reader) (string, map[string]string, error) {
	var size int64

	// Common case: if the reader is a file (or trivial wrapper), avoid
	// loading it all into memory.
	var ra io.ReaderAt
	if f, ok := r.(interface {
		io.ReaderAt
		Stat() (os.FileInfo, error)
	}); ok {
		si, err := f.Stat()
		if err != nil {
			return "", nil, err
		}
		size = si.Size()
		ra = f
	} else {
		b, err := io.ReadAll(io.LimitReader(r, maxBytes))
		if err != nil {
			return "", nil, nil
		}
		size = int64(len(b))
		ra = bytes.NewReader(b)
	}

	zr, err := zip.NewReader(ra, size)
	if err != nil {
		return "", nil, fmt.Errorf("error unzipping data: %v", err)
	}
	var text string
	for _, file := range zr.File {
		rc, err := file.Open()
		if err != nil {
			continue
		}

		str, err := XlsxXMLToText(rc)
		if err != nil {
			continue
		}
		text += str + "\n"
		_ = rc.Close()
	}
	return text, nil, err
}

// XlsxXMLToText converts Docx XML into plain text.
func XlsxXMLToText(r io.Reader) (string, error) {
	return XMLToText(r, []string{"v", "t"}, []string{"instrText", "script"}, true)
}
