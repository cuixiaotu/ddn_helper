package behavior

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func append2File() error {
	// 创建 CSV 文件
	file, err := os.OpenFile("members.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	return nil
}

// 假设 formMemberHtml 解析 HTML，返回 [][]string
func DoFormMemberHtml(resp string) ([][]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	table := doc.Find("#tblData")
	var rows [][]string

	// 遍历表格行，跳过表头和“合计”行
	table.Find("tr").NextAll().Each(func(i int, tr *goquery.Selection) {
		if strings.Contains(tr.Text(), "合计") {
			return
		}

		var row []string
		tr.Find("td").Each(func(_ int, td *goquery.Selection) {
			text := strings.TrimSpace(td.Text())
			row = append(row, text)
		})

		if len(row) > 0 {
			rows = append(rows, row)
		}
	})

	return rows, nil
}

func appendToCSV(filePath string, rows [][]string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(rows)
}
