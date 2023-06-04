package my

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	spacev1 "ecodepost/pb/space/v1"

	"github.com/spf13/cast"
)

type NotificationAuditPassRequest struct {
	OpReason string `json:"opReason"` // 管理员理由
}

// NotificationAuditPass 查询审核的通知列表
// @Tags Notification
func NotificationAuditPass(c *bffcore.Context) {
	auditId := cast.ToInt64(c.Param("auditId"))
	if auditId == 0 {
		c.JSONE(1, "audit id cant empty", nil)
		return
	}

	var req NotificationAuditPassRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	_, err := invoker.GrpcSpace.AuditPassSpaceMember(c.Ctx(), &spacev1.AuditPassSpaceMemberReq{
		Id:         auditId,
		OperateUid: c.Uid(),
		OpReason:   req.OpReason,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	c.JSONOK()
}

type NotificationAuditRejectRequest struct {
	OpReason string `json:"opReason"` // 管理员理由
}

// NotificationAuditReject 查询审核的通知列表
// @Tags Notification
func NotificationAuditReject(c *bffcore.Context) {
	auditId := cast.ToInt64(c.Param("auditId"))
	if auditId == 0 {
		c.JSONE(1, "audit id cant empty", nil)
		return
	}

	var req NotificationAuditRejectRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	_, err := invoker.GrpcSpace.AuditRejectSpaceMember(c.Ctx(), &spacev1.AuditRejectSpaceMemberReq{
		Id:         auditId,
		OperateUid: c.Uid(),
		OpReason:   req.OpReason,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
