// @Author TrandLiu
// @Date 2024/12/15 0:53:00
// @Desc
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写加密
func Md5Encode(data string) string {

	h := md5.New()
	h.Write([]byte(data))
	temStr := h.Sum(nil)
	return hex.EncodeToString(temStr)
}

// 大写加密
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt string, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
