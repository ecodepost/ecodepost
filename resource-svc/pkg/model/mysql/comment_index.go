package mysql

import (
	"fmt"

	commentv1 "ecodepost/pb/comment/v1"
	commonv1 "ecodepost/pb/common/v1"
	userv1 "ecodepost/pb/user/v1"

	"gorm.io/gorm"
)

// CommentIndex 索引表
type CommentIndex struct {
	Id              int64             `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	Guid            string            `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:comment index guid"`
	SubjectId       int64             `gorm:"not null; default:0; comment:SubjectId"`
	Uid             int64             `gorm:"not null;" json:"uid"`          // 用户id
	ReplyToUid      int64             `gorm:"not null;" json:"replayUserId"` // 回复的id
	ReplyToGuid     string            `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:comment index guid"`
	ReplyToRootGuid string            `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:comment index guid"`
	BizType         commonv1.CMN_BIZ  // 业务类型
	BizGuid         string            `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:comment index guid"`
	ActionType      commonv1.FILE_ACT // 行动类型, 置顶、加精，提及，评论
	ActionGuid      string            `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:comment index guid"`
	CntChildComment int64             `gorm:"not null;default:0;comment:子评论数"` // 子评论总数
	CntStar         int32             `gorm:"not null;default:0;comment:点赞" json:"cntStar"`
	Status          int8              `gorm:"not null;index:idx_status;comment:0正常、1审核不通过，2审核中，3删除等" json:"status"`
	Ctime           int64             `gorm:"comment:创建时间" json:"ctime"`
	Utime           int64             `gorm:"comment:更新时间" json:"utime"` // 更新时间
	Dtime           int64             `gorm:"comment:删除时间;" json:"dtime"`
	Ip              string            `gorm:"type:varchar(15) not null;default:'';comment:IP"` // 后续可以考虑用unsigned int存储
}

func (CommentIndex) TableName() string {
	return "comment_index"
}

type CommentIndexStoreRedis struct {
	Id              int64 `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	Guid            string
	SubjectId       int64
	Uid             int64            `gorm:"not null;" json:"uid"` // 用户id
	BizType         commonv1.CMN_BIZ // 业务类型
	BizGuid         string
	ReplyToUid      int64 `gorm:"not null;" json:"replayUserId"` // 回复的id
	Ctime           int64 `gorm:"comment:创建时间" json:"ctime"`
	ActionGuid      string
	ActionType      commonv1.FILE_ACT // 行动类型, 置顶、加精，提及，评论
	ReplyToRootGuid string
	ReplyToGuid     string
	CntChildComment int64 `gorm:"not null;default:0;comment:子评论数"` // 子评论总数
	Ip              string
}

func (CommentIndexStoreRedis) TableName() string {
	return "comment_index"
}

// CommentIndexPage 只给部分数据
type CommentIndexPage struct {
	Guid string
}

func (CommentIndexPage) TableName() string {
	return "comment_index"
}

type CommentIndexList []*CommentIndex

func findUserByUid(list []*userv1.UserInfo, uid int64) (uInfo *userv1.UserInfo) {
	for _, u := range list {
		if u.Uid == uid {
			uInfo = u
		}
	}
	return uInfo
}

func (c CommentIndexStoreRedis) ToPBDetail(uReply *userv1.ListRes, content string) *commentv1.CommentDetail {
	var (
		userNickname  string
		userAvatar    string
		replyNickname string
		replyAvatar   string
		flag          bool
	)

	if c.Uid != 0 {
		// uInfo, flag := uReply.UserList[c.Uid]
		uInfo := findUserByUid(uReply.UserList, c.Uid)
		if flag {
			userNickname = uInfo.Nickname
			userAvatar = uInfo.Avatar
		}
	}

	// rInfo := &userv1.InfoRes{}
	isReply := 0
	// 回复需要返回回复头像的信息
	if c.ReplyToUid != 0 {
		rInfo := findUserByUid(uReply.UserList, c.ReplyToUid)
		if rInfo == nil {
			replyNickname = rInfo.Nickname
			replyAvatar = rInfo.Avatar
		}
		isReply = 1
	}

	return &commentv1.CommentDetail{
		CommentGuid:         c.Guid,
		BizGuid:             c.BizGuid,
		BizType:             c.BizType,
		Content:             content,
		ReplyToGuid:         c.ReplyToGuid,
		ReplyToRootGuid:     c.ReplyToRootGuid,
		Uid:                 c.Uid,
		UserNickname:        userNickname,
		UserAvatar:          userAvatar,
		ReplyToUid:          c.ReplyToUid,
		ReplyNickname:       replyNickname,
		ReplyAvatar:         replyAvatar,
		Ctime:               c.Ctime,
		IsReply:             int32(isReply),
		ActionGuid:          c.ActionGuid,
		ActionType:          c.ActionType,
		Children:            nil,
		CntChildComment:     c.CntChildComment,
		HasMoreChildComment: false,
		//IpLocation:          ipdb.IpDB.GetCountryProvice(c.Ip),
	}
}

// CommentInfoByField 获取评论信息
func CommentInfoByField(db *gorm.DB, field, guid string) (info CommentIndex, err error) {
	if err = db.Select(field).Where("guid = ?", guid).Find(&info).Error; err != nil {
		err = fmt.Errorf("CommentInfoByField fail, err: %w", err)
		return
	}
	return
}
