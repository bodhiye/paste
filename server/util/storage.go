package util

import (
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// StorageConfig 存储相关配置
type StorageConfig struct {
	uploadDir string
	uRLPrefix string
}

var (
	storageConfig StorageConfig
	storageMutex  sync.RWMutex // 用于保护文件操作的读写锁
)

// InitializeStorage 初始化存储配置
func InitializeStorage() {
	storageMutex.Lock()
	defer storageMutex.Unlock()

	storageConfig = StorageConfig{
		uploadDir: viper.GetString("paste.storage.upload_dir"),
		uRLPrefix: viper.GetString("paste.storage.url_prefix"),
	}

	// 确保上传目录存在
	if err := os.MkdirAll(storageConfig.uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
}

// GetUploadDir 获取上传目录
func GetUploadDir() string {
	storageMutex.RLock()
	defer storageMutex.RUnlock()
	return storageConfig.uploadDir
}

// GetURLPrefix 获取URL前缀
func GetURLPrefix() string {
	storageMutex.RLock()
	defer storageMutex.RUnlock()
	return storageConfig.uRLPrefix
}

// GetImageURL 根据文件名获取完整的URL
func GetImageURL(filename string) string {
	storageMutex.RLock()
	defer storageMutex.RUnlock()
	return storageConfig.uRLPrefix + filename
}

// GetImagePath 根据文件名获取完整的文件路径
func GetImagePath(filename string) string {
	storageMutex.RLock()
	defer storageMutex.RUnlock()
	return filepath.Join(storageConfig.uploadDir, filename)
}
