package space

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	commonv1 "ecodepost/pb/common/v1"
	spacev1 "ecodepost/pb/space/v1"
)

type CreateOrUpdateGroupRequest struct {
	Name                 string             `json:"name" binding:"required" label:"组名"`
	Icon                 string             `json:"icon"`
	Visibility           commonv1.CMN_VISBL `json:"visibility"`
	IsOpenReadMemberList bool               `json:"isOpenReadMemberList"` // 如果打开，属于这个分组下的用户，可以看到用户列表
}

func GroupInfo(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid不能为空", nil)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcSpace.SpaceGroupInfo(c.Ctx(), &spacev1.SpaceGroupInfoReq{
		Uid:            c.Uid(),
		SpaceGroupGuid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

func CreateGroup(c *bffcore.Context) {
	var req CreateOrUpdateGroupRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	resp, err := invoker.GrpcSpace.CreateSpaceGroup(c.Ctx(), &spacev1.CreateSpaceGroupReq{
		OperateUid:            c.Uid(),
		Name:                  req.Name,
		Icon:                  req.Icon,
		Visibility:            req.Visibility,
		IsAllowReadMemberList: req.IsOpenReadMemberList,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

func UpdateGroup(c *bffcore.Context) {
	var req CreateOrUpdateGroupRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid不能为空", nil)
		return
	}
	resp, err := invoker.GrpcSpace.UpdateSpaceGroup(c.Ctx(), &spacev1.UpdateSpaceGroupReq{
		OperateUid:            c.Uid(),
		SpaceGroupGuid:        guid,
		Name:                  req.Name,
		Icon:                  req.Icon,
		Visibility:            req.Visibility,
		IsAllowReadMemberList: req.IsOpenReadMemberList,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

func DeleteGroup(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid不能为空", nil)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcSpace.DeleteSpaceGroup(c.Ctx(), &spacev1.DeleteSpaceGroupReq{
		OperateUid:     c.Uid(),
		SpaceGroupGuid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// type ChangeGroupSortRequest struct {
//	AfterSpaceGroupGuid string `json:"afterSpaceGroupGuid"`
//	SpaceGroupGuid      string `json:"spaceGroupGuid" label:"组ID"`
// }
//
// func ChangeGroupSort(c *bffcore.Context) {
//	var req ChangeGroupSortRequest
//	if err := c.Bind(&req); err != nil {
//		c.JSONE(1, "参数错误:"+err.Error(), err)
//		return
//	}
//	resp, err := invoker.GrpcSpace.ChangeSpaceGroupSort(c.Ctx(), &spacev1.ChangeSpaceGroupSortReq{
//		OperateUid:          c.Uid(),
//		CmtGuid:             c.CmtGuid(),
//		SpaceGroupGuid:      req.SpaceGroupGuid,
//		AfterSpaceGroupGuid: req.AfterSpaceGroupGuid,
//	})
//	if err != nil {
//		c.EgoJsonI18N(err)
//		return
//	}
//	c.JSONOK(resp)
// }

type CreateMemberRequest struct {
	AddUids []int64 `json:"addUids"`
}

func CreateGroupMember(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	var req CreateMemberRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcSpace.AddSpaceGroupMember(c.Ctx(), &spacev1.AddSpaceGroupMemberReq{
		OperateUid:     c.Uid(),
		SpaceGroupGuid: guid,
		AddUids:        req.AddUids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type GroupMemberListRequest struct {
	CurrentPage int32 `form:"currentPage"` // 当前页数
}

func GroupMemberList(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	var req GroupMemberListRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcSpace.SpaceGroupMemberList(c.Ctx(), &spacev1.SpaceGroupMemberListReq{
		OperateUid:     c.Uid(),
		SpaceGroupGuid: guid,
		Pagination: &commonv1.Pagination{
			CurrentPage: req.CurrentPage,
		},
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}

type DeleteMemberRequest struct {
	DelUids []int64 `json:"delUids"`
}

func DeleteGroupMember(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	var req DeleteMemberRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcSpace.DeleteSpaceGroupMember(c.Ctx(), &spacev1.DeleteSpaceGroupMemberReq{
		OperateUid:     c.Uid(),
		SpaceGroupGuid: guid,
		DeleteUids:     req.DelUids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type SearchMemberRequest struct {
	Keyword string `form:"keyword"` // 当前页数
}

//
//func SearchGroupMember(c *bffcore.Context) {
//	guid := c.Param("guid")
//	if guid == "" {
//		c.JSONE(1, "guid not empty", nil)
//		return
//	}
//	var req SearchMemberRequest
//	if err := c.Bind(&req); err != nil {
//		c.JSONE(1, "参数错误:"+err.Error(), err)
//		return
//	}
//	c.InjectGuid(guid)
//	res, err := invoker.GrpcSearch.Query(c.Ctx(), &searchv1.QueryReq{
//		Keyword: req.Keyword,
//		BizType: commonv1.CMN_BIZ_USER,
//		Params:  map[string]string{"cmtNicknames.guid": c.CmtGuid()},
//	})
//	if err != nil {
//		c.EgoJsonI18N(err)
//		return
//	}
//	// 如果为空，直接返回空数据
//	if len(res.List) == 0 {
//		c.JSONListPage([]*commonv1.MemberRole{}, res.Pagination)
//		return
//	}
//	uids := make([]int64, 0, len(res.List))
//	for _, v := range res.List {
//		u, ok := v.Payload.(*searchv1.Item_User)
//		if ok {
//			uids = append(uids, u.User.Uid)
//		}
//	}
//	resp, err := invoker.GrpcSpace.SpaceGroupMemberList(c.Ctx(), &spacev1.SpaceGroupMemberListReq{
//		OperateUid:     c.Uid(),
//		CmtGuid:        c.CmtGuid(),
//		SpaceGroupGuid: guid,
//		Uids:           uids,
//	})
//	if err != nil {
//		c.EgoJsonI18N(err)
//		return
//	}
//	c.JSONListPage(resp.List, nil)
//}

type ApplyGroupMemberRequest struct {
	Reason string `json:"reason"`
}

// ApplyGroupMember 申请成为成员
func ApplyGroupMember(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	var req ApplyGroupMemberRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectGuid(guid)
	// resp, err := invoker.GrpcSpace.AuditApplySpaceMember(c.Ctx(), &spacev1.AuditApplySpaceMemberReq{
	//	OperateUid: c.Uid(),
	//	CmtGuid:    c.CmtGuid(),
	//	TargetGuid: guid,
	//	AuditType:  commonv1.AUDIT_TYPE_SPACE_GROUP,
	//	Reason:     req.Reason,
	// })
	// if err != nil {
	//	c.EgoJsonI18N(err)
	//	return
	// }
	c.JSONOK()
}

// GroupPermission 权限
func GroupPermission(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	resp, err := invoker.GrpcSpace.GetSpacePermissionByUid(c.Ctx(), &spacev1.GetSpacePermissionByUidReq{
		OperateUid: c.Uid(),
		TargetGuid: guid,
		GuidType:   commonv1.CMN_GUID_SPACE_GROUP,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type ListSpaceAndGroupRes struct {
	SpaceList      []*spacev1.TreeSpace      `json:"spaceList"`
	SpaceGroupList []*spacev1.TreeSpaceGroup `json:"spaceGroupList"`
}

func ListSpaceAndGroup(c *bffcore.Context) {
	res, err := invoker.GrpcSpace.ListSpaceAndGroup(c.Ctx(), &spacev1.ListSpaceAndGroupReq{
		OperateUid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(ListSpaceAndGroupRes{
		SpaceList:      res.SpaceList,
		SpaceGroupList: res.SpaceGroupList,
	})
}
