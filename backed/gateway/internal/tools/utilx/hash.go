package utilx

import (
	"math/rand"
	"time"
)

const (
	charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// GenerateRandomString 生成包含数字或字母的随机字符串
func GenerateRandomString(length int) string {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
