package docconv

import (
	"bytes"
	"fmt"
	"github.com/EndFirstCorp/doc2txt"
	"io"
	"time"

	"github.com/richardlehane/mscfb"
	"github.com/richardlehane/msoleps"
)

// ConvertDoc converts an MS Word .doc to text.
func ConvertDoc(r io.Reader) (string, map[string]string, error) {
	f, err := NewLocalFile(r)
	if err != nil {
		return "", nil, fmt.Errorf("error creating local file: %v", err)
	}
	defer f.Done()

	// Meta data
	meta := make(map[string]string)
	doc, err := mscfb.New(f)
	if err != nil {
		// TODO: Propagate error.
		return "", nil, fmt.Errorf("error read meta %v", err)
	}

	props := msoleps.New()
	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		if msoleps.IsMSOLEPS(entry.Initial) {
			if err := props.Reset(doc); err != nil {
				// TODO: Propagate error.
				break
			}

			for _, prop := range props.Property {
				meta[prop.Name] = prop.String()
			}
		}
	}

	const defaultTimeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

	// Convert parsed meta
	if tmp, ok := meta["LastSaveTime"]; ok {
		if t, err := time.Parse(defaultTimeFormat, tmp); err == nil {
			meta["ModifiedDate"] = fmt.Sprintf("%d", t.Unix())
		}
	}
	if tmp, ok := meta["CreateTime"]; ok {
		if t, err := time.Parse(defaultTimeFormat, tmp); err == nil {
			meta["CreatedDate"] = fmt.Sprintf("%d", t.Unix())
		}
	}

	// Body data
	buf, err := doc2txt.ParseDoc(f)
	if err != nil {
		return "", nil, fmt.Errorf("error read body file: %v", err)
	}

	return buf.(*bytes.Buffer).String(), meta, nil
}
