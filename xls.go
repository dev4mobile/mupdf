package docconv

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/sergeilem/xls"
)

// ConvertXls converts an Excel xls file to text.
func ConvertXls(r io.Reader) (string, map[string]string, error) {
	// Convert io.Reader to io.ReadSeeker
	data, err := io.ReadAll(r)
	slog.Info("xls==>", "size=", len(data))
	if err != nil {
		return "", nil, err
	}
	reader := bytes.NewReader(data)

	xlFile, err := xls.OpenReader(reader, "utf-8")
	if err != nil || xlFile == nil {
		return "", nil, err
	}

	var buf bytes.Buffer
	for n := 0; n < xlFile.NumSheets(); n++ {
		if sheet := xlFile.GetSheet(n); sheet != nil {
			sheetTitle := fmt.Sprintf("Sheet \"%s\" (%d rows):\n", sheet.Name, sheet.MaxRow)
			buf.WriteString(sheetTitle)

			for m := 0; m <= int(sheet.MaxRow); m++ {
				row := sheet.Row(m)
				if row == nil {
					continue
				}

				var rowText []string
				for c := row.FirstCol(); c < row.LastCol(); c++ {
					if text := row.Col(c); text != "" {
						// Clean cell text
						text = strings.ReplaceAll(text, "\n", " ")
						text = strings.ReplaceAll(text, "\r", "")
						text = strings.TrimSpace(text)

						rowText = append(rowText, text)
					}
				}

				if len(rowText) > 0 {
					buf.WriteString(strings.Join(rowText, ", "))
					buf.WriteString("\n")
				}
			}
		}
	}

	return buf.String(), nil, nil
}
