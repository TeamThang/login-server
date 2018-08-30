package util

import (
	"fmt"
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
func GetUUID() (string, error) {
	u1, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("uuid gerator wrong: %s", err)
	}
	return u1.String(), nil
}
