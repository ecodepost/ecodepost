package count

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	countv1 "ecodepost/pb/count/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	"ecodepost/user-svc/pkg/code"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"ecodepost/user-svc/pkg/server/count/cache"
	"ecodepost/user-svc/pkg/util"
	"github.com/gotomicro/ego/core/elog"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// LockKey Redis分布式锁Key
func LockKey(biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, tid string) string {
	return fmt.Sprintf("%s_%s", util.BizActKey(biz, act), tid)
}

// checkLimit 检测是否被限制拦截住
func (eng *Engine) checkLimit(in *countv1.SetReq) bool {
	switch in.Biz {
	case commonv1.CMN_BIZ_ARTICLE, commonv1.CMN_BIZ_ANSWER, commonv1.CMN_BIZ_QUESTION:
		day := time.Now().Format("20060102")
		if in.GetDid() == "" || in.Tid == "" {
			return true
		}
		key := util.BizActKey(in.Biz, in.Act) + "_" + day + "_" + in.GetDid() + ":" + in.Tid
		return eng.CheckLimit(key)
	default:
	}
	return false
}

// checkBidParams 业务传参检测
func (eng *Engine) checkBidParams(bidType commonv1.CMN_BIZ, act commonv1.CNT_ACT) bool {
	if bidType.String() == "" || act.String() == "" {
		return false
	}

	return true
}

// checkOpType 检测操作类型是否合法
func (eng *Engine) checkOpType(actType commonv1.CNT_ACT, acti commonv1.CNT_ACTI) bool {
	switch actType {
	case commonv1.CNT_ACT_COLLECT, commonv1.CNT_ACT_LIKE, commonv1.CNT_ACT_DISLIKE, commonv1.CNT_ACT_FOLLOW:
		if acti == commonv1.CNT_ACTI_ADD || acti == commonv1.CNT_ACTI_SUB {
			return true
		}
	case commonv1.CNT_ACT_VIEW, commonv1.CNT_ACT_CNT:
		if acti == commonv1.CNT_ACTI_ADD || acti == commonv1.CNT_ACTI_SUB || acti == commonv1.CNT_ACTI_UPDATE || acti == commonv1.CNT_ACTI_RESET {
			return true
		}
	default:
		elog.Error("unsupported act")
	}
	return false
}

func (eng *Engine) Set(ctx context.Context, in *countv1.SetReq) (res *countv1.SetRes, err error) {
	res = &countv1.SetRes{}
	// 计数时增加频率限制拦截
	if eng.checkLimit(in) {
		return nil, code.RateLimtError
	}

	if !eng.checkConfig(in.Biz, in.Act, in.Acti) {
		return nil, code.ConfigNotExists
	}

	lock, succ := eng.tryLock(LockKey(in.Biz, in.Act, in.Tid), 1*time.Second)
	if in.MaxVal > 0 && !succ {
		elog.Warn("SetCount.tryLock.fail", elog.String("fid", in.Fid), elog.String("tid", in.Tid), elog.String("biz", in.Biz.String()), zap.Any("params", in))
		return nil, code.ServerBusyError
	}

	if in.MaxVal > 0 {
		defer eng.unLock(lock)
	}

	if in.MaxVal > 0 {
		info, err := eng.db.GetStatInfo(ctx, &mysql.CounterInfo{Tid: in.Tid, Fid: in.Fid, Biz: in.Biz, Act: in.Act})
		if err != nil {
			elog.Error("SetCount.GetStatInfo.fail", elog.String("fid", in.Fid), elog.String("tid", in.Tid), elog.String("biz", in.Biz.String()), zap.Any("params", in))
			return nil, errcodev1.ErrDbError().WithMessage("GetStatInfo fail," + err.Error())
		}
		if info.RealNum >= int64(in.MaxVal) {
			elog.Info("SetCount.GetStatInfo.limit", elog.String("fid", in.Fid), elog.String("tid", in.Tid), elog.String("biz", in.Biz.String()), zap.Any("params", in))
			return nil, code.CountLimitError
		}
	}

	return eng.DealCountEvent(ctx, in)
}

func (eng *Engine) GetTdetailsByTids(ctx context.Context, in *countv1.GetTdetailsByTidsReq) (res *countv1.GetTdetailsByTidsRes, err error) {
	if !eng.checkBidParams(in.Biz, in.Act) {
		return nil, code.BidTypeParamsError
	}
	if len(in.Tids) == 0 {
		return nil, code.TidZero
	}
	if len(in.Tids) > 100 {
		return nil, code.TidIsLimit
	}
	if in.MaxFids > cache.TargetMaxFromNum {
		return nil, code.TidIsLimit
	}

	data, err := eng.countCache.GetMap(ctx, in)
	if err != nil {
		return nil, code.RedisOpError
	}
	return &countv1.GetTdetailsByTidsRes{Map: data}, nil
}

func (eng *Engine) GetTdetailsByFid(ctx context.Context, in *countv1.GetTdetailsByFidReq) (res *countv1.GetTdetailsByFidRes, err error) {
	if !eng.checkBidParams(in.Biz, in.Act) {
		return nil, code.BidTypeParamsError
	}

	if in.GetLimit() > 30 {
		return nil, code.TidIsLimit
	}

	data, err := eng.countCache.GetTidByFid(ctx, in)
	if err != nil {
		return nil, code.RedisOpError
	}

	return &countv1.GetTdetailsByFidRes{List: data}, nil
}

