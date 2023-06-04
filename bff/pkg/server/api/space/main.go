package space

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	columnv1 "ecodepost/pb/column/v1"

	commonv1 "ecodepost/pb/common/v1"
	spacev1 "ecodepost/pb/space/v1"
)

type TreeChangeSortRequest struct {
	Type                 commonv1.CMN_GUID `json:"type"`
	TargetSpaceGuid      *string           `json:"targetSpaceGuid"`
	TargetSpaceGroupGuid *string           `json:"targetSpaceGroupGuid"`
	ParentSpaceGroupGuid *string           `json:"parentSpaceGroupGuid"`
	DropPosition         *string           `json:"dropPosition"`
	SpaceGuid            string            `json:"spaceGuid" label:"组ID"`
	SpaceGroupGuid       string            `json:"spaceGroupGuid" label:"组ID"`
}

func TreeChangeSort(c *bffcore.Context) {
	var req TreeChangeSortRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectSpc(req.SpaceGuid)
	if req.Type == commonv1.CMN_GUID_SPACE {
		_, err := invoker.GrpcSpace.ChangeSpaceSort(c.Ctx(), &spacev1.ChangeSpaceSortReq{
			OperateUid:           c.Uid(),
			SpaceGuid:            req.SpaceGuid,
			TargetSpaceGuid:      req.TargetSpaceGuid,
			DropPosition:         req.DropPosition,
			ParentSpaceGroupGuid: req.ParentSpaceGroupGuid,
		})
		if err != nil {
			c.EgoJsonI18N(err)
			return
		}
		c.JSONOK()
	} else if req.Type == commonv1.CMN_GUID_SPACE_GROUP {
		_, err := invoker.GrpcSpace.ChangeSpaceGroupSort(c.Ctx(), &spacev1.ChangeSpaceGroupSortReq{
			OperateUid:           c.Uid(),
			SpaceGroupGuid:       req.SpaceGroupGuid,
			TargetSpaceGroupGuid: *req.TargetSpaceGroupGuid,
			DropPosition:         *req.DropPosition,
		})
		if err != nil {
			c.EgoJsonI18N(err)
			return
		}
		c.JSONOK()
	} else {
		c.JSONE(1, "not found type", nil)
	}
}

type TreeChangeGroupRequest struct {
	SpaceGuid           string `json:"spaceGuid" label:"组ID"`
	AfterSpaceGroupGuid string `json:"afterSpaceGroupGuid"`
}

