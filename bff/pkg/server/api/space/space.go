package space

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	commonv1 "ecodepost/pb/common/v1"
	filev1 "ecodepost/pb/file/v1"
	spacev1 "ecodepost/pb/space/v1"
)

// CreateMember 创建空间成员
func CreateMember(c *bffcore.Context) {
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
	c.InjectSpc(guid)
	resp, err := invoker.GrpcSpace.AddSpaceMember(c.Ctx(), &spacev1.AddSpaceMemberReq{
		OperateUid: c.Uid(),
		SpaceGuid:  guid,
		AddUids:    req.AddUids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// MemberList @LHP @2022-10-26
func MemberList(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	var req GroupMemberListRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	resp, err := invoker.GrpcSpace.SpaceMemberList(c.Ctx(), &spacev1.SpaceMemberListReq{
		OperateUid: c.Uid(),
		SpaceGuid:  guid,
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

// DeleteMember 删除空间成员
func DeleteMember(c *bffcore.Context) {
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
	c.InjectSpc(guid)
	resp, err := invoker.GrpcSpace.DeleteSpaceMember(c.Ctx(), &spacev1.DeleteSpaceMemberReq{
		OperateUid: c.Uid(),
		SpaceGuid:  guid,
		DeleteUids: req.DelUids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

//
//func SearchMember(c *bffcore.Context) {
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
//	c.InjectSpc(guid)
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
//	resp, err := invoker.GrpcSpace.SpaceMemberList(c.Ctx(), &spacev1.SpaceMemberListReq{
//		OperateUid: c.Uid(),
//		CmtGuid:    c.CmtGuid(),
//		SpaceGuid:  guid,
//		Uids:       uids,
//	})
//	if err != nil {
//		c.EgoJsonI18N(err)
//		return
//	}
//	c.JSONListPage(resp.List, nil)
//}

func Emojis(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcFile.EmojiList(c.Ctx(), &filev1.EmojiListReq{
		Uid:       c.Uid(),
		SpaceGuid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type QuitMemberReq struct {
	Reason string `json:"reason"`
}

func QuitMember(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	var req QuitMemberReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectSpc(guid)
	_, err := invoker.GrpcSpace.QuitSpaceMember(c.Ctx(), &spacev1.QuitSpaceMemberReq{
		Uid:       c.Uid(),
		SpaceGuid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type ApplyMemberRequest struct {
	Reason string `json:"reason"`
}

// ApplyMember 申请成为成员
func ApplyMember(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid not empty", nil)
		return
	}
	var req ApplyMemberRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectSpc(guid)
	resp, err := invoker.GrpcSpace.AuditApplySpaceMember(c.Ctx(), &spacev1.AuditApplySpaceMemberReq{
		OperateUid: c.Uid(),
		TargetGuid: guid,
		AuditType:  commonv1.AUDIT_TYPE_SPACE,
		Reason:     req.Reason,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// Permission 空间的状态
// 1 INTERNAL 内部可见这个空间，并且能看到内容
//
//	isView  只要进入到社区，那么他就是看得见这个空间
//	isWrite 进入社区，虽然可以看到空间内容，但是他没办法写入，需要在页面上有这个提示，让他加入该空间，才能够写入 调用这个接口加入，/spaces/:guid/apply， internal空间调用后，就直接加入了
//	isMember 判断是否为成员
//
// 2 PRIVATE  内部可见这个空间，不能看到内容 （不存在这个类型了）
//
//	isView  进入到社区，但不在该空间，那么他看不见这个空间
//	isWrite 只有申请成功，加入到空间，才能写入 ，/spaces/:guid/apply， internal空间调用后，需要审核加入
//	isMember 判断是否为成员
//	auditStatus 如果不是空间成员，需要判断审核状态情况
//
// 3 Secret   不是成员，那么他看不见这个空间
//
//	isView  进入到社区，但不在该空间，那么他看不见这个空间
//	isWrite 只有申请成功，加入到空间，才能写入
//	isMember 判断是否为成员
func Permission(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "space guid not empty", nil)
		return
	}
	c.InjectSpc(guid)

	resp, err := invoker.GrpcSpace.GetSpacePermissionByUid(c.Ctx(), &spacev1.GetSpacePermissionByUidReq{
		OperateUid: c.Uid(),
		TargetGuid: guid,
		GuidType:   commonv1.CMN_GUID_SPACE,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}
