package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/satori/go.uuid"
)

// 获取字符串的md5
func GetMD5(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 生成uuid
func GetUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}
