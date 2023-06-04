package cron

import (
	"context"
	"fmt"

	enotify2 "ecodepost/sdk/enotify"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"github.com/ego-component/egorm"
	"github.com/ego-component/eredis"
	"github.com/ego-component/eredis/ecronlock"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ecron"
)

var (
	notifyComp *enotify2.Component
)

func NotifyCron() ecron.Ecron {
	// todo sync.Once
	invoker.Redis = eredis.Load("redis").Build()
	invoker.Db = egorm.Load("mysql").Build()
	notifyComp = enotify2.DefaultContainer().Build(enotify2.WithDb(invoker.Db))
	notifyComp.Start()
	// 构造分布式任务锁，目前已实现redis版本. 如果希望自定义，可以实现 ecron.Lock 接口
	locker := ecronlock.DefaultContainer().Build(ecronlock.WithClient(invoker.Redis))
	cron := ecron.Load("user-svc.cron.notify").Build(
		// 设置分布式锁
		ecron.WithLock(locker.NewLock("user:notify:cron")),
		ecron.WithJob(notify),
	)
	return cron
}

func notify(ctx context.Context) error {
	list, err := mysql.NotifyInitList(invoker.Db.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("notify msg receiver init list fail, err: %w", err)
	}
	for _, value := range list {
		err := notifyComp.Send(*value)
		if err != nil {
			elog.Error("notify send fail", elog.FieldErr(err))
			continue
		}
	}
	return nil
}
