package cache

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/user-svc/pkg/util"
)

const (
	// KeyHashCntTid 目标ID(t:tid)被Act(a:ACT)次数 <bizKey>:<tid>
	KeyHashCntTid = "cnt_t_a:%s:%s"
	// KeyListCntTidLasttime 目标ID(t:tid)被lasttime(t:lasttime)Act(a:ACT)的来源ID列表(有序集合) <bizKey>:<tid>
	KeyListCntTidLasttime = "cnt_t_ta:%s:%s"
	// KeyHashCntFidActStatus 来源(f:fid)对目标(t:tid)的操作状态(s:actStatus)  <bizKey>:<fid>
	KeyHashCntFidActStatus = "cnt_f_ts:%s:%s"
	// KeyListCntFid 来源ID(f:fid)ACT(a:ACT)列表 sort集合(l:list)排序 <bizKey>:<fid>
	KeyListCntFid = "cnt_f_al:%s:%s"

	// KeyFieldNum 点赞/支持属性
	KeyFieldNum = "num"
	// KeyFieldRealNum 真实点赞/支持属性
	KeyFieldRealNum = "real_num"

	// TargetMaxFromNum 目标ID最多保存10个来源ID
	TargetMaxFromNum = 50
)

// KeyHashTargetAct 目标ID被ACT次数 KEY
func KeyHashTargetAct(biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, tid string) string {
	return fmt.Sprintf(KeyHashCntTid, util.BizActKey(biz, act), tid)
}

// KeyListTargetLasttimeAct 目标ID被last time ACT的来源ID列表 KEY
func KeyListTargetLasttimeAct(biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, tid string) string {
	return fmt.Sprintf(KeyListCntTidLasttime, util.BizActKey(biz, act), tid)
}

// KeyHashFromActStatus 来源ID对目标ID操作状态
func KeyHashFromActStatus(biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, fid string) string {
	return fmt.Sprintf(KeyHashCntFidActStatus, util.BizActKey(biz, act), fid)
}

// KeyListFromAct 来源ID的ACT目标列表 KEY
func KeyListFromAct(biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, fid string) string {
	return fmt.Sprintf(KeyListCntFid, util.BizActKey(biz, act), fid)
}
