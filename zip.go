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
	count := 0
	for e := a.Entry(); e == nil; e = a.Entry() {
		if count > 10 {
			slog.Warn("==>5. too many files", "count", count, "name", a.Name())
			break
		}
		count++
		if data, err := a.ReadAll(); err == nil && len(data) < 5*1024*1024 {
			slog.Warn("==>1. convert zip", "name", a.Name(), "size", len(data), "mime", MimeTypeByExtension(a.Name()))
			if res, err := Convert(bytes.NewReader(data), MimeTypeByExtension(a.Name()), false); err == nil {
				slog.Warn("==>2. convert zip", "name", a.Name(), "size", len(data), "mime", MimeTypeByExtension(a.Name()), "error", err)
				if len(res.Body) <= 0 {
					continue
				}
				text += a.Name() + "\r\n" + res.Body + "\r\n"
				meta = res.Meta
			} else {
				slog.Warn("==>3. convert zip", "name", a.Name(), "size", len(data), "mime", MimeTypeByExtension(a.Name()), "error", err)
			}
		}
	}
	slog.Warn("==>4. convert zip", "text", text, "meta", meta)
	return text, meta, nil
}
