package util

import (
	"crypto/md5"
	"fmt"
	"regexp"
)

func Md5HexStr(payload []byte) string {
	hasher := md5.New()
	hasher.Write(payload)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// 匹配规则
// ^1第一位为一
// [345789]{1} 后接一位345789 的数字
// \\d \d的转义 表示数字 {9} 接9位
// $ 结束符
var phoneRgx, _ = regexp.Compile(`^1[3456789]{1}\d{9}$`)

// CheckMobile 检验手机号
func CheckMobile(phone string) bool {
	// 返回 MatchString 是否匹配
	return phoneRgx.MatchString(phone)
}
