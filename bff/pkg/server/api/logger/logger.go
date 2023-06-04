package logger

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	commonv1 "ecodepost/pb/common/v1"
	loggerv1 "ecodepost/pb/logger/v1"
)

type ListPageRequest struct {
	CurrentPage int32              `form:"currentPage"` // 当前页数
	Event       commonv1.LOG_EVENT `form:"event"`
	Group       commonv1.LOG_GROUP `form:"group"`
	OperateUid  int64              `form:"operateUid"`
	TargetUid   int64              `form:"targetUid"`
}

func ListPage(c *bffcore.Context) {
	var req ListPageRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}

	grpcReq := &loggerv1.ListPageReq{
		OperateUid: c.Uid(),
		I18N:       c.GetString(bffcore.ContextLanguage),
		Pagination: &commonv1.Pagination{
			CurrentPage: req.CurrentPage,
		},
	}

	if req.Event != commonv1.LOG_EVENT_INVALID {
		grpcReq.SearchEvent = req.Event
	}

	if req.Group != commonv1.LOG_GROUP_INVALID {
		grpcReq.SearchGroup = req.Group
	}

	if req.OperateUid != 0 {
		grpcReq.SearchOperateUid = req.OperateUid
	}

	if req.TargetUid != 0 {
		grpcReq.SearchTargetUid = req.TargetUid
	}

	resp, err := invoker.GrpcLogger.ListPage(c.Ctx(), grpcReq)
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}

func EventAndGroupList(c *bffcore.Context) {
	resp, err := invoker.GrpcLogger.ListEventAndGroup(c.Ctx(), &loggerv1.ListEventAndGroupReq{
		OperateUid: c.Uid(),
		I18N:       c.GetString(bffcore.ContextLanguage),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}
