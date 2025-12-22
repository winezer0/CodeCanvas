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

// CountFileStats 分析文件并返回其统计信息
func CountFileStats(path string) (FileStats, error) {
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

	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		trimmedLine := strings.TrimSpace(line)

		// 处理空白行
		if trimmedLine == "" {
			if inBlockComment {
				stats.Comment++
			} else {
				stats.Blank++
			}
			continue
		}

		// 处理块注释
		if inBlockComment {
			stats.Comment++
			// 检查块注释结束
			if strings.Contains(line, "*/") {
				inBlockComment = false
			}
			continue
		}

		// 检查块注释开始
		if strings.Contains(line, "/*") {
			stats.Comment++
			// 检查块注释是否在同一行结束
			if !strings.Contains(line, "*/") {
				inBlockComment = true
			}
			continue
		}

		// 检查行注释
		if strings.HasPrefix(trimmedLine, "//") || strings.HasPrefix(trimmedLine, "#") {
			stats.Comment++
			continue
		}

		// 检查行内注释
		if strings.Contains(line, "//") || strings.Contains(line, "#") {
			stats.Code++
			continue
		}

		// 其他情况视为代码
		stats.Code++
	}

	return stats, scanner.Err()
}
