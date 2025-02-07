package zip_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Servicewall/docconv/v2"
)

func TestConvertZip(t *testing.T) {
	f, err := os.Open("./testdata/docx.zip")
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}
	resp, _, err := docconv.ConvertZip(f)
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}

	if want := "sample.docx"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
	if want := "sample_2.docx"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
	if want := "sample_3.docx"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
}

func TestConvertRar(t *testing.T) {
	f, err := os.Open("./testdata/test.rar")
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}
	resp, _, err := docconv.ConvertZip(f)
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}

	if want := "13805313105"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
}

func TestConvert7z(t *testing.T) {
	f, err := os.Open("./testdata/test.7z")
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}
	resp, _, err := docconv.ConvertZip(f)
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}

	if want := "13805313105"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
}

func TestConvertZipInZip(t *testing.T) {
	f, err := os.Open("./testdata/zip.zip")
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}
	resp, _, err := docconv.ConvertZip(f)
	if err != nil {
		t.Fatalf("got error = %v, want nil", err)
	}

	if want := "13805313105"; !strings.Contains(resp, want) {
		t.Errorf("expected %v to contains %v", resp, want)
	}
}
