package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	countv1 "ecodepost/pb/count/v1"
	"ecodepost/user-svc/pkg/util"
	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"gorm.io/gorm"
)

// CountNum 响应计数值
type CountNum struct {
	// RealNum 真实计数值
	RealNum int64
	// Num 前台计数值 当不存在策略时，真实计数值=前台计数值
	Num int64
}

// CounterInfo 计数信息
type CounterInfo struct {
	// Tid 计数的目标ID
	Tid string
	// Fid 计数的来源ID
	Fid string
	// Biz 计数的事件类型或者业务ID 对应 bid_type
	Biz commonv1.CMN_BIZ
	// Act 计数的事件类型或者业务ID 对应 bid_type
	Act commonv1.CNT_ACT
	// Acti 动作指令
	Acti commonv1.CNT_ACTI
	// RealNum 真实计数值
	RealNum int64
	// Num 前台计数值 当不存在策略时，真实计数值=前台计数值
	Num int64
	// BaseNum 基础打底计数值
	BaseNum int64
	// Utime 更新时间
	Utime int64
}

type Ups = map[string]any

// CountTable 计数表
type CountTable struct {
	CountStat        CountStat
	CountLog         CountLog
	CountBehaviorLog CountBehaviorLog
}

// NewCountTable 初始化一个计数服务
func NewCountTable() *CountTable {
	return &CountTable{
		// CountLog:   newCountLogTable(),
		// CountStat:        newCountStatTable(),
		// CountBehaviorLog: newBehaviorLogInfoTable(),
	}
}

// IncrUpdateStat 增量更新统计表计数
func (ct *CountTable) IncrUpdateStat(db *gorm.DB, num, realNum int64, info *CounterInfo) (err error) {
	num, err = ct.CountStat.InsertOrUpdate(db, CountStat{
		Tid:     info.Tid,
		Biz:     info.Biz,
		Act:     info.Act,
		Num:     num,
		RealNum: realNum,
		Utime:   time.Now().Unix(),
	})
	if err != nil {
		elog.Error("CountTable.IncrUpdateStat", elog.String("tid", info.Tid), elog.Any("info", info), elog.FieldErr(err))
		return
	}
	return nil
}

// SetUpdateStatRealNum 设置更新统计表真实计数值
func (ct *CountTable) SetUpdateStatRealNum(ctx context.Context, info *CounterInfo) (err error) {
	err = ct.CountStat.Update(ctx, egorm.Conds{
		"biz": int32(info.Biz),
		"tid": info.Tid,
		"act": int32(info.Act),
	}, Ups{
		"ctime":    time.Now().Unix(),
		"real_num": info.RealNum,
		"num":      info.Num,
	})
	return err
}

// SetUpdateStatFrontNum 设置更新统计表前台计数值
func (ct *CountTable) SetUpdateStatFrontNum(ctx context.Context, info *CounterInfo) (err error) {
	ups := Ups{
		"num":   info.Num,
		"ctime": time.Now().Unix(),
	}
	if info.BaseNum > 0 {
		ups["base_num"] = info.BaseNum
	}
	err = ct.CountStat.Update(ctx, egorm.Conds{"tid": info.Tid, "biz": int32(info.Biz), "act": int32(info.Act)}, ups)
	return nil
}

// StoreCountStat 计数表新增一条记录
func (ct *CountTable) StoreCountStat(ctx context.Context, info *CounterInfo) (err error) {
	err = ct.CountStat.Create(ctx, &CountStat{
		Tid:     info.Tid,
		Biz:     info.Biz,
		Act:     info.Act,
		Num:     info.Num,
		RealNum: info.RealNum,
		BaseNum: info.BaseNum,
		Utime:   info.Utime,
	})
	return err
}

// GetRealNum 获取日志表计数值
func (ct *CountTable) GetRealNum(db *gorm.DB, info *CounterInfo) (realNum int64) {
	// TODO 性能优化只取realNum
	res, err := ct.CountLog.Find(db, egorm.Conds{
		"tid": info.Tid,
		"biz": int32(info.Biz),
		"act": int32(info.Act),
		"fid": info.Fid,
	})
	if err != nil {
		elog.Error("CountLog load fail", elog.FieldErr(err))
		return 0
	}
	return res.RealNum
}

// GetStatInfo 获取统计表基础随机值
func (ct *CountTable) GetStatInfo(ctx context.Context, info *CounterInfo) (backInfo *CounterInfo, err error) {
	res, err := ct.CountStat.Find(ctx, egorm.Conds{"biz": int32(info.Biz), "act": int32(info.Act), "tid": info.Tid}, CountStat{
		Num:     info.Num,
		RealNum: info.RealNum,
		BaseNum: info.BaseNum,
		Utime:   info.Utime,
	})
	if err != nil {
		return nil, fmt.Errorf("CountStat load fail, %w", err)
	}
	return &CounterInfo{
		Tid:     res.Tid,
		Biz:     res.Biz,
		Act:     res.Act,
		RealNum: res.RealNum,
		Num:     res.Num,
		BaseNum: res.BaseNum,
		Utime:   res.Utime,
	}, nil
}

