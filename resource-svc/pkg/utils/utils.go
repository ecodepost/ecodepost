package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gotomicro/ego/core/econf"
	"github.com/microcosm-cc/bluemonday"
)

// CheckFileName TODO validate
func CheckFileName(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("file name is empty")
	}
	if utf8.RuneCountInString(fileName) > 300 {
		return fmt.Errorf("file name is over size")
	}
	return nil
}

func UUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func StripHtmlTags(html string) string {
	return bluemonday.StripTagsPolicy().Sanitize(html)
}

func GetArticleMentionByGuid(content string) []string {
	ids := make([]string, 0)
	articleMap := make(map[string]struct{})
	var re *regexp.Regexp
	// TODO
	if econf.GetString("regex.host") == "forum.gocn.vip" {
		re = regexp.MustCompile(`https:\/\/forum\.gocn\.vip\/topics\/\d+`)
	}
	if econf.GetString("regex.host") == "gocn.vip" {
		re = regexp.MustCompile(`https:\/\/gocn\.vip\/topics\/\d+`)
	}

	slices := re.FindAllString(content, -1)
	for _, value := range slices {
		if value == "" {
			continue
		}
		arrs := strings.Split(value, "/")
		articleId := arrs[len(arrs)-1]
		if articleId == "" {
			continue
		}
		_, ok := articleMap[articleId]
		if ok {
			continue
		}
		ids = append(ids, articleId)
		// 去重复
		articleMap[articleId] = struct{}{}
	}
	return ids
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// GetRandomSalt return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(8)
}

// GetRandomString 生成随机字符串
func GetRandomString(num int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetFileSuffix(fn string) (suffix string) {
	splitRes := strings.Split(fn, `.`)
	if len(splitRes) > 1 {
		return splitRes[len(splitRes)-1]
	}

	return
}

func NewVal[T comparable](val T) *T {
	return &val
}

type Meta = map[string]any

func NewMeta(meta Meta) []byte {
	metaBytes, _ := json.Marshal(meta)
	return metaBytes
}
