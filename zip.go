package docconv

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/gen2brain/go-unarr"
)

// ConvertZip converts an archive file to text.
func ConvertZip(r io.Reader) (string, map[string]string, error) {
	a, err := unarr.NewArchiveFromReader(r)
	if err != nil {
		return "", nil, err
	}
	defer a.Close()

	var text string
	var meta map[string]string
	// iterate files and extract text
	for e := a.Entry(); e == nil; e = a.Entry() {
		if data, err := a.ReadAll(); err == nil {
			slog.Warn("convert zip", "name", a.Name(), "size", len(data), "mime", MimeTypeByExtension(a.Name()))
			if res, err := Convert(bytes.NewReader(data), MimeTypeByExtension(a.Name()), false); err == nil {
				text += a.Name() + "\r\n" + res.Body + "\r\n"
				meta = res.Meta
			} else {
				slog.Warn("convert zip", "name", a.Name(), "size", len(data), "mime", MimeTypeByExtension(a.Name()), "error", err)
			}
		}
	}
	slog.Warn("convert zip", "text", text, "meta", meta)
	return text, meta, nil
}
