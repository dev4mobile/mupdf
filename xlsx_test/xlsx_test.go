package xlsx_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Servicewall/docconv/v2"
)

func TestConvertXlsx(t *testing.T) {
	f, err := os.Open("./testdata/Excel表格文件.xlsx")
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}

	resp, _, err := docconv.ConvertXlsx(f)
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}

	if want := "客户名称"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
	if want := "张强"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
	if want := "15806400671"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
}
