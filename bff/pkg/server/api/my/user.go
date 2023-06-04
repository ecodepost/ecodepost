package my

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	countv1 "ecodepost/pb/count/v1"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/spf13/cast"
)

type UserInfo struct {
	Uid int64 `json:"uid"`
	// 昵称
	Nickname string `json:"nickname,omitempty"`
	// 头像
	Avatar string `json:"avatar,omitempty"`
	// 用户名称
	Name string `json:"name,omitempty"`
	// 有资格创建社区的审核状态
	CmtIdentifyStatus commonv1.AUDIT_STATUS `json:"cmtIdentifyStatus"`
}

type OauthUserInfoRes struct {
	User UserInfo `json:"user"`
}

// OauthUserInfo 获取当前登录用户Handler
func OauthUserInfo(c *bffcore.Context) {
	if c.Uid() == 0 {
		c.JSONOK(OauthUserInfoRes{
			User: UserInfo{
				Uid: 0,
			},
		})
		return
	}

	userOauthInfo, err := invoker.GrpcUser.OauthInfo(c.Ctx(), &userv1.OauthInfoReq{Uid: c.Uid()})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	c.JSONOK(OauthUserInfoRes{
		User: UserInfo{
			Uid:               c.Uid(),
			Nickname:          userOauthInfo.Nickname,
			Avatar:            userOauthInfo.Avatar,
			Name:              userOauthInfo.Name,
			CmtIdentifyStatus: userOauthInfo.GetCmtIdentifyStatus(),
		},
	})
}

type CommunityChooseRequest struct {
	Guid string `json:"guid"`
}
type BlockListReq struct {
	Pagination commonv1.Pagination `json:"pagination"`
}

func FollowingCreate(c *bffcore.Context) {
	uid := cast.ToInt64(c.Param("uid"))
	if uid <= 0 {
		c.EgoJsonI18N(errcodev1.ErrUidEmpty())
		return
	}
	_, err := invoker.GrpcCount.Set(c.Ctx(), &countv1.SetReq{
		Biz:  commonv1.CMN_BIZ_USER,
		Fid:  cast.ToString(c.Uid()),
		Tid:  cast.ToString(uid),
		Act:  commonv1.CNT_ACT_FOLLOW,
		Acti: commonv1.CNT_ACTI_ADD,
		Ip:   c.ClientIP(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func FollowingDelete(c *bffcore.Context) {
	uid := cast.ToInt64(c.Param("uid"))
	if uid <= 0 {
		c.EgoJsonI18N(errcodev1.ErrUidEmpty())
		return
	}
	_, err := invoker.GrpcCount.Set(c.Ctx(), &countv1.SetReq{
		Biz:  commonv1.CMN_BIZ_USER,
		Fid:  cast.ToString(c.Uid()),
		Tid:  cast.ToString(uid),
		Act:  commonv1.CNT_ACT_FOLLOW,
		Acti: commonv1.CNT_ACTI_SUB,
		Ip:   c.ClientIP(),
	})

	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
