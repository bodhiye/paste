package util

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// StorageConfig 存储相关配置
type StorageConfig struct {
	UploadDir string
	URLPrefix string
}

var storageConfig StorageConfig

// InitializeStorage 初始化存储配置
func InitializeStorage() {
	storageConfig = StorageConfig{
		UploadDir: viper.GetString("paste.storage.upload_dir"),
		URLPrefix: viper.GetString("paste.storage.url_prefix"),
	}

	// 确保上传目录存在
	if err := os.MkdirAll(storageConfig.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
}

// GetUploadDir 获取上传目录
func GetUploadDir() string {
	return storageConfig.UploadDir
}

// GetURLPrefix 获取URL前缀
func GetURLPrefix() string {
	return storageConfig.URLPrefix
}

// GetImageURL 根据文件名获取完整的URL
func GetImageURL(filename string) string {
	return storageConfig.URLPrefix + filename
}

// GetImagePath 根据文件名获取完整的文件路径
func GetImagePath(filename string) string {
	return filepath.Join(storageConfig.UploadDir, filename)
}
