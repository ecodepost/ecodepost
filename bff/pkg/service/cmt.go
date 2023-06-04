package service

import (
	"fmt"
	"strings"
	"time"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	communityv1 "ecodepost/pb/community/v1"
	userv1 "ecodepost/pb/user/v1"
)

type cmt struct{}

// GetFirstVisitUrl 获取第一次访问的URL
func (*cmt) GetFirstVisitUrl(c *bffcore.Context) (pageOptionUrl string, err error) {
	// 获取社区的page options
	homeOption, err := invoker.GrpcCommunity.GetHomeOption(c.Ctx(), &communityv1.GetHomeOptionReq{
		Uid: c.Uid(),
	})
	if err != nil {
		return
	}
	// 如果没有登录
	if c.Uid() == 0 {
		pageOptionUrl, err = pageOption(homeOption.DefaultPageByNotLogin)
		if err != nil {
			return
		}
		return
	}
	// 如果登录了
	profileInfo, err := invoker.GrpcUser.ProfileInfo(c.Ctx(), &userv1.ProfileInfoReq{
		Uid: c.Uid(),
	})
	if err != nil {
		return
	}
	// 我们认为大于7天不在是新用户
	if profileInfo.RegisterTime > time.Now().Unix()+86400*7 {
		pageOptionUrl, err = pageOption(homeOption.DefaultPageByNewUser)
		if err != nil {
			return
		}
		return
	}
	pageOptionUrl, err = pageOption(homeOption.DefaultPageByLogin)
	if err != nil {
		return
	}
	return
}

// GetInfo 获取社区信息
func (c *cmt) GetInfo(ctx *bffcore.Context, isUpdateLastActiveCmt bool) (res dto.CmtInfo, err error) {
	homeRes, e := invoker.GrpcCommunity.Home(ctx.Ctx(), &communityv1.HomeReq{})
	if e != nil {
		return res, e
	}
	firstVisitUrl, e := c.GetFirstVisitUrl(ctx)
	if e != nil {
		return res, e
	}

	res = dto.CmtInfo{
		Name:        homeRes.GetName(),
		Description: homeRes.GetDescription(),
		Logo:        homeRes.GetLogo(),
		Ctime:       homeRes.GetCtime(),
		//IsSetHome:      homeRes.GetIsSetHome(),
		FirstVisitUrl:  firstVisitUrl,
		Access:         homeRes.GetAccess(),
		GongxinbuBeian: homeRes.GetGongxinbuBeian(),
	}
	return
}

const spaceUrl = "/s/%s"

func pageOption(pageInfo string) (page string, err error) {
	pageArr := strings.Split(pageInfo, "_")
	if len(pageArr) != 2 {
		return "", fmt.Errorf("page arr length not eq 2")
	}
	if pageArr[0] == "sys" {
		if pageArr[1] == "home" {
			return fmt.Sprintf("/home"), nil
		}
	}
	if pageArr[0] == "space" {
		return fmt.Sprintf(spaceUrl, pageArr[1]), nil
	}
	return "", fmt.Errorf("not suppport space option")
}
