package mysql

import (
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	spacev1 "ecodepost/pb/space/v1"
	userv1 "ecodepost/pb/user/v1"

	"gorm.io/gorm"
)

type AuditIndex struct {
	Id        int64                 `gorm:"type:int(10) NOT NULL AUTO_INCREMENT;comment:ID;" json:"id"`                                           // ID
	Uid       int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:提交人;" json:"uid"`                                           // 提交人
	CmtGuid   string                `gorm:"type:varchar(191) NOT NULL DEFAULT '';comment:团队GUID;UNIQUE_INDEX" json:"cmtGuid"`                     // 团队GUID
	AuditGuid string                `gorm:"type:varchar(191) NOT NULL DEFAULT 0;comment:审核的Id;unique_index:idx_audit_guid_type" json:"auditGuid"` // 审核的Id
	AuditType commonv1.AUDIT_TYPE   `gorm:"type:int(10) NOT NULL DEFAULT 0;comment:类型;unique_index:idx_audit_guid_type" json:"auditType"`         // 类型
	Status    commonv1.AUDIT_STATUS `gorm:"type:int(10) NOT NULL DEFAULT 0;comment:审核情况;" json:"status"`                                          // 只记录最新的审核情况
	OpUid     int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:操作人;" json:"opUid"`                                         // 提交人
	Reason    string                `gorm:"type:longtext NOT NULL" json:"reason"`                                                                 // 只记录最新的原因
	OpReason  string                `gorm:"type:longtext NOT NULL" json:"opReason"`                                                               // 只记录最新的原因
	Ctime     int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:创建时间;index" json:"ctime"`                                   // 创建时间
	Utime     int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:更新时间;" json:"utime"`                                        // 更新时间
	Dtime     int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:删除时间;INDEX" json:"dtime"`                                   // 删除时间
}

// AuditLog index对应多个log
type AuditLog struct {
	Id           int32                 `gorm:"type:int(10) NOT NULL AUTO_INCREMENT;comment:ID;" json:"id"`                       // ID
	AuditIndexId int64                 `gorm:"type:int(10) NOT NULL DEFAULT 0;comment:audit_index_id;index" json:"auditIndexId"` // auditIndexId
	Uid          int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:提交人;" json:"uid"`                       // 提交人
	CmtGuid      string                `gorm:"type:varchar(191) NOT NULL DEFAULT '';comment:团队GUID;UNIQUE_INDEX" json:"cmtGuid"` // 团队GUID
	Status       commonv1.AUDIT_STATUS `gorm:"type:int(10) NOT NULL DEFAULT 0;comment:审核情况;" json:"status"`                      // 审核情况
	IsDone       int8                  `gorm:"type:tinyint(2) NOT NULL DEFAULT 0;comment:是否完成;" json:"isDone"`                   // 是否完成
	Reason       string                `gorm:"type:longtext NOT NULL" json:"reason"`                                             // 原因
	OpReason     string                `gorm:"type:longtext NOT NULL" json:"opReason"`                                           // 只记录最新的原因
	OpUid        int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:操作人;" json:"opUid"`                     // 提交人
	Ctime        int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:创建时间;index" json:"ctime"`               // 创建时间
	Utime        int64                 `gorm:"type:bigint(20) NOT NULL DEFAULT 0;comment:更新时间;" json:"utime"`                    // 更新时间
}

type AuditLogs []*AuditIndex

// TableName 设置表名
func (t AuditIndex) TableName() string {
	return "audit_index"
}

// TableName 设置表名
func (t AuditLog) TableName() string {
	return "audit_log"
}

type AuditIndexes []AuditIndex

func (list AuditIndexes) ToPb(userMap map[int64]*userv1.UserInfo) []*spacev1.AuditMember {
	output := make([]*spacev1.AuditMember, 0)
	for _, value := range list {
		output = append(output, value.ToPb(userMap))
	}
	return output
}

func (list AuditIndexes) ToUids() []int64 {
	output := make([]int64, 0)
	for _, value := range list {
		output = append(output, value.Uid)
	}
	return output
}

