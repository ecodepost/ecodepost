package pms

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	errcodev1 "ecodepost/pb/errcode/v1"

	commonv1 "ecodepost/pb/common/v1"
	pmsv1 "ecodepost/pb/pms/v1"

	"github.com/spf13/cast"
)

// RoleList (0613) 获取全部role列表
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.GetRoleListRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func RoleList(c *bffcore.Context) {
	resp, err := invoker.GrpcPms.GetRoleList(c.Ctx(), &pmsv1.GetRoleListReq{
		OperateUid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// UserRoleIds (0613) 获取用户uid全部role ids
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.GetRoleListRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func UserRoleIds(c *bffcore.Context) {
	uid := cast.ToInt64(c.Param("uid"))
	if uid == 0 {
		c.EgoJsonI18N(errcodev1.ErrUidEmpty())
		return
	}

	resp, err := invoker.GrpcPms.GetRoleIds(c.Ctx(), &pmsv1.GetRoleIdsReq{
		Uid: uid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required" label:"角色名称"`
}

// CreateRole (0616,0617) 创建角色
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.CreateRoleRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func CreateRole(c *bffcore.Context) {
	var req CreateRoleRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}

	resp, err := invoker.GrpcPms.CreateRole(c.Ctx(), &pmsv1.CreateRoleReq{
		OperateUid: c.Uid(),
		Name:       req.Name,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// UpdateRole (0707) 更新角色
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.UpdateRoleRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func UpdateRole(c *bffcore.Context) {
	id := cast.ToInt64(c.Param("roleId"))
	if id == 0 {
		c.JSONE(1, "role id can't empty", nil)
		return
	}

	var req CreateRoleRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}

	resp, err := invoker.GrpcPms.UpdateRole(c.Ctx(), &pmsv1.UpdateRoleReq{
		RoleId:     id,
		OperateUid: c.Uid(),
		Name:       req.Name,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

func DeleteRole(c *bffcore.Context) {
	id := cast.ToInt64(c.Param("roleId"))
	if id == 0 {
		c.JSONE(1, "role id can't empty", nil)
		return
	}
	resp, err := invoker.GrpcPms.DeleteRole(c.Ctx(), &pmsv1.DeleteRoleReq{
		RoleId:     id,
		OperateUid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// RoleMemberList (0613) 获取某个role的成员列表
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.GetRoleMemberListRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func RoleMemberList(c *bffcore.Context) {
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}
	resp, err := invoker.GrpcPms.GetRoleMemberList(c.Ctx(), &pmsv1.GetRoleMemberListReq{
		OperateUid: c.Uid(),
		RoleId:     roleId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type RolePermissionResp struct {
	List []*commonv1.PmsItem `json:"list"`
}

// RolePermission (0613) 获取某个role的权限列表
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.GetRolePermissionRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func RolePermission(c *bffcore.Context) {
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}
	resp, err := invoker.GrpcPms.GetRolePermission(c.Ctx(), &pmsv1.GetRolePermissionReq{
		OperateUid: c.Uid(),
		RoleId:     roleId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(RolePermissionResp{
		List: resp.List,
	})
}

type RoleSpaceGroupPermissionRes struct {
	List []*pmsv1.SpaceGroupPmsItem `json:"list"`
}

// RoleSpaceGroupPermission (0613) 获取某个role的权限列表
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.GetRolePermissionRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func RoleSpaceGroupPermission(c *bffcore.Context) {
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}
	resp, err := invoker.GrpcPms.GetRolePermission(c.Ctx(), &pmsv1.GetRolePermissionReq{
		OperateUid: c.Uid(),
		RoleId:     roleId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(RoleSpaceGroupPermissionRes{
		List: resp.SpaceGroupList,
	})
}

type RoleSpacePermissionRes struct {
	List []*pmsv1.SpacePmsItem `json:"list"`
}

// RoleSpacePermission (0613) 获取某个role的权限列表
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.GetRolePermissionRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func RoleSpacePermission(c *bffcore.Context) {
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}
	resp, err := invoker.GrpcPms.GetRolePermission(c.Ctx(), &pmsv1.GetRolePermissionReq{
		OperateUid: c.Uid(),
		RoleId:     roleId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(RoleSpacePermissionRes{
		List: resp.SpaceList,
	})
}

type RoleSpaceGroupInitPermissionRes struct {
	List []*commonv1.PmsItem `json:"list"`
}

// RoleSpaceGroupInitPermission (0707) role的space group的初始权限点
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.RoleSpaceGroupInitPermissionRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func RoleSpaceGroupInitPermission(c *bffcore.Context) {
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}

	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid cant empty", nil)
		return
	}

	resp, err := invoker.GrpcPms.GetInitActionOptionPermission(c.Ctx(), &pmsv1.GetInitActionOptionPermissionReq{
		OperateUid: c.Uid(),
		Guid:       guid,
		Type:       commonv1.CMN_GUID_SPACE_GROUP,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(RoleSpaceGroupInitPermissionRes{
		List: resp.List,
	})
}

// RoleSpaceInitPermission (0707) role的space的初始权限点
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.RoleSpaceInitPermissionRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func RoleSpaceInitPermission(c *bffcore.Context) {
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}

	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid cant empty", nil)
		return
	}

	resp, err := invoker.GrpcPms.GetInitActionOptionPermission(c.Ctx(), &pmsv1.GetInitActionOptionPermissionReq{
		OperateUid: c.Uid(),
		Guid:       guid,
		Type:       commonv1.CMN_GUID_SPACE,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type PutRolePermissionRequest struct {
	List []*commonv1.PmsItem `json:"list"`
}

// PutRolePermission (0613) 修改某个role权限点
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.PutRolePermissionRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func PutRolePermission(c *bffcore.Context) {
	var req PutRolePermissionRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}

	resp, err := invoker.GrpcPms.PutRolePermission(c.Ctx(), &pmsv1.PutRolePermissionReq{
		OperateUid: c.Uid(),
		RoleId:     roleId,
		List:       req.List,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type PutRoleSpacePermissionRequest struct {
	List []*pmsv1.SpacePmsItem `json:"list"`
}

// PutRoleSpacePermission (0613) 修改某个role权限点
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.PutRoleSpacePermissionReq}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func PutRoleSpacePermission(c *bffcore.Context) {
	var req PutRoleSpacePermissionRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}

	resp, err := invoker.GrpcPms.PutRoleSpacePermission(c.Ctx(), &pmsv1.PutRoleSpacePermissionReq{
		OperateUid: c.Uid(),
		RoleId:     roleId,
		List:       req.List,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type PutRoleSpaceGroupPermissionRequest struct {
	List []*pmsv1.SpaceGroupPmsItem `json:"list"`
}

// PutRoleSpaceGroupPermission (0613) 修改某个role权限点
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.PutRoleSpaceGroupPermissionReq}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func PutRoleSpaceGroupPermission(c *bffcore.Context) {
	var req PutRoleSpaceGroupPermissionRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}

	resp, err := invoker.GrpcPms.PutRoleSpaceGroupPermission(c.Ctx(), &pmsv1.PutRoleSpaceGroupPermissionReq{
		OperateUid: c.Uid(),
		RoleId:     roleId,
		List:       req.List,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type CreateRoleMemberRequest struct {
	Uids []int64 `json:"uids"`
}

// CreateRoleMember (0613) 添加某个role成员
// @Tags Pms
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.CreateRoleMemberRes}
// @Failure 200 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func CreateRoleMember(c *bffcore.Context) {
	var req CreateRoleMemberRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}

	resp, err := invoker.GrpcPms.CreateRoleMember(c.Ctx(), &pmsv1.CreateRoleMemberReq{
		Uids:       req.Uids,
		OperateUid: c.Uid(),
		RoleId:     roleId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type DeleteRoleMemberRequest struct {
	Uids []int64 `json:"uids"`
}

// DeleteRoleMember 管理员基本数据
func DeleteRoleMember(c *bffcore.Context) {
	var req DeleteRoleMemberRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", nil)
		return
	}
	roleId := cast.ToInt64(c.Param("roleId"))
	if roleId == 0 {
		c.JSONE(1, "role id cant empty", nil)
		return
	}

	resp, err := invoker.GrpcPms.DeleteRoleMember(c.Ctx(), &pmsv1.DeleteRoleMemberReq{
		Uids:       req.Uids,
		OperateUid: c.Uid(),
		RoleId:     roleId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}