func TreeChangeGroup(c *bffcore.Context) {
	var req TreeChangeGroupRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}

	if req.SpaceGuid == "" {
		c.JSONE(1, "名称不能为空", nil)
		return
	}
	if req.AfterSpaceGroupGuid == "" {
		c.JSONE(1, "名称不能为空", nil)
		return
	}
	c.InjectSpc(req.SpaceGuid)

	resp, err := invoker.GrpcSpace.UpdateSpace(c.Ctx(), &spacev1.UpdateSpaceReq{
		OperateUid:     c.Uid(),
		SpaceGuid:      req.SpaceGuid,
		SpaceGroupGuid: &req.AfterSpaceGroupGuid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type CreateRequest struct {
	Name            string              `json:"name" binding:"required" label:"名称"`          // 空间名称
	SpaceGroupGuid  string              `json:"spaceGroupGuid" binding:"required" label:"组"` // 空间分组
	Icon            string              `json:"icon"`                                        // 图标
	SpaceType       commonv1.CMN_APP    `json:"spaceType"`                                   // 空间类型
	SpaceLayout     commonv1.SPC_LAYOUT `json:"spaceLayout"`                                 // 空间布局
	Visibility      commonv1.CMN_VISBL  `json:"visibility"`                                  // visibility
	ColumnAuthorUid *int64              `json:"columnAuthorUid"`
	Link            string              `json:"link"`
	Cover           string              `json:"cover"`
}

func Create(c *bffcore.Context) {
	var req CreateRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}

	// 设置spaceType为文章作为缺省值
	if req.SpaceType == 0 {
		req.SpaceType = commonv1.CMN_APP_ARTICLE
	}
	// 根据spaceType类型，设置其spaceLayout缺省值
	if req.SpaceLayout == 0 {
		switch req.SpaceType {
		case commonv1.CMN_APP_ARTICLE:
			req.SpaceLayout = commonv1.SPC_LAYOUT_ARTICLE_FEED
		case commonv1.CMN_APP_COLUMN:
			req.SpaceLayout = commonv1.SPC_LAYOUT_ARTICLE_TREE
		case commonv1.CMN_APP_QA:
			// 不做任何操作
		}
	}
	resp, err := invoker.GrpcSpace.CreateSpace(c.Ctx(), &spacev1.CreateSpaceReq{
		OperateUid:     c.Uid(),
		Name:           req.Name,
		SpaceGroupGuid: req.SpaceGroupGuid,
		Icon:           req.Icon,
		SpaceType:      req.SpaceType,
		SpaceLayout:    req.SpaceLayout,
		Visibility:     req.Visibility,
		Link:           req.Link,
		Cover:          req.Cover,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 执行各种space其他业务操作
	switch req.SpaceType {
	case commonv1.CMN_APP_COLUMN:
		// 更新teachers
		_, err = invoker.GrpcColumn.CreateSpaceInfo(c.Ctx(), &columnv1.CreateSpaceInfoReq{
			SpaceGuid: resp.GetInfo().GetGuid(),
			Uid:       c.Uid(),
			AuthorUid: *req.ColumnAuthorUid,
		})
		if err != nil {
			c.EgoJsonI18N(err)
			return
		}
	}

	c.JSONOK(resp)
}

type CheckMembershipReq struct {
	Guids []string `form:"guids" binding:"required" label:"空间guid列表"`
}

type CheckMembershipRes []*spacev1.MemberStatus

func GetMemberStatus(c *bffcore.Context) {
	var req CheckMembershipReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid params, "+err.Error(), nil)
		return
	}
	res, err := invoker.GrpcSpace.GetMemberStatus(c.Ctx(), &spacev1.GetMemberStatusReq{
		Uid:        c.Uid(),
		SpaceGuids: req.Guids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res.List)
}

func Info(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid不能为空", nil)
		return
	}
	res, err := invoker.GrpcSpace.SpaceInfo(c.Ctx(), &spacev1.SpaceInfoReq{
		OperateUid: c.Uid(),
		SpaceGuid:  guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res.SpaceInfo)
}

type ColumnInfoResp struct {
	// GUID
	Guid string `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
	// 名称
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// 成员个数
	MemberCnt int64 `protobuf:"varint,9,opt,name=memberCnt,proto3" json:"memberCnt,omitempty"`
	// 空间分组guid
	SpaceGroupGuid string `protobuf:"bytes,10,opt,name=spaceGroupGuid,proto3" json:"spaceGroupGuid,omitempty"`
	// 空间简介或描述
	Desc string `protobuf:"bytes,15,opt,name=desc,proto3" json:"desc,omitempty"`
	// 空间封面
	Cover          string `protobuf:"bytes,17,opt,name=cover,proto3" json:"cover,omitempty"`
	AuthorUid      int64  `json:"authorUid"`
	AuthorNickname string `json:"authorNickname"`
	AuthorAvatar   string `json:"authorAvatar"`
}

func ColumnInfo(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid不能为空", nil)
		return
	}
	res, err := invoker.GrpcSpace.SpaceInfo(c.Ctx(), &spacev1.SpaceInfoReq{
		OperateUid: c.Uid(),
		SpaceGuid:  guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	resColumn, err := invoker.GrpcColumn.GetSpaceInfo(c.Ctx(), &columnv1.GetSpaceInfoReq{
		Uid:       c.Uid(),
		SpaceGuid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(ColumnInfoResp{
		Guid:           res.GetSpaceInfo().Guid,
		Name:           res.GetSpaceInfo().Name,
		Desc:           res.GetSpaceInfo().Desc,
		Cover:          res.GetSpaceInfo().Cover,
		AuthorUid:      resColumn.AuthorUid,
		AuthorNickname: resColumn.AuthorNickname,
		AuthorAvatar:   resColumn.AuthorAvatar,
	})
}

type UpdateRequest struct {
	// SpaceGroupGuid        string                 `json:"spaceGroupGuid" binding:"required" label:"组"`
	SpaceGroupGuid        *string                  `json:"spaceGroupGuid" label:"组"`
	Name                  *string                  `json:"name" label:"名称"`
	Icon                  *string                  `json:"icon" label:"图标"`
	ChargeType            *commonv1.SPC_CT         `json:"chargeType" label:"收费方式"`
	Price                 *int64                   `json:"price" label:"价格"`
	Desc                  *string                  `json:"desc" label:"简介或描述"`
	HeadImage             *string                  `json:"headImage" label:"头图"`
	Teachers              []string                 `json:"teachers" label:"教师guid列表"`
	Cover                 *string                  `json:"cover" label:"封面"`
	SpaceType             *commonv1.CMN_APP        `json:"spaceType"`             // 空间类型
	SpaceLayout           *commonv1.SPC_LAYOUT     `json:"spaceLayout"`           // 空间布局
	Visibility            *commonv1.CMN_VISBL      `json:"visibility"`            // visibility
	IsAllowReadMemberList *bool                    `json:"isAllowReadMemberList"` // 为零值表示不修改
	SpaceOptions          *[]*commonv1.SpaceOption `json:"spaceOptions"`          // 为零值表示不修改
	Topics                *[]string                `json:"topics"`
	Access                *commonv1.SPC_ACS        `json:"access"`
	Link                  *string                  `json:"link"`
	ColumnAuthorUid       *int64                   `json:"columnAuthorUid"`
	// Status                *commonv1.SPC_STATUS     `json:"status"`
}

func Update(c *bffcore.Context) {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid不能为空", nil)
		return
	}
	c.InjectSpc(guid)

	// 查询空间属性
	spaceInfo, err := invoker.GrpcSpace.SpaceInfo(c.Ctx(), &spacev1.SpaceInfoReq{
		OperateUid: c.Uid(),
		SpaceGuid:  guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 更新空间属性
	upSpaceReq := &spacev1.UpdateSpaceReq{
		OperateUid:            c.Uid(),
		SpaceGuid:             guid,
		Name:                  req.Name,
		Icon:                  req.Icon,
		SpaceType:             req.SpaceType,
		SpaceLayout:           req.SpaceLayout,
		Visibility:            req.Visibility,
		IsAllowReadMemberList: req.IsAllowReadMemberList,
		ChargeType:            req.ChargeType,
		OriginPrice:           req.Price,
		Price:                 req.Price,
		Desc:                  req.Desc,
		HeadImage:             req.HeadImage,
		Cover:                 req.Cover,
		Access:                req.Access,
		Link:                  req.Link,
		// Status:                req.Status,
	}
	if req.SpaceOptions != nil {
		upSpaceReq.SpaceOptions = *req.SpaceOptions
	}
	res, err := invoker.GrpcSpace.UpdateSpace(c.Ctx(), upSpaceReq)
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 执行各种space其他业务操作
	switch spaceInfo.SpaceInfo.SpaceType {
	case commonv1.CMN_APP_COLUMN:
		// 更新teachers
		_, err = invoker.GrpcColumn.UpdateSpaceInfo(c.Ctx(), &columnv1.UpdateSpaceInfoReq{
			SpaceGuid: guid,
			Uid:       c.Uid(),
			AuthorUid: *req.ColumnAuthorUid,
		})
		if err != nil {
			c.EgoJsonI18N(err)
			return
		}
	}
	c.JSONOK(res)
}

func Delete(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "guid不能为空", nil)
		return
	}
	c.InjectSpc(guid)
	resp, err := invoker.GrpcSpace.DeleteSpace(c.Ctx(), &spacev1.DeleteSpaceReq{
		SpaceGuid:  guid,
		OperateUid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

// type ChangeSortRequest struct {
//	AfterSpaceGuid string `json:"afterSpaceGuid"`
//	SpaceGuid      string `json:"spaceGuid" label:"组ID"`
// }
//
// func ChangeSort(c *bffcore.Context) {
//	var req ChangeSortRequest
//	if err := c.Bind(&req); err != nil {
//		c.JSONE(1, "参数错误", err)
//		return
//	}
//	resp, err := invoker.GrpcSpace.ChangeSpaceSort(c.Ctx(), &spacev1.ChangeSpaceSortReq{
//		OperateUid:     c.Uid(),
//		CmtGuid:        c.CmtGuid(),
//		SpaceGuid:      req.SpaceGuid,
//		AfterSpaceGuid: req.AfterSpaceGuid,
//	})
//	if err != nil {
//		c.EgoJsonI18N(err)
//		return
//	}
//	c.JSONOK(resp)
// }
