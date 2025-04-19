package util

import (
	"os"
	"path/filepath"
)

// DeleteImage 删除图片文件
// image 参数是图片的URL路径，例如: "/uploads/123456_abcd.jpg"
func DeleteImage(image string) error {
	// 如果路径为空，直接返回
	if image == "" {
		return nil
	}

	// 从URL路径中提取文件名
	filename := filepath.Base(image)

	// 使用配置的上传目录构建完整的文件路径
	filePath := GetImagePath(filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在则直接返回
	}

	// 删除文件
	return os.Remove(filePath)
}
