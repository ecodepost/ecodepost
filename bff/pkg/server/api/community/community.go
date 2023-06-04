package community

import (
	"errors"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/service"
	commonv1 "ecodepost/pb/common/v1"
	communityv1 "ecodepost/pb/community/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	pmsv1 "ecodepost/pb/pms/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eerrors"
)

type UpdateRequest struct {
	Name       *string             `json:"name" binding:"required" label:"名称"`
	Logo       *string             `json:"logo" label:"图标"`
	Visibility *commonv1.CMN_VISBL `json:"visibility"`
}

func Update(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "invalid param", nil)
		return
	}
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid param,"+err.Error(), err)
		return
	}
	c.InjectGuid(guid)
	// 更新社区
	res, err := invoker.GrpcCommunity.Update(c.Ctx(), &communityv1.UpdateReq{
		Name:       req.Name,
		Uid:        c.Uid(),
		Logo:       req.Logo,
		Visibility: req.Visibility,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res)
}

type UpdatePriceRequest struct {
	IsSetPrice        bool  `json:"isSetPrice" label:"图片"`
	AnnualPrice       int64 `json:"annualPrice" binding:"required" label:"标题"`
	AnnualOriginPrice int64 `json:"annualOriginPrice" label:"内容"`
}

type UpdateBannerRequest struct {
	Img   string `json:"img" label:"图片"`
	Title string `json:"title" binding:"required" label:"标题"`
}

func UpdateBanner(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "invalid param", nil)
		return
	}
	var req UpdateBannerRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid param,"+err.Error(), err)
		return
	}
	c.InjectGuid(guid)

	// 更新社区
	res, err := invoker.GrpcCommunity.PutHomeOption(c.Ctx(), &communityv1.PutHomeOptionReq{
		Uid:         c.Uid(),
		BannerImg:   &req.Img,
		BannerTitle: &req.Title,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res)
}

type MemberListReq struct {
	CurrentPage int32 `form:"currentPage"` // 当前页数
}

// ListLogos 社区推荐logo列表
func ListLogos(c *bffcore.Context) {
	c.JSONOK(dto.ListRes{
		List: econf.GetStringSlice("recommend.communityLogos"),
	})
}

// ListCovers 社区推荐封面
func ListCovers(c *bffcore.Context) {
	c.JSONOK(dto.ListRes{
		List: econf.GetStringSlice("recommend.communityCovers"),
	})
}

func Managers(c *bffcore.Context) {
	res, err := invoker.GrpcPms.GetManagerMemberList(c.Ctx(), &pmsv1.GetManagerMemberListReq{})
	egoErr := eerrors.FromError(err)
	if err != nil && !errors.Is(egoErr, errcodev1.ErrNotFound()) {
		c.EgoJsonI18N(err)
		return
	}

	if errors.Is(egoErr, errcodev1.ErrNotFound()) {
		c.JSONOK([]struct{}{})
		return
	}

	c.JSONOK(res)
}

// SearchMember 搜索社区成员
func SearchMember(c *bffcore.Context) {
	var req SearchMemberReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectSpc(req.SpaceGuid)
	res, err := invoker.GrpcUser.Search(c.Ctx(), &userv1.SearchReq{
		Nickname: req.Keyword,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 如果为空，直接返回空数据
	if len(res.List) == 0 {
		c.JSONListPage([]*userv1.UserInfo{}, nil)
		return
	}

	c.JSONListPage(res.List, nil)
}

// MemberList @LHP @2022-10-26
func MemberList(c *bffcore.Context) {
	var req MemberListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid param,"+err.Error(), err)
		return
	}
	// 社区成员
	resp, err := invoker.GrpcUser.ListPage(c.Ctx(), &userv1.ListPageReq{
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

type CustomDomain struct {
	Host    string
	CmtGuid string
}

// Detail 社区的详细信息
func Detail(c *bffcore.Context) {
	res := dto.CmtDetail{}
	// 获取该信息的时候，需要更新一次用户表，标记用户上一次访问社区的信息
	cmtInfo, err := service.Cmt.GetInfo(c, true)
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	homeOptionRes, err := invoker.GrpcCommunity.GetHomeOption(c.Ctx(), &communityv1.GetHomeOptionReq{})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	cmtInfo.IsSetHome = homeOptionRes.IsSetHome
	res.CmtInfo = cmtInfo

	themeRes, err := invoker.GrpcCommunity.GetTheme(c.Ctx(), &communityv1.GetThemeReq{})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	res.CmtTheme = dto.CmtTheme{
		IsCustom:          themeRes.IsCustom,
		ThemeName:         themeRes.ThemeName,
		CustomColor:       []byte(themeRes.CustomColor),
		DefaultAppearance: themeRes.DefaultAppearance,
	}

	if c.Uid() == 0 {
		res.UserInfo = dto.CmtUserInfo{IsLogin: false}
		c.JSONOK(res)
		return
	}

	// 查询社区授权数据
	permissionRes, err := invoker.GrpcPms.CommunityPermission(c.Ctx(), &pmsv1.CommunityPermissionReq{
		Uid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	output := dto.CmtPermission{
		IsAllowManageCommunity:  permissionRes.GetIsAllowManageCommunity(),
		IsAllowCreateSpaceGroup: permissionRes.GetIsAllowCreateSpaceGroup(),
		IsAllowCreateSpace:      permissionRes.GetIsAllowCreateSpace(),
		IsAllowUpgradeEdition:   permissionRes.GetIsAllowUpgradeEdition(),
	}
	// 如果可以创建空间，那么获取下可以创建的应用列表
	if permissionRes.GetIsAllowCreateSpace() {
		appList := make([]dto.AppInfo, 0)
		appList = append(appList, dto.AppInfo{
			AppId: commonv1.CMN_APP_ARTICLE,
			Name:  "帖子",
		})
		appList = append(appList, dto.AppInfo{
			AppId: commonv1.CMN_APP_QA,
			Name:  "问答",
		})
		appList = append(appList, dto.AppInfo{
			AppId: commonv1.CMN_APP_COLUMN,
			Name:  "栏目",
		})
		appList = append(appList, dto.AppInfo{
			AppId: commonv1.CMN_APP_LINK,
			Name:  "链接",
		})

		output.AppList = appList
	}

	res.Permission = output
	c.JSONOK(res)
	return

}

type SearchMemberReq struct {
	Keyword   string `form:"keyword" binding:"required" label:"关键字"` // 当前页数
	SpaceGuid string `form:"spaceGuid"`                              // 空间guid，可以为空
	bffcore.Pagination
}
