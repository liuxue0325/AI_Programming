package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetRelativePath 获取相对路径
func GetRelativePath(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

// GetAbsolutePath 获取绝对路径
func GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	execDir, err := os.Getwd()
	if err != nil {
		return path
	}

	return filepath.Join(execDir, path)
}

// IsHiddenFile 判断是否为隐藏文件
func IsHiddenFile(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".")
}

// GetFileExt 获取文件扩展名
func GetFileExt(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	return ext
}

// GetFileNameWithoutExt 获取不含扩展名的文件名
func GetFileNameWithoutExt(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
