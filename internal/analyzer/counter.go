package analyzer

import (
	"bufio"
	"os"
	"strings"
)

// FileStats 保存单个文件的行数统计
type FileStats struct {
	Code    int64
	Comment int64
	Blank   int64
	Lines   int64
}

// CountFile 分析文件并返回其统计信息
func CountFile(path string, lang *LanguageDefinition) (FileStats, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileStats{}, err
	}
	defer file.Close()

	stats := FileStats{}
	scanner := bufio.NewScanner(file)
	// 增加长行的缓冲区大小
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	inBlockComment := false
	var currentBlockEnd string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		stats.Lines++

		if line == "" {
			if inBlockComment {
				stats.Comment++
			} else {
				stats.Blank++
			}
			continue
		}

		if inBlockComment {
			stats.Comment++
			if strings.Contains(line, currentBlockEnd) {
				inBlockComment = false
				currentBlockEnd = ""
			}
			continue
		}

		// 检查块注释开始
		blockStarted := false
		for _, block := range lang.MultiLine {
			if strings.HasPrefix(line, block[0]) {
				inBlockComment = true
				currentBlockEnd = block[1]
				stats.Comment++
				blockStarted = true
				// 检查它是否也在同一行结束
				if strings.Contains(line[len(block[0]):], block[1]) {
					inBlockComment = false
					currentBlockEnd = ""
				}
				break
			}
		}
		if blockStarted {
			continue
		}

		// 检查行注释
		isLineComment := false
		for _, c := range lang.LineComments {
			if strings.HasPrefix(line, c) {
				stats.Comment++
				isLineComment = true
				break
			}
		}
		if isLineComment {
			continue
		}

		// 如果不是空行，也不是注释，就是代码
		stats.Code++
	}

	return stats, scanner.Err()
}
