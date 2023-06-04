package cron

import (
	"context"
	"math"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"github.com/ego-component/eredis/ecronlock"
	"github.com/gotomicro/ego/task/ecron"
)

var (
	locker *ecronlock.Component
)

// CronJob 异常任务
// https://blog.csdn.net/grafx/article/details/70495369
func CronJob() ecron.Ecron {
	job := func(ctx context.Context) error {
		var maxId int64
		invoker.Db.Model(mysql.File{}).Select("max(id) as maxid").Find(&maxId)
		for i := 0; i < int(maxId); i++ {
			articleId := i + 1
			var info mysql.File
			err := invoker.Db.Select("id,cnt_view,cnt_collect,cnt_comment,cnt_like,ctime,utime").Where("id = ?", articleId).Find(&info).Error
			if err != nil {
				continue
			}
			if info.Id == 0 {
				continue
			}

			score := calcTopicHotScore(info.CntView, info.CntCollect, info.CntComment, info.CntLike, info.Ctime, info.Utime)
			invoker.Db.Model(mysql.File{}).Where("id =?", articleId).Updates(map[string]interface{}{
				"recommend_score": score,
				"hot_score":       info.CntView + info.CntCollect + info.CntComment + info.CntLike,
			})
		}
		return nil
	}
	locker = ecronlock.DefaultContainer().Build(ecronlock.WithClient(invoker.Redis))
	cron := ecron.Load("cron.score").Build(
		ecron.WithJob(job),
		// 设置分布式锁
		ecron.WithLock(locker.NewLock("resource:cronjob:score")),
	)
	return cron
}

// calcTopicHotScore 计算topic的热度得分
func calcTopicHotScore(cntView, cntCollect, cntComment, cntLike, ctime, utime int64) float64 {
	// 新文章不能让分数为0，需要加个初始化积分
	// log10，需要+1，防止0对数为无穷小
	cnt := math.Log10(float64(cntView+1))*4 + float64(cntCollect) + float64(cntComment) + float64(cntLike)/5 + 1
	hour := math.Pow(float64((time.Now().Unix()-ctime)/60+1)-float64(utime/60-ctime/60), 1.5)
	if hour == 0 {
		hour = 1
	}
	return cnt / hour
}
