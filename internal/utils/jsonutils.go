package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ToJson 将任意 map 转换为格式化的 JSON 字符串（用于输出）
func ToJson(v interface{}) string {
	return string(ToJsonBytes(v))
}

// ToJsonBytes  将任意 map 转换为格式化的 JSON 字符串（用于输出）
func ToJsonBytes(v interface{}) []byte {
	data, _ := json.MarshalIndent(v, "", "  ")
	return data
}

// EnsureDir 确保目录存在，如果不存在则创建。
// 如果 isFile 为 true，则 dirPath 被视为文件路径，函数会确保其所在目录存在；
// 如果 isFile 为 false，则 dirPath 被视为目录路径，函数会确保该目录存在。
func EnsureDir(dirPath string, isFile bool) error {
	targetDir := dirPath
	if isFile {
		targetDir = filepath.Dir(dirPath)
	}
	return os.MkdirAll(targetDir, 0755)
}
func WriteJSON(filePath string, v interface{}) error {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("serialization of JSON data failed: %v", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write the JSON file: %v", err)
	}

	return nil
}
