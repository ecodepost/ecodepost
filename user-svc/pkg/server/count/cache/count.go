package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	countv1 "ecodepost/pb/count/v1"
	"ecodepost/user-svc/pkg/code"
	"github.com/gotomicro/ego/core/elog"
	"github.com/samber/lo"

	"github.com/ego-component/eredis"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

// CountCacheI ...
type CountCacheI interface {
	// Incr Incr(ADD/SUB)操作
	Incr(context.Context, *countv1.SetReq) (*countv1.SetRes, error)
	// Update 支持
	Update(context.Context, *countv1.SetReq) (*countv1.SetRes, error)
	// Reset 重置用户所有的关系列表
	Reset(context.Context, *countv1.SetReq) (*countv1.SetRes, error)
	// GetMap 根据目标ID获取目标ID的状态
	GetMap(context.Context, *countv1.GetTdetailsByTidsReq) (map[string]*countv1.Tdetail, error)
	// GetTidByFid 获取目标ID的点赞/推荐列表
	GetTidByFid(context.Context, *countv1.GetTdetailsByFidReq) ([]*countv1.Tdetail, error)
	// GetFids 获取目标ID最近的来源ID列表，取LastTime最新50个
	GetFids(context.Context, *countv1.GetFidsTdetailByTidReq) (*countv1.GetFidsTdetailByTidRes, error)
}

var _ CountCacheI = CountCache{}

// CountCache 缓存句柄
type CountCache struct {
	*eredis.Component
}

// NewCountCache 初始化缓存cache
func NewCountCache(cache *eredis.Component) *CountCache {
	return &CountCache{
		Component: cache,
	}
}

