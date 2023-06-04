package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base32"
	"encoding/base64"
	"time"
)

type S2s = map[string]string
type S2i = map[string]int
type S2a = map[string]interface{}

func Aes256CfbEncrypt(secret string, plaintext string) (ciphertext string, err error) {
	if secret == "" {
		return
	}
	secret = RepeatStrToLen(secret, 32)

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return
	}
	iv := []byte(secret[:block.BlockSize()])

	dst := make([]byte, len(plaintext))
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(dst, []byte(plaintext))

	ciphertext = base64.URLEncoding.WithPadding(base32.NoPadding).EncodeToString(dst)
	return
}

func Aes256CfbDecrypt(secret string, ciphertext string) (plaintext string, err error) {
	if secret == "" {
		return
	}
	secret = RepeatStrToLen(secret, 32)
	cipherTextBytes, err := base64.URLEncoding.WithPadding(base32.NoPadding).DecodeString(ciphertext)
	if err != nil {
		return
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return
	}
	iv := []byte(secret[:block.BlockSize()])

	dst := make([]byte, len(cipherTextBytes))
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(dst, cipherTextBytes)

	plaintext = string(dst)
	return
}

func RepeatStrToLen(str string, length int) string {
	if str == "" {
		panic("RepeatStrToLen: empty str")
	}

	for len(str) < length {
		str += str
	}

	str = str[:length]
	return str
}

// ToDate 获取日期
func ToDate(now time.Time) string {
	return now.Format("2006-01-02")
}

func FromBase64(config string) (result string, err error) {
	bytes, err := base64.StdEncoding.DecodeString(config)
	if err != nil {
		return
	}
	result = string(bytes)
	return
}

func ToBase64(config string) (result string) {
	return base64.StdEncoding.EncodeToString([]byte(config))
}
