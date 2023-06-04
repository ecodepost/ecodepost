package file

import (
	"encoding/json"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	filev1 "ecodepost/pb/file/v1"
	statv1 "ecodepost/pb/stat/v1"
	"github.com/samber/lo"
)

// GetInfo 获取文档内容
func GetInfo(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	res, err := invoker.GrpcFile.GetShowInfo(c.Ctx(), &filev1.GetShowInfoReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(dto.FileShow{
		Guid:                 res.GetFile().GetGuid(),
		Name:                 res.GetFile().GetName(),
		Uid:                  res.GetFile().GetUid(),
		Nickname:             res.GetFile().GetNickname(),
		Avatar:               res.GetFile().GetAvatar(),
		Ctime:                res.GetFile().GetCtime(),
		CntComment:           res.GetFile().GetCntComment(),
		CntView:              res.GetFile().GetCntView(),
		CntCollect:           res.GetFile().GetCntCollect(),
		HeadImage:            res.GetFile().GetHeadImage(),
		SpaceGuid:            res.GetFile().GetSpaceGuid(),
		IsAllowCreateComment: res.GetFile().GetIsAllowCreateComment(),
		IsSiteTop:            res.GetFile().GetIsSiteTop(),
		IsRecommend:          res.GetFile().GetIsRecommend(),
		Format:               res.GetFile().GetFormat(),
		EmojiList:            res.GetFile().GetEmojiList(),
		IpLocation:           res.GetFile().GetIpLocation(),
		Content:              json.RawMessage(res.GetFile().GetContent()),
	})
}

type StatReq struct {
	FileGuids []string `form:"fileGuids"`
}

type StatRes struct {
	EmojiList   []dto.MyEmojiInfo   `json:"emojiList"`
	CollectList []dto.MyCollectInfo `json:"collectList"`
}

// Stat 根据用户uid，文件guids，获取对应的状态数据
func Stat(c *bffcore.Context) {
	if c.Uid() == 0 {
		c.JSONOK(StatRes{
			EmojiList:   []dto.MyEmojiInfo{},
			CollectList: []dto.MyCollectInfo{},
		})
		return
	}

	req := StatReq{}
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "错误", err)
		return
	}
	if len(req.FileGuids) > 100 {
		c.JSONE(1, "错误", "fileGuids最多100个")
		return
	}

	emojiList, err := invoker.GrpcFile.MyEmojiListByFileGuids(c.Ctx(), &filev1.MyEmojiListByFileGuidsReq{
		Uid:   c.Uid(),
		Guids: req.FileGuids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	dtoEmoji := lo.Map(emojiList.List, func(t *filev1.MyEmojiInfo, i int) dto.MyEmojiInfo {
		return dto.MyEmojiInfo{
			Guid: t.Guid,
			List: t.List,
		}
	})

	collectList, err := invoker.GrpcStat.MyCollectionListByFileGuids(c.Ctx(), &statv1.MyCollectionListByFileGuidsReq{
		Uid:       c.Uid(),
		FileGuids: req.FileGuids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	dtoCollect := lo.Map(collectList.List, func(t *statv1.CollectionInfo, i int) dto.MyCollectInfo {
		isCollect := false
		if t.Id > 0 {
			isCollect = true
		}
		return dto.MyCollectInfo{
			Guid:      t.BizGuid,
			IsCollect: isCollect,
		}
	})
	c.JSONOK(&StatRes{
		EmojiList:   dtoEmoji,
		CollectList: dtoCollect,
	})
}

type IncreaseEmojiRequest struct {
	EmojiId int32 `json:"emojiId" binding:"required" label:"文档ID"`
}

// IncreaseEmoji 传入一个emoji id，点赞
// @Tags Emoji
// @Success 200 {object} bffcore.Res{}
func IncreaseEmoji(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}

	var req IncreaseEmojiRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcFile.CreateEmoji(c.Ctx(), &filev1.CreateEmojiReq{
		Uid:  c.Uid(),
		Guid: guid,
		V:    req.EmojiId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type DecreaseEmojiRequest struct {
	EmojiId int32 `json:"emojiId" binding:"required" label:"文档ID"`
}

// DecreaseEmoji 传入一个emoji id，去掉点赞
// @Tags Emoji
// @Success 200 {object} bffcore.Res{}
func DecreaseEmoji(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}

	var req DecreaseEmojiRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcFile.DeleteEmoji(c.Ctx(), &filev1.DeleteEmojiReq{
		Uid:  c.Uid(),
		Guid: guid,
		V:    req.EmojiId,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

// Permission 权限
func Permission(c *bffcore.Context) {
	if c.Uid() == 0 {
		c.JSONOK(dto.FilePermission{
			IsAllowWrite:         false,
			IsAllowDelete:        false,
			IsAllowSiteTop:       false,
			IsAllowRecommend:     false,
			IsAllowSetComment:    false,
			IsAllowCreateComment: false,
		})
		return
	}

	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcFile.Permission(c.Ctx(), &filev1.PermissionReq{
		Uid:      c.Uid(),
		FileGuid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(dto.FilePermission{
		IsAllowWrite:         resp.GetIsAllowWrite(),
		IsAllowDelete:        resp.GetIsAllowDelete(),
		IsAllowSiteTop:       resp.GetIsAllowSiteTop(),
		IsAllowRecommend:     resp.GetIsAllowRecommend(),
		IsAllowSetComment:    resp.GetIsAllowSetComment(),
		IsAllowCreateComment: resp.GetIsAllowCreateComment(),
	})
}

type PermissionListReq struct {
	Guids []string `form:"guids"`
}

// PermissionList 权限
func PermissionList(c *bffcore.Context) {
	var req PermissionListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}

	if c.Uid() == 0 {
		output := make([]filev1.PermissionRes, 0)
		for _, guid := range req.Guids {
			output = append(output, filev1.PermissionRes{
				Guid: guid,
			})
		}
		c.JSONOK(output)
		return
	}

	if len(req.Guids) == 0 {
		c.EgoJsonI18N(errcodev1.ErrFileGuidEmpty())
		return
	}

	res, err := invoker.GrpcFile.PermissionList(c.Ctx(), &filev1.PermissionListReq{
		Uid:      c.Uid(),
		FileGuid: req.Guids,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res)
}

type BatchGetUrlsReq struct {
	SpaceGuid   string               `form:"spaceGuid" json:"spaceGuid" binding:"required"`
	ContentKeys []string             `form:"contentKeys" json:"contentKeys" binding:"required"`
	UploadType  commonv1.CMN_UP_TYPE `form:"uploadType" json:"uploadType" binding:"required"`
}

type BatchGetUrlsResWrapper struct {
	Urls map[string]string `json:"urls"`
}
