package count

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	countv1 "ecodepost/pb/count/v1"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"github.com/gotomicro/ego/core/elog"
	"gorm.io/gorm"
)

// IncrOnce 添加计数操作
func (eng *Engine) IncrOnce(db *gorm.DB, tid, fid string, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, num, realNum int64, ct, did, ip string) (err error) {
	info := &mysql.CounterInfo{
		Tid: tid,
		Fid: fid,
		Biz: biz,
		Act: act,
	}

	err = eng.db.IncrUpdateLog(db, num, realNum, info, ct, did, ip)
	if err != nil {
		elog.Error("OperateCounter.IncrOnce.IncrUpdateLog", elog.String("tid", info.Tid), elog.Any("info", info))
		return
	}

	log.Printf("22--------------->"+"%+v\n", 22)
	err = eng.db.IncrUpdateStat(db, num, realNum, info)
	if err != nil {
		elog.Error("OperateCounter.IncrOnce.IncrUpdateStat", elog.String("tid", info.Tid), elog.Any("info", info))
		return
	}
	return nil
}

// RepeatedUpdate 直接更新计数值可以重复计数
func (eng *Engine) RepeatedUpdate(ctx context.Context, tid, fid string, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, num, realNum int64, ct, did, ip string) (err error) {
	info := &mysql.CounterInfo{
		Tid:     tid,
		Fid:     fid,
		Biz:     biz,
		Act:     act,
		Num:     num,
		RealNum: realNum,
	}

	if err = eng.db.IncrUpdateLog(invoker.Db.WithContext(ctx), num, realNum, info, ct, did, ip); err != nil {
		elog.Error("OperateCounter.RepeatedUpdate.IncrUpdateLog", elog.String("tid", info.Tid), elog.Any("info", info))
		return
	}

	if err = eng.db.IncrUpdateStat(invoker.Db.WithContext(ctx), num, realNum, info); err != nil {
		elog.Error("OperateCounter.RepeatedUpdate.IncrUpdateStat", elog.String("tid", info.Tid), elog.Any("info", info))
		return
	}

	// 处理异常数据
	eng.tryCatchCount(ctx, num, realNum, info)
	// 处理随机事件 TODO 根据配置决定是否开启
	// if num > 0 && act == commonv1.CNT_ACT_VIEW {
	// 	eng.updateRandNumEvent(ctx, info, num)
	// }
	return nil
}

// updateRandNumEvent 更新需要随机加成的事件类型
func (eng *Engine) updateRandNumEvent(ctx context.Context, info *mysql.CounterInfo, num int64) {
	var err error
	// 默认第一条新增数据 则初始化数据
	var baseNum = eng.randBaseNumInt64()
	if num > 1 {
		baseInfo, err := eng.db.GetStatInfo(ctx, &mysql.CounterInfo{
			Tid: info.Tid,
			Biz: info.Biz,
			Act: info.Act,
		})
		if err != nil {
			elog.Error("OperateCounter.updateRandNumEvent.GetStatInfo", elog.String("tid", info.Tid), elog.Any("info", info))
		}
		baseNum = baseInfo.BaseNum
	}

	err = eng.db.SetUpdateStatFrontNum(ctx, &mysql.CounterInfo{
		Tid:     info.Tid,
		Fid:     info.Fid,
		Biz:     info.Biz,
		Act:     info.Act,
		Num:     eng.computeUpNum(baseNum, num),
		BaseNum: baseNum,
	})

	if err != nil {
		elog.Error("OperateCounter.updateRandNumEvent.SetUpdateStatFrontNum", elog.String("tid", info.Tid), elog.Any("info", info))
	}
}

// computeUpNum 计算浏览数加成值
func (eng *Engine) computeUpNum(baseNum, realUpNum int64) int64 {
	return baseNum + realUpNum*eng.multipleNum
}

// randBaseNumInt64 计算基础随机分
func (eng *Engine) randBaseNumInt64() int64 {
	if eng.baseMinNum >= eng.baseMaxNum || eng.baseMinNum == 0 || eng.baseMaxNum == 0 {
		return eng.baseMaxNum
	}
	return rand.Int63n(eng.baseMaxNum-eng.baseMinNum) + eng.baseMinNum
}

