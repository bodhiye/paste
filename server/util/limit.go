package util

import (

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type limitConfig struct {
    codeLength int
    codeCount int
    imageSize int
    imageCount int
}

var LimitConfig limitConfig

func (lc *limitConfig) CodeLength() int {
    return lc.codeLength
}

func (lc *limitConfig) CodeCount() int {
    return lc.codeCount
}

func (lc *limitConfig) ImageSize() int {
    return lc.imageSize
}

func (lc *limitConfig) ImageCount() int {
    return lc.imageCount
}

func init(){
    newViper := viper.Sub("limit")
    err := newViper.Unmarshal(&LimitConfig)
    if err != nil {
        log.Fatalf("viper unmarshal failed: +%v",err)
        return
    }
}