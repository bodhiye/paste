package util

import (
	"os"
	"path/filepath"
	"sync"
)

// 保护图片操作的互斥锁
var imageMutex sync.Mutex

// DeleteImage 删除图片文件
// image 参数是图片的URL路径，例如: "/uploads/123456_abcd.jpg"
func DeleteImage(image string) error {
	// 如果路径为空，直接返回
	if image == "" {
		return nil
	}

	// 使用互斥锁保护文件操作
	imageMutex.Lock()
	defer imageMutex.Unlock()

	// 从URL路径中提取文件名
	filename := filepath.Base(image)

	// 使用配置的上传目录构建完整的文件路径
	filePath := GetImagePath(filename)

	// 直接尝试删除，不先检查是否存在
	// 如果文件不存在，os.Remove会返回os.ErrNotExist
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		// 只有非"文件不存在"的错误才返回
		return err
	}

	return nil
}