// InitOtsCountBaseNum 初始化更新ots计数值
func (eng *Engine) InitOtsCountBaseNum(ctx context.Context, info mysql.CounterInfo) (err error) {
	baseNum := eng.randBaseNumInt64()
	baseInfo, err := eng.db.GetStatInfo(ctx, &mysql.CounterInfo{Tid: info.Tid, Biz: info.Biz, Act: info.Act})
	if err != nil {
		return
	}
	// 检测数据是否存在初始化
	if baseInfo.Tid != "" {
		return
	}
	err = eng.db.StoreCountStat(ctx, &mysql.CounterInfo{
		Tid:     info.Tid,
		Biz:     info.Biz,
		Act:     info.Act,
		Num:     baseNum,
		BaseNum: baseNum,
		Utime:   time.Now().Unix(),
	})
	return
}

// tryCatchCount 处理异常、防止计数为负数的情况
func (eng *Engine) tryCatchCount(ctx context.Context, num, realNum int64, info *mysql.CounterInfo) {
	var err error
	// 处理异常情况
	if realNum <= 0 {
		err = eng.db.SetUpdateLogRealNum(invoker.Db.WithContext(ctx), &mysql.CounterInfo{
			Tid:     info.Tid,
			Fid:     info.Fid,
			Biz:     info.Biz,
			Act:     info.Act,
			RealNum: 0,
		})
		if err != nil {
			elog.Error("OperateCounter.tryCatchCount.SetUpdateLogRealNum", elog.String("tid", info.Tid), elog.Any("info", info))
		}
	}

	if num <= 0 {
		err = eng.db.SetUpdateStatRealNum(ctx, &mysql.CounterInfo{
			Tid:     info.Tid,
			Biz:     info.Biz,
			Act:     info.Act,
			RealNum: 0,
			Num:     0,
		})
		if err != nil {
			elog.Error("OperateCounter.tryCatchCount.SetUpdateStatRealNum", elog.String("tid", info.Tid), elog.Any("info", info))
		}
	}
}

// AddBehaviorLog 添加行为日志
func (eng *Engine) AddBehaviorLog(ctx context.Context, info *mysql.CountBehaviorLog) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := eng.db.CountBehaviorLog.StoreBehaviorLog(ctx, info)
	if err != nil {
		elog.Error("AddBehaviorLog.Create.Fail", elog.String("fid", info.Fid), elog.String("tid", info.Tid), elog.String("biz", info.Biz.String()), elog.String("acti", info.Acti.String()), elog.FieldErr(err))
		return
	}

	elog.Info("AddBehaviorLog.Create.Success", elog.String("fid", info.Fid), elog.String("tid", info.Tid), elog.String("biz", info.Biz.String()), elog.String("acti", info.Acti.String()))
}

// OtsBatchGetTargetCnt 前台批量获取目标计数值
// key:tid val:countv1.TargetActsCount
func (eng *Engine) OtsBatchGetTargetCnt(ctx context.Context, in *countv1.GetTnumByBaksAndTidsReq) (res map[string]*countv1.TargetActsCount, err error) {
	res = make(map[string]*countv1.TargetActsCount)
	var (
		eg sync.WaitGroup
		lc sync.RWMutex
	)
	for _, tid := range in.Tids {
		tid := tid
		go func() {
			eg.Add(1)
			defer eg.Done()
			cntMap, err := eng.db.ScanCountStatNumByActKeys(invoker.Db.WithContext(ctx), tid, in.Baks)
			// todo err
			if err != nil {
				elog.Error("ScanCountStatNum fail", elog.FieldErr(err))
				return
			}

			eventFrontNum := make(map[string]uint64)
			for event, cntNum := range cntMap {
				eventFrontNum[event] = uint64(cntNum.Num)
			}

			lc.Lock()
			res[tid] = &countv1.TargetActsCount{Map: eventFrontNum}
			lc.Unlock()
		}()
	}
	eg.Wait()
	return res, err
}
