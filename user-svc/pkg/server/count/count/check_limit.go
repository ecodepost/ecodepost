package count

import (
	"context"
	"fmt"
	"time"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
)

// CheckLimit 检测是否受到限制 响应true时需要限制
func (eng *Engine) CheckLimit(check string) bool {
	reply, err := eng.checkCache.Incr(context.Background(), check)
	if err != nil {
		elog.Error("Engine.CheckLimit.INCR", elog.String("check", check), elog.FieldErr(err))
		return false
	}

	// 限制大于0时才走限制规则
	if eng.LimitNum > 0 && reply > eng.LimitNum {
		elog.Info("Engine.CheckLimit.limit.success", elog.String("check", check))
		return true
	}
	_, err = eng.checkCache.Expire(context.Background(), check, 86400*time.Second)
	if err != nil {
		elog.Error("Engine.CheckLimit.EXPIRE", elog.String("check", check), elog.FieldErr(err))
		return false
	}
	return false
}

// CheckMQMsgID 检查消息是否存在
func (eng *Engine) CheckMQMsgID(msgID, checkType string) (isExists bool) {
	// isExists, err := redis.Bool(eng.checkCache.DoOnSlave("SISMEMBER", getMQMsgKey(checkType), msgID))
	isExists, err := eng.checkCache.SIsMember(context.Background(), getMQMsgKey(checkType), msgID)
	if err != nil {
		elog.Error("CheckMQMsgID", elog.FieldErr(err), elog.String("msgId", msgID))
	}

	return
}

// SetMQMsgID 添加成功消费消息
func (eng *Engine) SetMQMsgID(msgID, checkType string) (isExists bool) {
	_, err := eng.checkCache.SAdd(context.Background(), getMQMsgKey(checkType), msgID)
	if err != nil {
		elog.Error("SetMQMsgID SADD", elog.FieldErr(err), elog.String("msgId", msgID))
		return false
	}

	// _, err = eng.checkCache.DoOnMaster("EXPIRE", getMQMsgKey(checkType), 86400)
	_, err = eng.checkCache.Expire(context.Background(), getMQMsgKey(checkType), 86400*time.Second)
	if err != nil {
		elog.Error("SetMQMsgID EXPIRE", elog.FieldErr(err), elog.String("msgId", msgID))
		return false
	}
	return
}

// getMQMsgKey 消息处理key
func getMQMsgKey(checkType string) string {
	day := time.Now().Format("20060102")
	return fmt.Sprintf("ct_msg_%s:%s:%s", eapp.AppInstance(), checkType, day)
}
