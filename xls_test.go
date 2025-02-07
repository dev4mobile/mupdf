package docconv

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertXls(t *testing.T) {
	tests := []struct {
		file            string
		wantTrimmedText string
		wantMeta        map[string]string
		wantErr         bool
	}{
		{
			file: "002-test.xls",
			wantTrimmedText: `Sheet "Sheet1" (5 rows):
客户名称, 联系方式
张强, 13805313105
王磊, 13156016177
张三, 15098846688
李四, 15865263302
王武, 15806400671`,
			wantMeta: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			f, err := os.Open(path.Join("testdata", tt.file))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			gotText, gotMeta, err := ConvertXls(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotText = strings.TrimSpace(gotText)
			if gotText != tt.wantTrimmedText {
				t.Errorf("ConvertDoc() text = %v, want %v", gotText, tt.wantTrimmedText)
			}
			if !cmp.Equal(tt.wantMeta, gotMeta, maybeTimeComparer) {
				t.Errorf("ConvertDoc() meta mismatch (-want +got):\n%v", cmp.Diff(tt.wantMeta, gotMeta, maybeTimeComparer))
			}
		})
	}
}
