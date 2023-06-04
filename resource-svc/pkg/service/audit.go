package service

import (
	"context"
	"fmt"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	spacev1 "ecodepost/pb/space/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/ego-component/egorm"
	"gorm.io/gorm"
)

type audit struct {
}

func (*audit) ListPage(db *gorm.DB, guid string, auditType commonv1.AUDIT_TYPE, reqList *commonv1.Pagination) (respList mysql.AuditIndexes, err error) {
	if reqList.PageSize == 0 {
		reqList.PageSize = 20
	}
	if reqList.CurrentPage == 0 {
		reqList.CurrentPage = 1
	}
	conds := egorm.Conds{}
	conds["dtime"] = 0
	conds["audit_type"] = auditType.Number()
	if reqList.Sort == "" {
		reqList.Sort = "id desc"
	}
	conds["audit_guid"] = guid
	sql, binds := egorm.BuildQuery(conds)

	listDb := db.Select("id,ctime,uid,cmt_guid,audit_guid,audit_type,status,reason").
		Joins("left join `audit_log` as c ON a.id = c.audit_index_id").Where(sql, binds...)
	err = listDb.Count(&reqList.Total).Error
	if err != nil {
		return nil, fmt.Errorf("audit list page fail,err: %w", err)
	}
	err = listDb.Order(reqList.Sort).Offset(int(reqList.CurrentPage-1) * int(reqList.PageSize)).Limit(int(reqList.PageSize)).Find(&respList).Error
	if err != nil {
		return nil, fmt.Errorf("audit list page fail2,err: %w", err)
	}
	return
}

func (*audit) Log(db *gorm.DB, auditId int32, auditType commonv1.AUDIT_TYPE) (list []mysql.AuditLog, err error) {
	err = db.Select("b.id,b.ctime,b.utime,b.status,b.reason,b.audit_index_id,b.is_done").Table("audit_index as a").
		Joins("left join `audit_log` as b ON a.id = b.audit_index_id").
		Where("a.id = ? and a.audit_type = ?", auditId, int32(auditType)).Order("b.id asc").Find(&list).Error
	return
}

func (*audit) Pass(ctx context.Context, auditId int64, opUid int64, reason string) (auditIndexInfo mysql.AuditIndex, err error) {
	// todo 验证这个用户是不是这个社区的
	auditIndexInfo, err = mysql.AuditIndexInfo(invoker.Db.WithContext(ctx), auditId)
	if err != nil {
		err = fmt.Errorf("audit pass get AuditIndexInfo failed, err: %w", err)
		return
	}

	if auditIndexInfo.Id == 0 {
		err = fmt.Errorf("audit pass get AuditIndexInfo failed, err: not exist")
		return
	}

	if auditIndexInfo.Status == commonv1.AUDIT_STATUS_PASS {
		err = fmt.Errorf("audit already pass")
		return
	}

	auditLogInfo, err := mysql.AuditLogInfoByNotDone(invoker.Db.WithContext(ctx), auditIndexInfo.Id)
	if err != nil {
		err = fmt.Errorf("audit pass get AuditLogInfoByNotDone failed, err: %w", err)
		return
	}
	if auditLogInfo.Id == 0 {
		err = fmt.Errorf("audit pass get AuditLogInfoByNotDone failed, err: not exist")
		return
	}

	tx := invoker.Db.WithContext(ctx).Begin()
	err = mysql.AuditIndexUpdate(tx, auditIndexInfo.Id, map[string]any{
		"status":    commonv1.AUDIT_STATUS_PASS.Number(),
		"op_uid":    opUid,
		"op_reason": reason,
	})
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("audit pass get AuditIndexUpdate failed, err: %w", err)
		return
	}
	err = mysql.AuditLogUpdate(tx, auditLogInfo.Id, map[string]interface{}{
		"status":    commonv1.AUDIT_STATUS_PASS.Number(),
		"op_uid":    opUid,
		"op_reason": reason,
		"is_done":   1,
	})
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("audit pass get AuditLogUpdate failed, err: %w", err)
		return
	}
	nowTime := time.Now().Unix()
	userInfo, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{
		Uid: auditIndexInfo.Uid,
	})
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("audit pass get AuditLogUpdate failed, err: %w", err)
		return
	}

	switch auditIndexInfo.AuditType {
	case commonv1.AUDIT_TYPE_SPACE:
		err = Space.CreateMember(ctx, tx, &mysql.SpaceMember{
			Ctime:     nowTime,
			Utime:     nowTime,
			Uid:       auditIndexInfo.Uid,
			Nickname:  userInfo.User.GetNickname(),
			Guid:      auditIndexInfo.AuditGuid,
			CreatedBy: opUid,
		})
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("audit change good status failed, err: %w", err)
			return
		}
	case commonv1.AUDIT_TYPE_SPACE_GROUP:
		err = mysql.SpaceGroupMemberCreate(tx, &mysql.SpaceGroupMember{
			Ctime:     nowTime,
			Utime:     nowTime,
			Uid:       auditIndexInfo.Uid,
			Nickname:  userInfo.User.GetNickname(),
			Guid:      auditIndexInfo.AuditGuid,
			CreatedBy: opUid,
		})
		if err != nil {
			tx.Rollback()
			err = fmt.Errorf("audit change good status failed, err: %w", err)
			return
		}
	}

	tx.Commit()
	return
}

