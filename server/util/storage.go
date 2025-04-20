package util

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// StorageConfig 存储相关配置
type StorageConfig struct {
	uploadDir string
	uRLPrefix string
	host      string
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
		host:      viper.GetString("server.host"),
	}

	// 确保上传目录存在
	if err := os.MkdirAll(storageConfig.uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
}

// GetServerHost 尝试替换原本从配置文件中设置的 服务器地址的端口号
func GetServerHost() string {
	// 获取环境变量 PORT_HTTP 的值，通常用于指定服务器监听的端口号
	p := os.Getenv("PORT_HTTP")
	// 使用 net.SplitHostPort 函数将传入的地址 sh 分割为主机名（h）和端口号（_）
	h, _, err := net.SplitHostPort(storageConfig.host)
	if err != nil || len(p) == 0 {
		return storageConfig.host
	}
	return fmt.Sprintf("%s:%s", h, p)
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
	// 0.0.0.0 是一个特殊的IP地址，表示监听所有网络接口，但客户端无法直接使用它来访问
	// 这里我们检查host中是否包含0.0.0.0，如果包含则替换为localhost
	host := strings.Replace(storageConfig.host, "0.0.0.0", "localhost", 1)
	return filepath.Join("http://", host, storageConfig.uRLPrefix, filename)
}

// GetImagePath 根据文件名获取完整的文件路径
func GetImagePath(filename string) string {
	storageMutex.RLock()
	defer storageMutex.RUnlock()
	return filepath.Join(storageConfig.uploadDir, filename)
}
