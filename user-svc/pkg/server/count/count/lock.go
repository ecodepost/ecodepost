package count

import (
	"context"
	"fmt"
	"time"

	"github.com/ego-component/eredis"
	"github.com/gotomicro/ego/core/elog"
)

// tryLock ...
func (eng *Engine) tryLock(key string, lockTime time.Duration) (*eredis.Lock, bool) {
	var lockKey = fmt.Sprintf("cnt_lock:%s", key)
	lock, err := eng.limitCache.LockClient().Obtain(context.Background(), lockKey, lockTime)
	if err != nil {
		elog.Error("tryLock.SET", elog.String("key", key), elog.String("lock_key", lockKey), elog.FieldErr(err))
		return lock, false
	}
	return lock, true
}

// unLock ...
func (eng *Engine) unLock(lock *eredis.Lock) {
	if err := lock.Release(context.Background()); err != nil {
		elog.Error("tryLock.DEL", elog.String("key", lock.Key()), elog.String("token", lock.Token()), elog.FieldErr(err))
	}
}
