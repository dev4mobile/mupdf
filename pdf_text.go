package docconv

import (
	"fmt"
	"strings"

	"github.com/gen2brain/go-fitz"
)

// BodyResult 表示 PDF 内容
type BodyResult struct {
	body string
}

type MetaResult struct {
	meta map[string]string
	err  error
}

// ConvertPDFText 使用 go-fitz 库对 PDF 进行解析，提取所有页面的文本内容，并返回 PDF 内容和元数据。
func ConvertPDFText(path string) (BodyResult, MetaResult, error) {
	// 打开 PDF 文件
	doc, err := fitz.New(path)
	if err != nil {
		return BodyResult{}, MetaResult{}, fmt.Errorf("无法打开 PDF 文件: %v", err)
	}
	defer doc.Close()

	var builder strings.Builder
	// 遍历 PDF 所有页面，提取文本
	for n := 0; n < doc.NumPage(); n++ {
		pageText, err := doc.Text(n)
		if err != nil {
			return BodyResult{}, MetaResult{}, fmt.Errorf("无法从第 %d 页提取文本: %v", n, err)
		}
		builder.WriteString(pageText)
		builder.WriteString("\n")
	}

	// 构造元数据
	meta := MetaResult{
		meta: map[string]string{
			"num_pages": fmt.Sprintf("%d", doc.NumPage()),
		},
	}

	body := BodyResult{
		body: builder.String(),
	}

	return body, meta, nil
}
