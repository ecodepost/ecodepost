package my

import (
	"context"
	"encoding/json"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	commonv1 "ecodepost/pb/common/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	spacev1 "ecodepost/pb/space/v1"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type NotificationListReq struct {
	bffcore.Pagination
	Types []commonv1.NTF_TYPE `form:"types" json:"types"` // 为空表示不限定查询类型
}

type NotificationItem struct {
	Type             commonv1.NTF_TYPE   `json:"type,omitempty"`             // 消息类型
	TargetId         string              `json:"targetId,omitempty"`         // 对象ID
	Link             string              `json:"link,omitempty"`             // 消息外链
	Meta             json.RawMessage     `json:"meta,omitempty"`             // 附属数据
	Status           commonv1.NTF_STATUS `json:"status,omitempty"`           // 消息状态
	NotificationId   int64               `json:"notificationId,omitempty"`   // 通知ID
	NotificationTime int64               `json:"notificationTime,omitempty"` // 通知时间
	SenderId         string              `json:"senderId,omitempty"`         // 发送者Id, 设置成string方便后续扩展, 目前可以全部当成uid来处理
	Id               int64               `json:"-"`                          // 唯一id
}

// NotificationList 查询通知列表
// @Tags Notification
// @Failure 200 {object} bffcore.Res{data=bffcore.ResPageData{list=[]notifyv1.ListUserNotificationItem}}
func NotificationList(c *bffcore.Context) {
	if c.Uid() == 0 {
		c.JSONListPage([]struct{}{}, nil)
		return
	}

	var req NotificationListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	res, err := invoker.GrpcNotify.ListUserNotification(c.Ctx(), &notifyv1.ListUserNotificationReq{
		Uid:        c.Uid(),
		Pagination: req.Pagination.ToPb(),
		Types:      req.Types,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	auditNis := make(map[int64]NotificationItem)
	list := make([]NotificationItem, 0, len(res.List))
	for _, v := range res.List {
		ni := NotificationItem{
			Type:             v.Type,
			TargetId:         v.TargetId,
			Link:             v.Link,
			Status:           v.Status,
			NotificationId:   v.NotificationId,
			NotificationTime: v.NotificationTime,
			SenderId:         v.SenderId,
			Id:               v.Id,
		}
		if len(v.Meta) != 0 {
			ni.Meta = v.Meta
		}
		if lo.Contains([]commonv1.NTF_TYPE{commonv1.NTF_TYPE_APPLY_SPACE, commonv1.NTF_TYPE_APPLY_SPACE_GROUP}, v.Type) {
			auditNis[v.Id] = ni
		}
		list = append(list, ni)
	}

	// 对meta等数据根据不同业务做额外处理
	list, e := HandleNis(c.Ctx(), list, auditNis)
	if e != nil {
		c.EgoJsonI18N(e)
		return
	}

	c.JSONListPage(list, res.Pagination)
}

// HandleNis 对nis做额外的业务处理, 追加meta字段, 后续可以增加其他map[int64]NotificationItem
func HandleNis(ctx context.Context, list []NotificationItem, auditNis map[int64]NotificationItem) ([]NotificationItem, error) {
	// 处理审核类型notificationItem
	if len(auditNis) != 0 {
		auditIds := make([]int64, 0, len(auditNis))
		for _, v := range auditNis {
			auditIds = append(auditIds, cast.ToInt64(v.TargetId))
		}
		auditMap, err := invoker.GrpcSpace.AuditMapByIds(ctx, &spacev1.AuditMapByIdsReq{AuditIds: auditIds})
		if err != nil {
			return nil, err
		}

		for k, oldNi := range auditNis {
			newNi := oldNi
			if len(oldNi.Meta) != 0 {
				var m map[string]any
				if e := json.Unmarshal(oldNi.Meta, &m); e != nil {
					return nil, e
				}
				info := auditMap.AuditMap[cast.ToInt64(oldNi.TargetId)]
				m["status"] = info.Status
				m["reason"] = info.Reason
				m["opReason"] = info.OpReason
				newNi.Meta, _ = json.Marshal(m)
				auditNis[k] = newNi
			}
		}
	}

	// 修改返回的list数据
	for k := range list {
		switch list[k].Type {
		case commonv1.NTF_TYPE_APPLY_SPACE, commonv1.NTF_TYPE_APPLY_SPACE_GROUP:
			list[k] = auditNis[list[k].Id]
		}
	}

	return list, nil
}