func (a AuditIndex) ToPb(userMap map[int64]*userv1.UserInfo) *spacev1.AuditMember {
	nickname := ""
	avatar := ""
	if userMap != nil {
		nickname = userMap[a.Uid].GetNickname()
		avatar = userMap[a.Uid].GetAvatar()
	}
	return &spacev1.AuditMember{
		AuditId:  a.Id,
		Uid:      a.Uid,
		Nickname: nickname,
		Avatar:   avatar,
		Reason:   a.Reason,
	}
}

// AuditLogCreate 创建一条记录
func AuditLogCreate(db *gorm.DB, data *AuditIndex, reason string) (err error) {
	auditInfo, err := AuditIndexInfoByUidAndGuid(db, data.OpUid, data.AuditGuid, commonv1.AUDIT_TYPE_SPACE)
	if err != nil {
		return errcodev1.ErrDbError().WithMessage("AuditApplySpaceMember fail, err: " + err.Error())
	}

	nowTime := time.Now().Unix()
	data.Ctime = nowTime
	data.Utime = nowTime
	data.Reason = reason
	if auditInfo.Id == 0 {
		data.Status = commonv1.AUDIT_STATUS_APPLY
		if err = db.Create(data).Error; err != nil {
			err = fmt.Errorf("audit_index create err: %w", err)
			return
		}
	} else {
		data.Status = commonv1.AUDIT_STATUS_RE_APPLY
		db.Model(data).Where("id = ?", auditInfo.Id).Updates(data)
	}

	if err = db.Create(&AuditLog{
		AuditIndexId: data.Id,
		Reason:       reason,
		Status:       data.Status,
		Ctime:        nowTime,
		Utime:        nowTime,
	}).Error; err != nil {
		err = fmt.Errorf("audit_log create err: %w", err)
		return
	}
	return
}

// AuditLogInfoByNotDone 根据PRI查询单条记录
func AuditLogInfoByNotDone(db *gorm.DB, auditIndexId int64) (resp AuditLog, err error) {
	if err = db.Where("audit_index_id = ? and is_done = 0", auditIndexId).Find(&resp).Error; err != nil {
		err = fmt.Errorf("AuditIndexInfo err: %w", err)
		return
	}
	return
}

// AuditIndexInfo 根据PRI查询单条记录
func AuditIndexInfo(db *gorm.DB, id int64) (resp AuditIndex, err error) {
	if err = db.Model(AuditIndex{}).Where("id = ?", id).Find(&resp).Error; err != nil {
		err = fmt.Errorf("audit_index info err: %w", err)
		return
	}
	return
}

// AuditIndexInfoByUidAndGuid 根据PRI查询单条记录
func AuditIndexInfoByUidAndGuid(db *gorm.DB, uid int64, targetGuid string, auditType commonv1.AUDIT_TYPE) (resp AuditIndex, err error) {
	if err = db.Model(AuditIndex{}).Where("uid = ? and audit_guid = ? and audit_type = ? ", uid, targetGuid, auditType.Number()).Find(&resp).Error; err != nil {
		err = fmt.Errorf("audit_index info err: %w", err)
		return
	}
	return
}

// AuditIndexUpdate 根据主键更新一条记录
func AuditIndexUpdate(db *gorm.DB, id int64, ups map[string]any) (err error) {
	ups["utime"] = time.Now().Unix()
	if err = db.Model(AuditIndex{}).Where("`id`=?", id).Updates(ups).Error; err != nil {
		err = fmt.Errorf("audit_index update err: %w", err)
		return
	}
	return
}

func AuditLogUpdate(db *gorm.DB, id int32, ups map[string]any) (err error) {
	ups["utime"] = time.Now().Unix()
	if err = db.Model(AuditLog{}).Where("`id`=?", id).Updates(ups).Error; err != nil {
		err = fmt.Errorf("audit_log update err: %w", err)
		return
	}
	return
}
