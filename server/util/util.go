package util

import (
	"log"
	"math/rand"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/crypto/bcrypt"
)

// 定义用于生成随机字符串的常量
const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 每个字符占用的比特数
	letterIdxMask = 1<<letterIdxBits - 1 // 用于位运算的掩码
	letterIdxMax  = 63 / letterIdxBits   // 每个 63 位随机数可生成的字符数
)

// 全局随机数生成器及互斥锁
var (
	rng   = rand.New(rand.NewSource(time.Now().UnixNano())) // 全局随机数生成器
	rngMu sync.Mutex                                        // 保护随机数生成器的互斥锁
)

// RandString 返回一个指定长度的随机字符串
func RandString(n int) string {
	b := make([]byte, n) // 用于存储生成的随机字符

	// 使用互斥锁保护随机数生成器
	rngMu.Lock()
	defer rngMu.Unlock()

	// 填充字节切片
	for i, cache, remain := n-1, rng.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rng.Int63(), letterIdxMax // 当缓存用尽时，重新生成随机数
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx] // 从字符集中选择对应的字符
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	// 将字节切片转换为字符串并返回
	return *(*string)(unsafe.Pointer(&b))
}

// String2bcryt 返回输入字符串的 bcrypt 哈希值
func String2bcrypt(str string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Encryption failure: %+v", err)
		return ""
	}
	return string(hashedPassword)
}
