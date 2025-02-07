//go:build ocr

package docconv

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/gen2brain/go-fitz"
)

var exts = []string{".jpg", ".tif", ".tiff", ".png", ".pbm"}

func compareExt(ext string, exts []string) bool {
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}

func ConvertPDFImages(path string) (BodyResult, error) {
	bodyResult := BodyResult{}

	// 打开PDF文档
	doc, err := fitz.New(path)
	if err != nil {
		return bodyResult, fmt.Errorf("error opening PDF: %v", err)
	}
	defer doc.Close()

	var wg sync.WaitGroup
	pageCount := doc.NumPage()
	data := make(chan string, pageCount)
	wg.Add(pageCount)

	// 遍历每一页提取图片
	for i := 0; i < pageCount; i++ {
		go func(pageNum int) {
			defer wg.Done()

			// 获取页面图片
			img, err := doc.Image(pageNum)
			if err != nil {
				return
			}

			// 转换图片为文本
			out, _, err := ConvertImage(img)
			if err != nil {
				return
			}

			data <- out
		}(i)
	}

	wg.Wait()
	close(data)

	for str := range data {
		bodyResult.body += str + " "
	}

	return bodyResult, nil
}

// PdfHasImage verify if `path` (PDF) has images
func PDFHasImage(path string) (bool, error) {
	doc, err := fitz.New(path)
	if err != nil {
		return false, fmt.Errorf("error opening PDF: %v", err)
	}
	defer doc.Close()

	// 检查前5页
	maxPages := 5
	if doc.NumPage() < maxPages {
		maxPages = doc.NumPage()
	}

	for i := 0; i < maxPages; i++ {
		images, err := doc.PageImageList(i)
		if err != nil {
			continue
		}
		if len(images) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func ConvertPDF(r io.Reader) (string, map[string]string, error) {
	f, err := NewLocalFile(r)
	if err != nil {
		return "", nil, fmt.Errorf("error creating local file: %v", err)
	}
	defer f.Done()

	// 打开PDF文档
	doc, err := fitz.New(f.Name())
	if err != nil {
		return "", nil, fmt.Errorf("error opening PDF: %v", err)
	}
	defer doc.Close()

	// 提取文本
	var bodyText strings.Builder
	for n := 0; n < doc.NumPage(); n++ {
		text, err := doc.Text(n)
		if err != nil {
			continue
		}
		bodyText.WriteString(text)
		bodyText.WriteString(" ")
	}

	// 检查是否包含图片
	hasImage, err := PDFHasImage(f.Name())
	if err != nil {
		return bodyText.String(), nil, fmt.Errorf("could not check if PDF has image: %w", err)
	}

	if !hasImage {
		return bodyText.String(), nil, nil
	}

	// 处理图片
	imageConvertResult, imageConvertErr := ConvertPDFImages(f.Name())
	if imageConvertErr != nil {
		return bodyText.String(), nil, nil // ignore error, return what we have
	}
	if imageConvertResult.err != nil {
		return bodyText.String(), nil, nil // ignore error, return what we have
	}

	fullBody := strings.Join([]string{bodyText.String(), imageConvertResult.body}, " ")
	return fullBody, nil, nil
}

var shellEscapePattern *regexp.Regexp

func init() {
	shellEscapePattern = regexp.MustCompile(`[^\w@%+=:,./-]`)
}

// shellEscape returns a shell-escaped version of the string s. The returned value
// is a string that can safely be used as one token in a shell command line.
func shellEscape(s string) string {
	if len(s) == 0 {
		return "''"
	}
	if shellEscapePattern.MatchString(s) {
		return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
	}

	return s
}
