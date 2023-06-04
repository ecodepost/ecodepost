package ssoservice

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"ecodepost/bff/pkg/invoker"
	"github.com/go-redis/redis/v8"

	"github.com/ego-component/eredis"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/util/xtime"
	"github.com/vmihailenco/msgpack"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type code struct {
	phoneCodeLength                int
	phoneCodeExpired               time.Duration
	phoneRegisterPrefixKey         string
	phoneLoginPrefixKey            string
	phoneRetrievePasswordPrefixKey string
	phoneMaxVerifyTimes            int
	phoneIpLimitKey                string
	phoneIpLimitRateNum            int64
	redis                          *eredis.Component
}

func InitCode(redis *eredis.Component) *code {
	return &code{
		redis:                          redis,
		phoneMaxVerifyTimes:            5,
		phoneCodeLength:                4,
		phoneIpLimitRateNum:            5, // 最大调用5次
		phoneIpLimitKey:                "sso:code:ipLimit:ip_%s_%d",
		phoneCodeExpired:               xtime.Duration("60s"),
		phoneRegisterPrefixKey:         "sso:code:register:phone_%s",
		phoneLoginPrefixKey:            "sso:code:login:phone_%s",
		phoneRetrievePasswordPrefixKey: "sso:code:retrieve_pwd:phone_%s",
	}
}

type CodeType string

const (
	RegisterCodeType         CodeType = "register"
	LoginCodeType            CodeType = "login"
	RetrievePasswordCodeType CodeType = "retrieve_password"
)

type codeStore struct {
	Code        string `msgpack:"c"`
	VerifyTimes int    `msgpack:"v"`
}

func (u codeStore) Marshal() ([]byte, error) {
	return msgpack.Marshal(u)
}

func (u *codeStore) Unmarshal(content []byte) error {
	return msgpack.Unmarshal(content, u)
}

// devCode 调试环境静态code
const devCode = "1111"

// SendCode 注册码
func (c *code) SendCode(ctx context.Context, phone string, ip string, codeType CodeType) (code string, codeTTL int64, isSent bool, err error) {
	var redisKey string
	switch codeType {
	case RegisterCodeType:
		redisKey = c.phoneRegisterPrefixKey
	case LoginCodeType:
		redisKey = c.phoneLoginPrefixKey
	case RetrievePasswordCodeType:
		redisKey = c.phoneRetrievePasswordPrefixKey
	}

	rateNum, err := c.setLimitRateByIp(ctx, ip)
	if err != nil {
		return "", 0, false, fmt.Errorf("setLimitRateByIp failed,err: %w", err)
	}
	if rateNum > c.phoneIpLimitRateNum {
		return "", 0, false, fmt.Errorf("rate limit over")
	}
	codeTTL, err = c.CodeTTLTime(ctx, phone, codeType)
	if err != nil {
		return "", 0, false, fmt.Errorf("CodeTTLTime failed,err: %w", err)
	}

	// 说明已经发过了，请等待
	if codeTTL != 0 {
		return "", codeTTL, true, nil
	}

	// 只有线上环境，才真正生成随机code，其他环境均使用devCode
	if econf.GetString("mode") == "pro" {
		code = genValidateCode(c.phoneCodeLength)
	} else {
		code = devCode
	}

	info, err := codeStore{Code: code, VerifyTimes: 0}.Marshal()
	if err != nil {
		return "", 0, false, fmt.Errorf("msgpack marshal failed,err: %w", err)
	}

	err = invoker.Redis.Set(ctx, fmt.Sprintf(redisKey, phone), info, c.phoneCodeExpired)
	if err != nil {
		return "", 0, false, fmt.Errorf("SendCode set redis failed, err: %w", err)
	}
	codeTTL = int64(c.phoneCodeExpired.Seconds())
	return
}

// CodeTTLTime 拿到该手机号的过期时间
func (c *code) CodeTTLTime(ctx context.Context, phone string, codeType CodeType) (ttlTime int64, err error) {
	var redisKey string
	switch codeType {
	case RegisterCodeType:
		redisKey = c.phoneRegisterPrefixKey
	case LoginCodeType:
		redisKey = c.phoneLoginPrefixKey
	case RetrievePasswordCodeType:
		redisKey = c.phoneRetrievePasswordPrefixKey
	}
	key := fmt.Sprintf(redisKey, phone)
	ttlInfo, err := invoker.Redis.TTL(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("VerifyCode get redis failed, err: %w", err)
	}

	// 说明没有这个数据
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return int64(ttlInfo.Seconds()), nil
}

func (c *code) VerifyCode(ctx context.Context, phone string, msgCode string, codeType CodeType) (showTips string, err error) {
	var redisKey string
	switch codeType {
	case RegisterCodeType:
		redisKey = c.phoneRegisterPrefixKey
	case LoginCodeType:
		redisKey = c.phoneLoginPrefixKey
	case RetrievePasswordCodeType:
		redisKey = c.phoneRetrievePasswordPrefixKey
	}
	key := fmt.Sprintf(redisKey, phone)
	codeInfo, err := invoker.Redis.GetBytes(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "不存在短信", fmt.Errorf("VerifyCode get redis failed, err: %w", err)
		}
		return "验证短信失败", fmt.Errorf("VerifyCode get redis failed, err: %w", err)
	}
	storeInfo := &codeStore{}
	if err = storeInfo.Unmarshal(codeInfo); err != nil {
		return "验证短信失败", fmt.Errorf("VerifyCode unmarshal failed, err: %w", err)
	}

	if storeInfo.VerifyTimes >= c.phoneMaxVerifyTimes {
		if _, err = invoker.Redis.Del(ctx, key); err != nil {
			return "", fmt.Errorf("VerifyCode delete failed, err: %w", err)
		}
		return "验证短信失败", fmt.Errorf("VerifyCode max verify times")
	}

	// 如果没验证成功
	if msgCode != storeInfo.Code {
		var info []byte
		storeInfo.VerifyTimes++
		info, err = storeInfo.Marshal()
		if err != nil {
			return "验证短信失败", fmt.Errorf("VerifyCode msgpack marshal failed,err: %w", err)
		}
		if err = invoker.Redis.Set(ctx, fmt.Sprintf(redisKey, phone), info, c.phoneCodeExpired); err != nil {
			return "验证短信失败", fmt.Errorf("VerifyCode set redis failed, err: %w", err)
		}
		return "验证短信失败", fmt.Errorf("VerifyCode not equal")
	}
	// 验证成功，删除redis信息
	_, _ = invoker.Redis.Del(ctx, key)
	return "", nil
}

// setLimitRateByIp 根据应用ID获取限制数
func (c *code) setLimitRateByIp(ctx context.Context, ip string) (num int64, err error) {
	limitKey := c.LimitRateKey(ip)
	num, err = c.redis.Incr(ctx, limitKey)
	if err != nil {
		err = fmt.Errorf("setLimitRateByIp failed1, %w", err)
		return
	}
	_, err = c.redis.Expire(ctx, limitKey, xtime.Duration("100s"))
	if err != nil {
		err = fmt.Errorf("setLimitRateByIp failed2, %w", err)
		return
	}
	return
}

func (c *code) LimitRateKey(ip string) string {
	return fmt.Sprintf(c.phoneIpLimitKey, ip, getNowMinTimestamps()) // ip_minuteTime
}

// genValidateCode 生成code码
func genValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func getNowMinTimestamps() int64 {
	local, _ := time.LoadLocation("Local")
	str := time.Now().Format("20060102 15:04:00")
	tt, _ := time.ParseInLocation("20060102 15:04:00", str, local)
	ts := tt.Unix()
	return ts
}