// Incr ADD/SUB时设置来源ID与目标ID的关系
func (cr CountCache) Incr(ctx context.Context, in *countv1.SetReq) (data *countv1.SetRes, err error) {
	data = &countv1.SetRes{}
	// 如果是VIEW/CNT操作,则会更新用户的记录信息
	if lo.Contains([]commonv1.CNT_ACT{commonv1.CNT_ACT_VIEW, commonv1.CNT_ACT_CNT}, in.Act) {
		return cr.Update(ctx, in)
	}

	// 处理COLLECT/LIKE/UNLIKE/WATCH操作
	// 添加前先检测来源ID是否已执行过Add操作
	_, st, err := cr.getCurrentField(ctx, in)
	if err != nil {
		elog.Error("CountCache.Add.getCurrentField", elog.FieldErr(err), elog.Any("in", in), elog.Any("st", st))
		return
	}
	var delta int
	var actStatus code.ActStatusType
	switch in.Acti {
	case commonv1.CNT_ACTI_ADD:
		delta = 1 // ADD 操作强制将增加量设置为1
		actStatus = code.ActStatusAdded
	case commonv1.CNT_ACTI_SUB:
		delta = -1 // SUB 操作强制将增加量设置为-1
		actStatus = code.ActStatusSubbed
	default:
		return nil, fmt.Errorf("invalid act, %s", in.Act.String())
	}

	// 如果已经操作,则直接返回结果
	if (in.Acti == commonv1.CNT_ACTI_ADD && st == code.ActStatusAdded) || (in.Acti == commonv1.CNT_ACTI_SUB && st == code.ActStatusSubbed) {
		num, realNum, e := cr.GetTargetCount(ctx, in.Biz, in.Act, in.Tid)
		if e != nil {
			return
		}
		return &countv1.SetRes{Fid: in.Fid, Tid: in.Tid, Num: num, RealNum: realNum, Status: int32(st)}, nil
	}

	// 否则进行缓存更新操作
	// 目标增加num
	num, err := cr.HIncrBy(ctx, KeyHashTargetAct(in.Biz, in.Act, in.Tid), KeyFieldNum, delta)
	if err != nil {
		elog.Error("CountCache.IncrOnce.HINCRBY.KeyHashTargetAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	// 目标增加realNum
	realNum, err := cr.HIncrBy(ctx, KeyHashTargetAct(in.Biz, in.Act, in.Tid), KeyFieldRealNum, delta)
	if err != nil {
		elog.Error("CountCache.IncrOnce.HINCRBY.KeyHashTargetAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	// 往有序集合里面添加目标ID最近被哪些来源ID关注或者点赞
	_, err = cr.ZAdd(ctx, KeyListTargetLasttimeAct(in.Biz, in.Act, in.Tid), &redis.Z{Score: float64(time.Now().Unix()), Member: in.Fid})
	if err != nil {
		elog.Error("CountCache.IncrOnce.ZADD.KeyListTargetLasttimeAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	// 保持有序集合里面长度始终为最大值
	_, err = cr.ZRemRangeByRank(ctx, KeyListTargetLasttimeAct(in.Biz, in.Act, in.Tid), 0, -TargetMaxFromNum-1)
	if err != nil {
		elog.Error("CountCache.IncrOnce.ZADZREMRANGEBYRANKD.KeyListTargetLasttimeAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	// 添加来源ID对目标ID的操作记录
	_, err = cr.ZAdd(ctx, KeyListFromAct(in.Biz, in.Act, in.Fid), &redis.Z{Score: float64(time.Now().Unix()), Member: in.Tid})
	if err != nil {
		elog.Error("CountCache.IncrOnce.ZADD.KeyListFromAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	// 添加来源ID对目标ID操作状态
	err = cr.HSet(ctx, KeyHashFromActStatus(in.Biz, in.Act, in.Fid), in.Tid, actStatus)
	if err != nil {
		elog.Error("CountCache.IncrOnce.HSET.KeyHashFromActStatus", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	return &countv1.SetRes{
		Fid:     in.Fid,
		Tid:     in.Tid,
		Num:     num,
		RealNum: realNum,
		Status:  int32(actStatus),
	}, nil
}

// Update 直接更新目标点赞的计数值,val表示增量值
func (cr CountCache) Update(ctx context.Context, in *countv1.SetReq) (data *countv1.SetRes, err error) {
	var num int64
	// 目标增加数量
	num, err = cr.HIncrBy(ctx, KeyHashTargetAct(in.Biz, in.Act, in.Tid), KeyFieldNum, int(in.Val))
	if err != nil {
		elog.Error("CountCache.Update.HINCRBY.KeyHashTargetAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	// 如果计数是减少操作、防止减少变为0的操作
	if num < 0 {
		num = 0
		if err = cr.HSet(ctx, KeyHashTargetAct(in.Biz, in.Act, in.Tid), KeyFieldNum, num); err != nil {
			elog.Error("CountCache.Update.HSET.KeyHashTargetAct", elog.Any("in", in), elog.FieldErr(err))
			return
		}
	}

	return &countv1.SetRes{
		Fid:    in.Fid,
		Tid:    in.Tid,
		Num:    num,
		Status: int32(code.ActStatusAdded),
	}, nil
}

// Reset 取消所有关系
func (cr CountCache) Reset(ctx context.Context, in *countv1.SetReq) (data *countv1.SetRes, err error) {
	data = &countv1.SetRes{Fid: in.Fid}
	// 暂时只有用户浏览支持清空操作
	if !lo.Contains([]commonv1.CNT_ACT{commonv1.CNT_ACT_VIEW, commonv1.CNT_ACT_CNT}, in.Act) {
		return
	}

	// 删除用户所有的操作列表
	_, err = cr.Del(ctx, KeyListFromAct(in.Biz, in.Act, in.Fid))
	if err != nil {
		elog.Error("CountCache.Reset.DEL.KeyListFromAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	_, err = cr.Del(ctx, KeyHashFromActStatus(in.Biz, in.Act, in.Fid))
	if err != nil {
		elog.Error("CountCache.Reset.DEL.KeyHashFromActStatus", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	return data, nil
}

// GetMap 根据来源ID获取目标ID的状态
func (cr CountCache) GetMap(ctx context.Context, in *countv1.GetTdetailsByTidsReq) (data map[string]*countv1.Tdetail, err error) {
	data = make(map[string]*countv1.Tdetail)
	// 查询出所有target id的操作数量
	targetMap, err := cr.batchGetTargetCount(ctx, in.Biz, in.Act, in.Tids)
	if err != nil {
		return
	}
	var statsMap map[string]int32
	// 如果请求中包含Fid，则查询出来Fid是否对Tid进行操作
	if in.Fid != "" {
		statsMap = cr.batchGetTidStatus(ctx, in.Biz, in.Act, in.Fid, in.Tids)
	}
	fidMap := cr.batchGetTidRecentFids(ctx, in.Biz, in.Act, in.Tids, in.MaxFids)
	for tid, reply := range targetMap {
		data[tid] = &countv1.Tdetail{
			Tid:     tid,
			Num:     reply[0],
			RealNum: reply[1],
			Status:  statsMap[tid],
			Fids:    fidMap[tid],
		}
	}

	return
}

// GetTidByFid 获取来源ID的点赞/推荐列表
func (cr CountCache) GetTidByFid(ctx context.Context, in *countv1.GetTdetailsByFidReq) (data []*countv1.Tdetail, err error) {
	data = make([]*countv1.Tdetail, 0)
	key := KeyListFromAct(in.Biz, in.Act, in.Fid)
	start := in.GetOffset()
	end := start + in.GetLimit() - 1

	res, err := cr.ZRevRange(ctx, key, int64(start), int64(end))
	if err != nil {
		elog.Error("CountCache.GetTidByFid.ZREVRANGE.KeyListFromAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}
	if len(res) == 0 {
		return
	}

	targetMap, err := cr.batchGetTargetCount(ctx, in.Biz, in.Act, res)
	if err != nil {
		return
	}
	for _, tid := range res {
		var num int64
		val, ok := targetMap[tid]
		if ok {
			num = val[0]
		}

		tmp := &countv1.Tdetail{
			Tid:    tid,
			Num:    num,
			Status: int32(code.ActStatusAdded),
		}
		data = append(data, tmp)
	}
	return
}

// GetFids 获取目标ID最近的来源ID列表
func (cr CountCache) GetFids(ctx context.Context, in *countv1.GetFidsTdetailByTidReq) (data *countv1.GetFidsTdetailByTidRes, err error) {
	// 查询当前的目标ID被哪些来源ID点赞
	res, err := cr.ZRevRange(ctx, KeyListTargetLasttimeAct(in.Biz, in.Act, in.Tid), 0, int64(in.GetLimit()-1))
	if err != nil {
		elog.Error("CountCache.GetFids.ZREVRANGE.KeyListTargetLasttimeAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}
	fids := make([]string, 0, len(res))
	for _, v := range res {
		fids = append(fids, v)
	}

	// 查询当前目标ID的赞踩数量
	nums, err := cr.HMGetString(ctx, KeyHashTargetAct(in.Biz, in.Act, in.Tid), []string{KeyFieldNum, KeyFieldRealNum})
	if err != nil {
		elog.Error("CountCache.GetFids.HGET.KeyHashTargetAct", elog.Any("in", in), elog.FieldErr(err))
		return
	}
	num := cast.ToInt64(nums[0])
	realNum := cast.ToInt64(nums[1])

	// 用户关注状态
	var st code.ActStatusType
	// 查询用户是否推荐当前目标ID
	if in.Fid != "" {
		checkIn := &countv1.SetReq{Biz: in.Biz, Act: in.Act, Fid: in.Fid, Tid: in.Tid}
		_, st, err = cr.getCurrentField(ctx, checkIn)
		if err != nil {
			elog.Error("CountCache.GetFids.getCurrentField", elog.Any("in", in), elog.FieldErr(err))
			return
		}
	}

	return &countv1.GetFidsTdetailByTidRes{
		Tid:     in.Tid,
		Num:     num,
		RealNum: realNum,
		Status:  int32(st),
		Fids:    fids,
	}, nil
}

// getCurrentField 获取当前redis hash状态字段, num/real_num
func (cr CountCache) getCurrentField(ctx context.Context, in *countv1.SetReq) (field string, st code.ActStatusType, err error) {
	// 添加来源ID对目标ID的操作记录 点赞或者支持
	resp, err := cr.HGet(ctx, KeyHashFromActStatus(in.Biz, in.Act, in.Fid), in.Tid)
	res := cast.ToUint8(resp)
	// redis执行失败
	if err != nil && !errors.Is(err, eredis.Nil) {
		elog.Error("CountCache.IncrOnce.HSET.KeyHashFromActStatus", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	// 表示不存在操作记录
	if errors.Is(err, eredis.Nil) {
		return "", code.ActStatusType(res), nil
	}

	if lo.Contains([]code.ActStatusType{code.ActStatusAdded, code.ActStatusSubbed}, code.ActStatusType(res)) {
		return KeyFieldNum, code.ActStatusType(res), nil
	}

	return "", code.ActStatusType(res), nil
}

// batchGetTidStatus 批量获取来源ID对目标的操作状态
func (cr CountCache) batchGetTidStatus(ctx context.Context, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, fid string, tidList []string) (statsMap map[string]int32) {
	// 状态Map, <tid>:<actStatus>
	statsMap = make(map[string]int32)
	if len(tidList) == 0 {
		return
	}

	// 1.先查询出来源ID是否对目标ID进行操作
	in := &countv1.SetReq{Biz: biz, Act: act, Fid: fid}
	cmdKeys := make([]string, 0, len(tidList))
	for _, v := range tidList {
		cmdKeys = append(cmdKeys, v)
	}
	res, err := cr.HMGetString(ctx, KeyHashFromActStatus(in.Biz, in.Act, in.Fid), cmdKeys)
	if err != nil && !errors.Is(err, eredis.Nil) {
		elog.Error("CountCache.batchGetTidStatus.HSET.KeyHashFromActStatus", elog.Any("in", in), elog.FieldErr(err))
		return
	}

	if len(res) > 0 && len(res) == len(tidList) {
		for k, tid := range tidList {
			statsMap[tid] = cast.ToInt32(res[k])
		}
	}
	return
}

// GetTargetCount 获取目标ID的ACT数量
func (cr CountCache) GetTargetCount(ctx context.Context, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, tid string) (num, realNum int64, err error) {
	key := KeyHashTargetAct(biz, act, tid)
	reply, err := cr.HMGetString(ctx, key, []string{KeyFieldNum, KeyFieldRealNum})
	if err != nil && !errors.Is(err, eredis.Nil) {
		elog.Error("getTargetCount.KeyHashTargetAct", elog.FieldErr(err))
		return
	}
	if len(reply) != 2 {
		return
	}
	return cast.ToInt64(reply[0]), cast.ToInt64(reply[1]), nil
}

// batchGetTargetCount 批量获取目标ID的ACT数量
func (cr CountCache) batchGetTargetCount(ctx context.Context, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, tidList []string) (res map[string][]int64, err error) {
	res = make(map[string][]int64)
	// 查询目标ID的点赞点踩数值
	for _, tid := range tidList {
		var num, realNum int64
		num, realNum, err = cr.GetTargetCount(ctx, biz, act, tid)
		if err != nil {
			elog.Error("batchGetTargetCount.getTargetCount", elog.FieldErr(err))
			return
		}
		tmp := make([]int64, 0)
		tmp = append(tmp, num)
		tmp = append(tmp, realNum)
		res[tid] = tmp
	}
	return
}

// batchGetTidRecentFids 批量获取对tid近期act过的fid列表
func (cr CountCache) batchGetTidRecentFids(ctx context.Context, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, tids []string, limit int32) (fidMap map[string][]string) {
	fidMap = make(map[string][]string)
	for _, tid := range tids {
		// 查出目标ＩＤ最近的来源ID列表
		// 查询当前的目标ID被哪些来源ID点赞
		key := KeyListTargetLasttimeAct(biz, act, tid)
		res, err := cr.ZRevRange(ctx, key, 0, int64(limit-1))
		if err != nil {
			elog.Error("CountCache.batchGetTidRecentFids.ZREVRANGE.KeyListTargetLasttimeAct", elog.String("key", key), elog.FieldErr(err))
			continue
		}
		fidMap[tid] = res
	}

	return
}

// UpdateTargetRealNum 更新目标的真实点赞数
func (cr CountCache) UpdateTargetRealNum(ctx context.Context, req *countv1.SetReq, num, realUpNum int64) (err error) {
	key := KeyHashTargetAct(req.Biz, req.Act, req.Tid)
	err = cr.HMSet(ctx, key, map[string]any{KeyFieldNum: num, KeyFieldRealNum: realUpNum}, -1*time.Second)
	if err != nil && !errors.Is(err, eredis.Nil) {
		elog.Error("UpdateTargetRealNum.KeyHashTargetAct", elog.Any("tid", req.Tid), elog.FieldErr(err))
		return
	}

	return nil
}

// GetTargetRealUpNum 获取目标的真实点赞数
func (cr CountCache) GetTargetRealUpNum(ctx context.Context, req *countv1.SetReq) (upNum, realUpNum int64, err error) {
	key := KeyHashTargetAct(req.Biz, req.Act, req.Tid)
	reply, err := cr.HMGetString(context.Background(), key, []string{KeyFieldNum, KeyFieldRealNum})
	if err != nil && !errors.Is(err, eredis.Nil) {
		elog.Error("GetTargetRealUpNum.KeyHashTargetAct", elog.FieldErr(err))
		return
	}
	if len(reply) != 2 {
		return
	}
	return cast.ToInt64(reply[0]), cast.ToInt64(reply[1]), nil
}

// BatchGetTargetUpNum 批量获取计数前台数据/真实数据
// func (cr CountCache) BatchGetTargetUpNum(ctx context.Context, req *countv1.BatchGetTidRealNumReq) (list *countv1.BatchGetTidRealNumRes) {
// 	list = &countv1.BatchGetTidRealNumRes{}
// 	bidTypeMap := make(map[string]*countv1.RealNums, 0)
// 	for key, item := range req.GetNums() {
// 		strs := strings.Split(key, ".")
// 		upNumMap := make(map[string]*countv1.RealNum, 0)
// 		for _, tid := range item.GetList() {
// 			biz := commonv1.CMN_BIZ(cast.ToInt32(strs[0]))
// 			act := commonv1.CNT_ACT(cast.ToInt32(strs[1]))
// 			upNum, realUpNum, err := cr.GetTargetRealUpNum(ctx, &countv1.SetReq{Bid: req.GetBid(), Biz: biz, Act: act, Tid: tid})
// 			if err != nil {
// 				elog.Error("CountCache.BatchGetTargetUpNum", elog.FieldErr(err))
// 				continue
// 			}
// 			tmp := &countv1.RealNum{
// 				RealNum: realUpNum,
// 				Num:     upNum,
// 			}
// 			upNumMap[tid] = tmp
// 		}
//
// 		bidTmp := &countv1.RealNums{Map: upNumMap}
// 		bidTypeMap[key] = bidTmp
// 	}
//
// 	list.Map = bidTypeMap
// 	return list
// }
