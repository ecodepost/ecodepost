package pms

import (
	"errors"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	errcodev1 "ecodepost/pb/errcode/v1"
	"github.com/gotomicro/ego/core/eerrors"

	pmsv1 "ecodepost/pb/pms/v1"

	"github.com/spf13/cast"
)

type MemberInfo struct {
	Uid      int64  `json:"uid,omitempty"`      // 用户uid
	Nickname string `json:"nickname,omitempty"` // 用户昵称
	Avatar   string `json:"avatar,omitempty"`   // 用户头像
	Email    string `json:"email,omitempty"`    // 用户email
	Ctime    int64  `json:"ctime,omitempty"`    // 创建时间
}

type ManagerMemberListResp struct {
	List []MemberInfo `json:"list"` // 成员列表
}

// ManagerMemberList 查看管理员的成员列表
// @Tags Pms
// @Param managerType  Path  string  true  "字符串类型，superAdmin为超级管理员，admin为管理员" Enums(superAdmin,admin)
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=ManagerMemberListResp}
// @Failure 400 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func ManagerMemberList(c *bffcore.Context) {
	resp, err := invoker.GrpcPms.GetManagerMemberList(c.Ctx(), &pmsv1.GetManagerMemberListReq{
		OperateUid: c.Uid(),
	})
	egoErr := eerrors.FromError(err)
	if err != nil && !errors.Is(egoErr, errcodev1.ErrNotFound()) {
		c.EgoJsonI18N(err)
		return
	}

	if errors.Is(egoErr, errcodev1.ErrNotFound()) {
		c.JSONOK([]struct{}{})
		return
	}

	output := make([]MemberInfo, 0)
	for _, value := range resp.GetList() {
		output = append(output, MemberInfo{
			Uid:      value.GetUid(),
			Nickname: value.GetNickname(),
			Avatar:   value.GetAvatar(),
			Ctime:    value.GetCtime(),
		})
	}

	c.JSONOK(ManagerMemberListResp{
		List: output,
	})
}

type CreateManagerMemberReq struct {
	Uids []int64 `json:"uids"`
}

// CreateManagerMember (0612) 添加权限成员
// @Tags Pms
// @Param managerType  Path  string  true  "字符串类型，superAdmin为超级管理员，admin为管理员" Enums(superAdmin,admin)
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.CreateManagerMemberRes}
// @Failure 400 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func CreateManagerMember(c *bffcore.Context) {
	var req CreateManagerMemberReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid param: "+err.Error(), err)
		return
	}

	resp, err := invoker.GrpcPms.CreateManagerMember(c.Ctx(), &pmsv1.CreateManagerMemberReq{
		Uids:       req.Uids,
		OperateUid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// DeleteManagerMember (0612) 删除权限成员
// @Tags Pms
// @Param managerType  Path  string  true  "字符串类型，superAdmin为超级管理员，admin为管理员" Enums(superAdmin,admin)
// @Success 200 {object} bffcore.Res{code=int,msg=string,data=pmsv1.DeleteManagerMemberRes}
// @Failure 400 {object} bffcore.Res{code=int,msg=string,data=string} "" bffcore.Res{code:1,msg:"fail",data="fail"}
func DeleteManagerMember(c *bffcore.Context) {
	uid := cast.ToInt64(c.Param("uid"))
	if uid == 0 {
		c.JSONE(1, "用户不能为空", nil)
		return
	}
	resp, err := invoker.GrpcPms.DeleteManagerMember(c.Ctx(), &pmsv1.DeleteManagerMemberReq{
		Uid:        uid,
		OperateUid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}
