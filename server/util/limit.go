package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type limitConfig struct {
	Snippets_Length int `mapstructure:"snippets_length"`
	Snippets_Count  int `mapstructure:"snippets_count"`
	Images_Size     int `mapstructure:"images_size"`
	Images_Count    int `mapstructure:"images_count"`
}

var LimitConfig limitConfig

func (lc *limitConfig) SnippetsLength() int {
	return lc.Snippets_Length
}

func (lc *limitConfig) SnippetsCount() int {
	return lc.Snippets_Count
}

func (lc *limitConfig) ImagesSize() int {
	return lc.Images_Size
}

func (lc *limitConfig) ImagesCount() int {
	return lc.Images_Count
}

// InitializeLimits 从 Viper 加载 limit 配置
// 这个函数应该在 LoadConfig 之后被调用
func InitializeLimits() {
	newViper := viper.Sub("limit")
	if newViper == nil {
		// 不再 Fatal，改为 Errorf，让调用者决定如何处理
		// 或者可以考虑提供默认值
		log.Errorf("viper sub for 'limit' failed, LimitConfig might be zero")
		return
	}
	err := newViper.Unmarshal(&LimitConfig)
	if err != nil {
		log.Errorf("viper unmarshal for 'limit' failed: %+v, LimitConfig might be zero", err)
		return
	}
	log.Info("Limit configuration initialized successfully")
}
