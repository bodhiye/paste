package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type limitClient struct {
	snippetLength int   `yaml:"snippet_length"` // 代码片段长度限制
	snippetCount  int   `yaml:"snippet_count"`  // 代码片段数量限制
	imageSize     int64 `yaml:"image_size"`     // 图片大小限制
	imageCount    int   `yaml:"image_count"`    // 图片数量限制
}

var LimitConfig limitClient

func (lc *limitClient) SnippetLength() int {
	return lc.snippetLength
}

func (lc *limitClient) SnippetCount() int {
	return lc.snippetCount
}

func (lc *limitClient) ImageSize() int64 {
	return lc.imageSize
}

func (lc *limitClient) ImageCount() int {
	return lc.imageCount
}

// InitLimitConfig 初始化限制配置
func InitLimitConfig() {
	LimitConfig.snippetCount = viper.GetInt("limit.snippet_count")
	LimitConfig.snippetLength = viper.GetInt("limit.snippet_length")
	LimitConfig.imageSize = viper.GetInt64("limit.image_size")
	LimitConfig.imageCount = viper.GetInt("limit.image_count")

	// 打印配置值以进行调试
	log.Infof("Limit config loaded: %+v", LimitConfig)
}
