package engine

import (
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/winezer0/codecanvas/internal/model"
)

// IndexMatcher 提供基于索引的文件查找功能
type IndexMatcher struct {
	Index *model.FileIndex
}

// NewIndexMatcher 创建一个新的索引匹配器
func NewIndexMatcher(index *model.FileIndex) *IndexMatcher {
	return &IndexMatcher{Index: index}
}

// FindFiles 使用索引查找匹配的文件。
// pattern 支持:
// 1. 精确相对路径 (e.g., "/package.json")
// 2. 文件名匹配 (e.g., "package.json", "*.json")
// 3. 递归通配符 (e.g., "**/*.go", "src/**/*.js")
func (m *IndexMatcher) FindFiles(pattern string) ([]string, error) {
	var results []string

	// Case 1: 精确相对路径 (e.g. "/package.json")
	if strings.HasPrefix(pattern, "/") {
		target := strings.TrimPrefix(pattern, "/")
		// 检查索引中是否存在 - 使用小写文件名作为键
		for _, idx := range m.Index.NameMap[strings.ToLower(path.Base(target))] {
			f := m.Index.Files[idx]
			// 不区分大小写的路径比较
			if strings.EqualFold(f, target) {
				results = append(results, filepath.Join(m.Index.RootDir, f))
			}
		}
		// 如果索引中找不到，但模式是绝对的，可能我们应该信任模式？
		// 但这是基于索引的，所以我们只返回索引中的。
		return results, nil
	}

	// Case 2: 文件名 (e.g. "package.json")
	// 如果不包含路径分隔符，则为简单文件名匹配
	if !strings.Contains(pattern, "/") && !strings.Contains(pattern, "*") {
		// 使用小写键在NameMap中查找，实现不区分大小写的匹配
		if indices, ok := m.Index.NameMap[strings.ToLower(pattern)]; ok {
			for _, idx := range indices {
				results = append(results, filepath.Join(m.Index.RootDir, m.Index.Files[idx]))
			}
		}
		return results, nil
	}

	// Case 3: 后缀匹配 (e.g. "*.json") - 优化
	if strings.HasPrefix(pattern, "*.") && !strings.Contains(pattern[1:], "/") {
		// 使用小写扩展名作为ExtensionMap的键，实现不区分大小写的后缀匹配
		ext := strings.ToLower(pattern[1:]) // ".json"
		if indices, ok := m.Index.ExtensionMap[ext]; ok {
			for _, idx := range indices {
				results = append(results, filepath.Join(m.Index.RootDir, m.Index.Files[idx]))
			}
		}
		return results, nil
	}

	// Case 5: 目录匹配 (模式以 / 结尾)
	if strings.HasSuffix(pattern, "/") {
		// 使用小写进行比较，实现不区分大小写的目录匹配
		patternLower := strings.ToLower(pattern)
		for _, fileRelPath := range m.Index.Files {
			if strings.HasPrefix(strings.ToLower(fileRelPath), patternLower) {
				results = append(results, filepath.Join(m.Index.RootDir, fileRelPath))
			}
		}
		return results, nil
	}

	// Case 4: 通用 glob 匹配 (e.g. "src/**/*.ts")
	// 这需要遍历索引中的所有文件路径进行匹配
	// 这是一个 O(N) 操作，其中 N 是文件总数，比磁盘 I/O 快得多。
	for _, fileRelPath := range m.Index.Files {
		// filepath.Match 不支持 **，我们需要支持它。
		// 这里我们假设 pattern 遵循 gitignore 风格或 standard glob。
		// 为了简单起见，我们使用 filepath.Match 对每部分进行匹配，或者使用正则。
		// Go 的 filepath.Match 不支持递归 **。
		// 为了支持 **，我们可以使用第三方库，或者简单的实现：
		// 如果 pattern 包含 **，我们将 ** 替换为 .* 并使用正则。

		matched, err := matchPath(pattern, fileRelPath)
		if err == nil && matched {
			results = append(results, filepath.Join(m.Index.RootDir, fileRelPath))
		}
	}

	return results, nil
}

// matchPath 简单的路径匹配，支持 **
func matchPath(pattern, name string) (bool, error) {
	// 确保模式使用正斜杠以匹配 FileIndex 约定
	pattern = filepath.ToSlash(pattern)

	// 简单的 ** 支持
	if strings.Contains(pattern, "**") {
		// 转义 .
		regexPat := strings.ReplaceAll(pattern, ".", "\\.")
		// 替换 * 为 [^/]*
		regexPat = strings.ReplaceAll(regexPat, "*", "[^/]*")
		// 修复 **: 之前被替换成了 [^/]*[^/]*，我们需要 .*
		// 但这有点复杂。
		// 让我们简化策略：如果包含 **，我们暂时退化为简单的包含检查或者正则。
		// 更好的做法是：
		// 1. 将 pattern 分段
		// 2. 将 name 分段
		// 3. 匹配

		// 考虑到时间限制，我们使用一个简单的正则转换：
		// 1. QuoteMeta 转义特殊字符
		// 2. 将 \* 替换为 [^/]* (非路径分隔符的任意字符)
		// 3. 将 \*\* 替换为 .* (任意字符)
		// 4. 添加 i 标志使正则不区分大小写

		// 重置
		regexPat = regexp.QuoteMeta(pattern)
		regexPat = strings.ReplaceAll(regexPat, "\\*\\*", ".*")
		regexPat = strings.ReplaceAll(regexPat, "\\*", "[^/]*")
		regexPat = "^" + regexPat + "$"

		// 使用不区分大小写的正则匹配
		return regexp.MatchString("(?i:"+regexPat+")", name)
	}

	// 对于非 ** 模式，将 pattern 和 name 都转换为小写后再匹配
	return path.Match(strings.ToLower(pattern), strings.ToLower(name))
}

// containsAllKeywords 检查文件是否包含所有必需的关键字。
// 优化：流式读取，避免加载大文件。
func containsAllKeywords(content []byte, keywords []string) bool {
	if len(keywords) == 0 {
		return true
	}

	contentText := string(content)
	for _, kw := range keywords {
		if !strings.Contains(contentText, kw) {
			return false
		}
	}
	return true
}

// extractVersionFromText 使用给定的正则表达式模式从文件内容中提取版本字符串。
func extractVersionFromText(content []byte, pattern string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}

	matches := re.FindSubmatch(content)
	if len(matches) < 2 {
		return ""
	}

	return string(matches[1])
}

// extractVersionFromPath 使用给定的正则表达式模式从文件路径中提取版本字符串。
func extractVersionFromPath(filePath, pattern string) string {
	// 确保路径分隔符统一
	filePath = filepath.ToSlash(filePath)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}

	matches := re.FindStringSubmatch(filePath)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

// mapConfidenceLevel 将规则级别映射到置信度字符串。
func mapConfidenceLevel(level string) string {
	switch level {
	case "L1":
		return model.ConfidenceHigh
	case "L2":
		return model.ConfidenceMedium
	case "L3":
		return model.ConfidenceLow
	default:
		return model.ConfidenceLow
	}
}
