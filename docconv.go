package docconv // import "github.com/dev4mobile/mupdf/v2"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// Response payload sent back to the requestor
type Response struct {
	Body  string            `json:"body"`
	Meta  map[string]string `json:"meta"`
	MSecs uint32            `json:"msecs"`
	Error string            `json:"error"`
}

// MimeTypeByExtension returns a mimetype for the given extension, or
// application/octet-stream if none can be determined.
func MimeTypeByExtension(filename string) string {
	switch strings.ToLower(path.Ext(filename)) {
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".odt":
		return "application/vnd.oasis.opendocument.text"
	case ".pages":
		return "application/vnd.apple.pages"
	case ".pdf":
		return "application/pdf"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".rtf":
		return "application/rtf"
	case ".xml":
		return "text/xml"
	case ".xhtml", ".html", ".htm":
		return "text/html"
	case ".jpg", ".jpeg", ".jpe", ".jfif", ".jfif-tbnl":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".tif":
		return "image/tif"
	case ".tiff":
		return "image/tiff"
	case ".txt":
		return "text/plain"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".zip", ".7z", ".rar":
		return "application/zip"
	}
	return "application/octet-stream"
}

// Convert a file to plain text.
func Convert(r io.Reader, mimeType string, readability bool) (*Response, error) {
	if r == nil {
		return &Response{
			Error: "reader is nil", // 返回错误信息
		}, fmt.Errorf("reader is nil")
	}

	start := time.Now()

	var body string
	var meta map[string]string
	var err error
	switch mimeType {
	case "application/msword", "application/vnd.ms-word":
		slog.Warn("==>application/msword", "mimeType", mimeType)
		// body, meta, err = ConvertDoc(r)

	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		slog.Warn("==>application/vnd.openxmlformats-officedocument.wordprocessingml.document", "mimeType", mimeType)
		// body, meta, err = ConvertDocx(r)

	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		slog.Warn("==>application/vnd.openxmlformats-officedocument.presentationml.presentation", "mimeType", mimeType)
		body, meta, err = ConvertPptx(r)

	case "application/vnd.ms-excel":
		slog.Warn("==>application/vnd.ms-excel", "mimeType", mimeType)
		// body, meta, err = ConvertXls(r)

	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		slog.Warn("==>application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "mimeType", mimeType)
		body, meta, err = ConvertXlsx(r)

	case "application/vnd.oasis.opendocument.text":
		slog.Warn("==>application/vnd.oasis.opendocument.text", "mimeType", mimeType)
		body, meta, err = ConvertODT(r)

	case "application/vnd.apple.pages", "application/x-iwork-pages-sffpages":
		slog.Warn("==>application/vnd.apple.pages", "mimeType", mimeType)
		body, meta, err = ConvertPages(r)

	case "application/pdf":
		slog.Warn("==>application/pdf", "mimeType", mimeType)
		body, meta, err = ConvertPDF(r)

	case "application/rtf", "application/x-rtf", "text/rtf", "text/richtext":
		slog.Warn("==>application/rtf", "mimeType", mimeType)
		body, meta, err = ConvertRTF(r)

	case "text/html":
		slog.Warn("==>text/html", "mimeType", mimeType)
		body, meta, err = ConvertHTML(r, readability)

	case "text/url":
		slog.Warn("==>text/url", "mimeType", mimeType)
		body, meta, err = ConvertURL(r, readability)

	case "text/xml", "application/xml":
		slog.Warn("==>text/xml", "mimeType", mimeType)
		body, meta, err = ConvertXML(r)

	case "image/jpeg", "image/png", "image/tif", "image/tiff":
		slog.Warn("==>image/jpeg", "mimeType", mimeType)
		body, meta, err = ConvertImage(r)

	case "application/zip":
		slog.Warn("==>application/zip", "mimeType", mimeType)
		// body, meta, err = ConvertZip(r)

	case "text/plain":
		var b []byte
		b, err = io.ReadAll(r)
		slog.Warn("==>text/plain", "mimeType", mimeType, "size=", len(b))
		body = string(b)

	default:
		// auto-detection from first 512 bytes
		b, _ := io.ReadAll(r)
		if detect := http.DetectContentType(b); mimeType != detect {
			// recursive call convert once
			slog.Warn("==>detect:", "mimeType", detect, "detectMimeType", mimeType)
			return Convert(bytes.NewReader(b), detect, readability)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error converting data: %v", err)
	}

	return &Response{
		Body:  strings.TrimSpace(body),
		Meta:  meta,
		MSecs: uint32(time.Since(start) / time.Millisecond),
	}, nil
}

// ConvertPath converts a local path to text.
func ConvertPath(path string) (*Response, error) {
	mimeType := MimeTypeByExtension(path)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Convert(f, mimeType, true)
}

// ConvertPathReadability converts a local path to text, with the given readability
// option.
func ConvertPathReadability(path string, readability bool) ([]byte, error) {
	mimeType := MimeTypeByExtension(path)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := Convert(f, mimeType, readability)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}
