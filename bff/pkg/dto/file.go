package dto

import (
	"encoding/json"

	commonv1 "ecodepost/pb/common/v1"
	userv1 "ecodepost/pb/user/v1"
)

type FileShow struct {
	// 文章GUID
	Guid string `json:"guid,omitempty"`
	// 文章标题
	Name string `json:"name,omitempty"`
	// 用户uid
	Uid int64 `json:"uid,omitempty"`
	// 用户昵称
	Nickname string `json:"nickname,omitempty"`
	// 用户头像
	Avatar string `json:"avatar,omitempty"`
	// 创建时间
	Ctime int64 `json:"ctime,omitempty"`
	// 回复总数
	CntComment int64 `json:"cntComment,omitempty"`
	// 查看总数
	CntView int64 `json:"cntView,omitempty"`
	// 收藏总数
	CntCollect int64 `json:"cntCollect,omitempty"`
	// headImage
	HeadImage string `json:"headImage,omitempty"`
	// 空间Guid
	SpaceGuid string `json:"spaceGuid,omitempty"`
	// 是否有readMore
	IsReadMore int32 `json:"isReadMore,omitempty"`
	// 是否允许创建评论
	IsAllowCreateComment int32 `json:"isAllowCreateComment,omitempty"`
	// 是否置顶
	IsSiteTop int32 `json:"isSiteTop,omitempty"`
	// 是否推荐
	IsRecommend int32 `json:"isRecommend,omitempty"`
	// 文档格式
	Format commonv1.FILE_FORMAT `json:"format,omitempty"`
	// emoji list
	EmojiList []*commonv1.EmojiInfo `json:"emojiList,omitempty"`
	// 内容
	Content json.RawMessage `json:"content"`
	//Content string `json:"content"`
	// ip定位的地址, 精确到省
	IpLocation string `json:"ipLocation,omitempty"`
}

type FilePermission struct {
	// 是否可以编辑
	IsAllowWrite bool `json:"isAllowWrite,omitempty"`
	// 是否可以删除
	IsAllowDelete bool `json:"isAllowDelete,omitempty"`
	// 是否可以置顶
	IsAllowSiteTop bool `json:"isAllowSiteTop,omitempty"`
	// 是否可以推荐
	IsAllowRecommend bool `json:"isAllowRecommend,omitempty"`
	// 是否可以打开评论或者关闭评论
	IsAllowSetComment bool `json:"isAllowSetComment,omitempty"`
	// 是否可以创建评论
	IsAllowCreateComment bool `json:"isAllowCreateComment,omitempty"`
}

type MetaAnswer struct {
	QuestionGuid  string          `json:"questionGuid"`
	QuestionTitle string          `json:"questionTitle"`
	AnswerAuthor  userv1.UserInfo `json:"answerAuthor"`
}