func (*audit) Reject(ctx context.Context, auditId int64, opUid int64, reason string) (err error) {
	auditIndexInfo, err := mysql.AuditIndexInfo(invoker.Db.WithContext(ctx), auditId)
	if err != nil {
		err = fmt.Errorf("audit pass get AuditIndexInfo failed, err: %w", err)
		return
	}

	if auditIndexInfo.Id == 0 {
		err = fmt.Errorf("audit pass get AuditIndexInfo failed, err: not exist")
		return
	}

	if auditIndexInfo.Status == commonv1.AUDIT_STATUS_REJECT {
		err = fmt.Errorf("audit already reject")
		return
	}
	auditLogInfo, err := mysql.AuditLogInfoByNotDone(invoker.Db.WithContext(ctx), auditIndexInfo.Id)
	if err != nil {
		err = fmt.Errorf("audit pass get AuditLogInfoByNotDone failed, err: %w", err)
		return
	}
	if auditLogInfo.Id == 0 {
		err = fmt.Errorf("audit pass get AuditLogInfoByNotDone failed, err: not exist")
		return
	}

	tx := invoker.Db.WithContext(ctx).Begin()
	err = mysql.AuditIndexUpdate(tx, auditIndexInfo.Id, map[string]interface{}{
		"status":    int(commonv1.AUDIT_STATUS_REJECT),
		"op_reason": reason,
		"op_uid":    opUid,
	})
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("audit pass get AuditIndexUpdate failed, err: %w", err)
		return
	}
	err = mysql.AuditLogUpdate(tx, auditLogInfo.Id, map[string]interface{}{
		"status":    int(commonv1.AUDIT_STATUS_REJECT),
		"op_reason": reason,
		"op_uid":    opUid,
		"is_done":   1,
	})
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("audit pass get AuditLogUpdate failed, err: %w", err)
		return
	}
	tx.Commit()
	return nil
}

func (*audit) Map(db *gorm.DB, ids []int64) (map[int64]*spacev1.AuditIndex, error) {
	var auditIndexes mysql.AuditIndexes
	err := db.Model(mysql.AuditIndex{}).Where("id in (?)", ids).Find(&auditIndexes).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Map Fail, err: " + err.Error())
	}
	// 响应map数据
	uMap := make(map[int64]*spacev1.AuditIndex)
	for _, v := range auditIndexes {
		uMap[v.Id] = &spacev1.AuditIndex{
			AuditId:  v.Id,
			Status:   v.Status,
			Reason:   v.Reason,
			OpReason: v.OpReason,
		}
	}
	return uMap, err
}