func (eng *Engine) GetFidsTdetailByTid(ctx context.Context, in *countv1.GetFidsTdetailByTidReq) (res *countv1.GetFidsTdetailByTidRes, err error) {
	if !eng.checkBidParams(in.Biz, in.Act) {
		return nil, code.BidTypeParamsError
	}

	if in.GetLimit() > cache.TargetMaxFromNum {
		return nil, code.TidIsLimit
	}

	data, err := eng.countCache.GetFids(ctx, in)
	if err != nil {
		return nil, code.RedisOpError
	}
	return data, nil
}

func (eng *Engine) GetTnumByBaksAndTids(ctx context.Context, in *countv1.GetTnumByBaksAndTidsReq) (res *countv1.GetTnumByBaksAndTidsRes, err error) {
	if len(in.Baks) == 0 {
		return nil, code.BidTypeLengthError
	}

	if len(in.Tids) == 0 {
		return nil, code.BidTypeLengthError
	}

	data, err := eng.OtsBatchGetTargetCnt(ctx, in)
	if err != nil {
		return nil, code.RedisOpError
	}

	return &countv1.GetTnumByBaksAndTidsRes{Map: data}, nil
}

func (eng *Engine) DBGetTdetailsByFid(ctx context.Context, in *countv1.DBGetTdetailsByFidReq) (*countv1.DBGetTdetailsByFidRes, error) {
	if !eng.checkBidParams(in.Biz, in.Act) {
		return nil, code.BidTypeParamsError
	}

	if in.Pagination.PageSize > 30 {
		return nil, code.TidIsLimit
	}

	// 查询countLog记录
	cls, err := eng.db.GetClsByFid(invoker.Db.WithContext(ctx), in.Fid, in.Biz, in.Act, in.Pagination)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetCountLogsByTargetId fail, err:" + err.Error())
	}
	log.Printf("cls--------------->"+"%+v\n", cls)
	tids := make([]string, 0, len(cls))
	for _, v := range cls {
		tids = append(tids, v.Tid)
	}
	// 根据tids查询tdetails
	data, err := eng.countCache.GetMap(ctx, &countv1.GetTdetailsByTidsReq{
		Tids: tids,
		Biz:  in.Biz,
		Act:  in.Act,
	})
	if err != nil {
		return nil, code.RedisOpError
	}
	tdetails := make([]*countv1.Tdetail, 0, len(tids))
	for _, v := range data {
		tdetails = append(tdetails, v)
	}
	return &countv1.DBGetTdetailsByFidRes{
		List:       tdetails,
		Pagination: in.Pagination,
	}, nil
}

// DBGetFidsTdetailByTid 根据target id。返回from id列表
// 例如查找某个用户的粉丝，target id为该用户uid，from id列表就是他的粉丝
func (eng *Engine) DBGetFidsTdetailByTid(ctx context.Context, in *countv1.DBGetFidsTdetailByTidReq) (*countv1.DBGetFidsTdetailByTidRes, error) {
	if !eng.checkBidParams(in.Biz, in.Act) {
		return nil, code.BidTypeParamsError
	}

	if in.Pagination.PageSize > 30 {
		return nil, code.TidIsLimit
	}

	// db 查询Fids列表
	countLogs, err := eng.db.GetCountLogsByTargetId(invoker.Db.WithContext(ctx), in.Tid, in.Biz, in.Act, in.Pagination)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetCountLogsByTargetId fail, err:" + err.Error())
	}

	// redis 查询Target总数
	nums, err := eng.countCache.HMGetString(ctx, cache.KeyHashTargetAct(in.Biz, in.Act, in.Tid), []string{cache.KeyFieldNum, cache.KeyFieldRealNum})
	if err != nil {
		elog.Error("CountCache.GetFids.HGET.KeyHashTargetAct", elog.Any("in", in), elog.FieldErr(err))
		return nil, errcodev1.ErrRedisError().WithMessage("HMGet fail, err:" + err.Error())
	}
	if len(nums) != 2 {
		return nil, errcodev1.ErrRedisError().WithMessage("HMGet fail, err:" + err.Error())
	}
	var res = &countv1.DBGetFidsTdetailByTidRes{
		Tid:        in.Tid,
		Num:        cast.ToInt64(nums[0]),
		RealNum:    cast.ToInt64(nums[1]),
		Status:     0, // TODO
		Fids:       make([]string, 0, len(countLogs)),
		Pagination: in.Pagination,
	}
	for _, v := range countLogs {
		res.Fids = append(res.Fids, v.Fid)
	}

	return res, nil
}

type GetTnumByFidItem struct {
	Fid string
	Num int64
	Err error
}

func (eng *Engine) GetTnumByFids(ctx context.Context, in *countv1.GetTnumByFidsReq) (res *countv1.GetTnumByFidsRes, err error) {
	if !eng.checkBidParams(in.Biz, in.Act) {
		return nil, code.BidTypeParamsError
	}

	if len(in.Fids) > 10 {
		return nil, code.FidIsLimit
	}

	// 批量
	ch := make(chan GetTnumByFidItem, 10)
	wg := sync.WaitGroup{}
	wg.Add(len(in.Fids))
	for _, v := range in.Fids {
		go func(fid string) {
			defer wg.Done()
			ret, e := eng.countCache.HLen(ctx, cache.KeyHashFromActStatus(in.Biz, in.Act, fid))
			if e != nil {
				ch <- GetTnumByFidItem{Err: e}
				return
			}
			ch <- GetTnumByFidItem{Fid: fid, Num: ret}
		}(v)
	}

	wg.Wait()
	close(ch)

	m := make(map[string]int64)
	for v := range ch {
		// 只要一个查询失败，则全部失败
		if v.Err != nil {
			return nil, fmt.Errorf("countCache fail, err:%w", v.Err)
		}
		m[v.Fid] = v.Num
	}
	return &countv1.GetTnumByFidsRes{Map: m}, nil
}
