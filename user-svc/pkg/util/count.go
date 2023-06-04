package util

import (
	"strconv"

	commonv1 "ecodepost/pb/common/v1"
)

type Biz interface {
	commonv1.CMN_BIZ | int32
}

type Act interface {
	commonv1.CNT_ACT | int32
}

// BizActKey BizActKey编码业务类型.动作类型的key，形如 1_1
func BizActKey[T1 Biz, T2 Act](biz T1, act T2) string {
	return strconv.Itoa(int(biz)) + "_" + strconv.Itoa(int(act))
}
