package util

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type limitConfig struct {
	mu              sync.RWMutex // 保护配置读写的互斥锁
	Snippets_Length int          `mapstructure:"snippets_length"`
	Snippets_Count  int          `mapstructure:"snippets_count"`
	Images_Size     int          `mapstructure:"images_size"`
	Images_Count    int          `mapstructure:"images_count"`
}

var LimitConfig limitConfig

func (lc *limitConfig) SnippetsLength() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.Snippets_Length
}

func (lc *limitConfig) SnippetsCount() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.Snippets_Count
}

func (lc *limitConfig) ImagesSize() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.Images_Size
}

func (lc *limitConfig) ImagesCount() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.Images_Count
}

// InitializeLimits 从 Viper 加载 limit 配置
func InitializeLimits() {
	LimitConfig.mu.Lock()
	defer LimitConfig.mu.Unlock()

	newViper := viper.Sub("limit")
	if newViper == nil {
		log.Errorf("viper sub for 'limit' failed, using default values")
		setDefaultLimits()
		return
	}

	err := newViper.Unmarshal(&LimitConfig)
	if err != nil {
		log.Errorf("viper unmarshal for 'limit' failed: %+v, using default values", err)
		setDefaultLimits()
		return
	}

	// 即使成功解析了配置，也要检查是否需要设置默认值
	ensureValidLimits()

	log.Info("Limit configuration initialized successfully")
}

// 设置默认限制值
func setDefaultLimits() {
	LimitConfig.Snippets_Length = 30000
	LimitConfig.Snippets_Count = 10
	LimitConfig.Images_Size = 10
	LimitConfig.Images_Count = 3
}

// 确保所有限制值都是有效的
func ensureValidLimits() {
	if LimitConfig.Snippets_Length <= 0 {
		LimitConfig.Snippets_Length = 30000
	}
	if LimitConfig.Snippets_Count <= 0 {
		LimitConfig.Snippets_Count = 10
	}
	if LimitConfig.Images_Size <= 0 {
		LimitConfig.Images_Size = 10
	}
	if LimitConfig.Images_Count <= 0 {
		LimitConfig.Images_Count = 3
	}
}
