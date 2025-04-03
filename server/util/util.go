package util

import (
	"log"
	"math/rand"
	"time"
	"unsafe"

	"golang.org/x/crypto/bcrypt"
)

// 定义用于生成随机字符串的常量
const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6      // 每个字符占用的比特数
	letterIdxMask = 1<<letterIdxBits - 1  // 用于位运算的掩码
	letterIdxMax  = 63 / letterIdxBits  // 每个 63 位随机数可生成的字符数
)

// RandString 返回一个指定长度的随机字符串
func RandString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano()) // 使用当前时间的纳秒数作为随机数种子
	b := make([]byte, n)                            // 用于存储生成的随机字符

	// 填充字节切片
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax // 当缓存用尽时，重新生成随机数
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

// String2bcryt 返回输入字符串的 MD5 哈希值的十六进制表示
func String2bcrypt(str string) string {
	// h := md5.New()
	// h.Write([]byte(str))
	// return hex.EncodeToString(h.Sum(nil))
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(str),bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Encryption failure: %+v", err)
		return ""
	}
	return string(hashedPassword)
}