// IncrUpdateLog 增量更新日志表计数值
func (ct *CountTable) IncrUpdateLog(db *gorm.DB, num, realNum int64, info *CounterInfo, client, did, ip string) (err error) {
	err = ct.CountLog.UpdateRow(db, CountLog{
		Tid:     info.Tid,
		Biz:     info.Biz,
		Act:     info.Act,
		Fid:     info.Fid,
		Num:     num,
		RealNum: realNum,
		Utime:   time.Now().Unix(),
		Ct:      client,
		Did:     did,
		Ip:      ip,
	})
	return err
}

// SetUpdateLogRealNum 设置更新日志表计数值
func (ct *CountTable) SetUpdateLogRealNum(db *gorm.DB, info *CounterInfo) (err error) {
	err = ct.CountLog.Update(db, egorm.Conds{
		"tid": info.Tid,
		"biz": info.Biz,
		"act": info.Act,
		"fid": info.Fid,
	}, Ups{
		"real_num": info.RealNum,
		"utime":    info.Utime,
	})
	return err
}

// ScanCountStatNum 获取目标对应事件最终计数值
func (ct *CountTable) ScanCountStatNum(db *gorm.DB, tid string, bizs []commonv1.CMN_BIZ, act commonv1.CNT_ACT) (cntMap map[string]*CountNum, err error) {
	bts := make([]int32, 0, len(bizs))
	for _, v := range bizs {
		bts = append(bts, int32(v))
	}
	res, err := ct.CountStat.List(db, egorm.Conds{"tid": tid, "biz": bts, "act": int32(act)}, nil)
	if err != nil {
		return nil, fmt.Errorf("CountStat scan fail, %w", err)
	}
	cntMap = make(map[string]*CountNum)
	for _, val := range res {
		cntMap[util.BizActKey(val.Biz, val.Act)] = &CountNum{
			RealNum: val.RealNum,
			Num:     val.Num,
		}
	}
	return
}

func ActKeyConds(aks []*countv1.BAK) (string, []any) {
	if len(aks) == 0 {
		return "", nil
	}
	sql := ""
	binds := []any{}
	for _, v := range aks {
		sql += fmt.Sprintf("OR (biz = ? AND act = ?) ")
		binds = append(binds, v.Biz.Number(), v.Act.Number())
	}
	return "(" + strings.TrimLeft(sql, "OR ") + ")", binds
}

func (ct *CountTable) ScanCountStatNumByActKeys(db *gorm.DB, tid string, actKeys []*countv1.BAK) (cntMap map[string]*CountNum, err error) {
	sql, binds := ActKeyConds(actKeys)
	res, err := ct.CountStat.ListX(db, egorm.Conds{"tid": tid}, nil, sql, binds)
	if err != nil {
		return nil, fmt.Errorf("CountStat scan fail, %w", err)
	}
	cntMap = make(map[string]*CountNum)
	for _, val := range res {
		cntMap[util.BizActKey(val.Biz, val.Act)] = &CountNum{
			RealNum: val.RealNum,
			Num:     val.Num,
		}
	}
	return
}

// ScanCountLogNum TODO
func (ct *CountTable) ScanCountLogNum(db *gorm.DB, tid string, bizs commonv1.CMN_BIZ, act commonv1.CNT_ACT) (cntMap map[string]int64, err error) {
	res, err := ct.CountLog.List(db, egorm.Conds{"tid": tid, "biz": int32(bizs), "act": int32(act)}, nil)
	if err != nil {
		return nil, fmt.Errorf("CountLog scan fail, %w", err)
	}
	cntMap = make(map[string]int64)
	for _, val := range res {
		cntMap[val.Fid] = val.RealNum
	}
	return
}

// GetCountLogsByTargetId 获取目标ID来源列表
func (ct *CountTable) GetCountLogsByTargetId(db *gorm.DB, tid string, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, page *commonv1.Pagination) (list []*CountLog, err error) {
	list = make([]*CountLog, 0)
	conds := egorm.Conds{
		"tid": tid,
		"biz": int32(biz),
		"act": int32(act),
		"num": egorm.Cond{Op: ">", Val: 0},
	}
	res, err := ct.CountLog.List(db, conds, page)
	if err != nil {
		return nil, fmt.Errorf("CountLog scan fail, %w", err)
	}
	list = make([]*CountLog, 0, len(res))
	for _, val := range res {
		list = append(list, &CountLog{
			Tid:     val.Tid,
			Biz:     val.Biz,
			Act:     val.Act,
			Fid:     val.Fid,
			Num:     val.Num,
			RealNum: val.RealNum,
			Utime:   val.Utime,
		})
	}
	return
}

// GetClsByFid 获取目标ID来源列表
func (ct *CountTable) GetClsByFid(db *gorm.DB, fid string, biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, page *commonv1.Pagination) (list []*CountLog, err error) {
	list = make([]*CountLog, 0)
	res, err := ct.CountLog.List(db, egorm.Conds{"fid": fid, "biz": int32(biz), "act": int32(act)}, page)
	if err != nil {
		return nil, fmt.Errorf("CountLog scan fail, %w", err)
	}
	list = make([]*CountLog, 0, len(res))
	for _, val := range res {
		list = append(list, &CountLog{
			Tid:     val.Tid,
			Biz:     val.Biz,
			Act:     val.Act,
			Fid:     val.Fid,
			Num:     val.Num,
			RealNum: val.RealNum,
			Utime:   val.Utime,
		})
	}
	return
}
