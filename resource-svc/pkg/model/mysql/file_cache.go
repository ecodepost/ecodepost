package mysql

import (
	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	questionv1 "ecodepost/pb/question/v1"
)

type FileCache struct {
	Guid                 string                `redis:"guid"`
	SpaceGuid            string                `redis:"sGuid"`
	ParentGuid           string                `redis:"pGuid"` // 只用于问答
	Name                 string                `redis:"n"`
	HeadImage            string                `redis:"hi"`
	Hash                 string                `redis:"h"`
	CreatedUid           int64                 `redis:"cUid"`
	UpdatedUid           int64                 `redis:"uUid"`
	Ctime                int64                 `redis:"ct"`
	Utime                int64                 `redis:"ut"`
	IsAllowCreateComment int32                 `redis:"isACC"`
	IsSiteTop            int32                 `redis:"isST"`
	IsRecommend          int32                 `redis:"isR"`
	FileFormat           int32                 `redis:"ff"`
	BizStatus            int32                 `redis:"bsts"`
	Content              string                `redis:"cont"`
	FileType             int32                 `redis:"ft"`
	Ip                   string                `redis:"ip"`
	EmojiList            []*commonv1.EmojiInfo `redis:"-"`
	Username             string                `redis:"-"`
	Nickname             string                `redis:"-"`
	Avatar               string                `redis:"-"`
	ContentUrl           string                `redis:"-"`
	CntComment           int64                 `redis:"-"`
	CntCollect           int64                 `redis:"-"`
	CntView              int64                 `redis:"-"`
	IsCollect            int32                 `redis:"-"`
	FileNode             commonv1.FILE_NODE    `redis:"-"`
	Icon                 string                `redis:"-"`
	Size                 int64                 `redis:"-"`
}

func (f *FileCache) ToMap() map[string]any {
	return map[string]any{
		"guid":  f.Guid,
		"sGuid": f.SpaceGuid,
		"pGuid": f.ParentGuid,
		"n":     f.Name,
		"hi":    f.HeadImage,
		"cUid":  f.CreatedUid,
		"uUid":  f.UpdatedUid,
		"ct":    f.Ctime,
		"ut":    f.Utime,
		"isACC": f.IsAllowCreateComment,
		"isST":  f.IsSiteTop,
		"isR":   f.IsRecommend,
		"ff":    f.FileFormat,
		"bsts":  f.BizStatus,
		"h":     f.Hash,
		"cont":  f.Content,
		"ft":    f.FileType,
		"ip":    f.Ip,
	}
}

type FileCaches []*FileCache

func (list FileCaches) ToArticlePb(layout commonv1.SPC_LAYOUT) []*articlev1.ArticleShow {
	output := make([]*articlev1.ArticleShow, 0)
	for _, value := range list {
		output = append(output, value.ToArticlePb(layout))
	}
	return output
}

func (list FileCaches) ToFileShowPb() []*commonv1.FileShow {
	output := make([]*commonv1.FileShow, 0)
	for _, value := range list {
		output = append(output, value.ToFileShowPb())
	}
	return output
}

func (f *FileCache) ToOneFileShowPb() *commonv1.FileShow {
	return &commonv1.FileShow{
		Guid:                 f.Guid,
		Name:                 f.Name,
		Uid:                  f.CreatedUid,
		Username:             f.Username,
		Nickname:             f.Nickname,
		Avatar:               f.Avatar,
		Ctime:                f.Ctime,
		CntComment:           f.CntComment,
		CntView:              f.CntView,
		CntCollect:           f.CntCollect,
		HeadImage:            f.HeadImage,
		SpaceGuid:            f.SpaceGuid,
		IsAllowCreateComment: f.IsAllowCreateComment,
		IsSiteTop:            f.IsSiteTop,
		IsRecommend:          f.IsRecommend,
		Format:               commonv1.FILE_FORMAT(f.FileFormat),
		EmojiList:            f.EmojiList,
		Content:              f.Content,
	}
}

func (f *FileCache) ToFileShowPb() *commonv1.FileShow {
	return &commonv1.FileShow{
		Guid:                 f.Guid,
		Name:                 f.Name,
		Uid:                  f.CreatedUid,
		Username:             f.Username,
		Nickname:             f.Nickname,
		Avatar:               f.Avatar,
		Ctime:                f.Ctime,
		CntComment:           f.CntComment,
		CntView:              f.CntView,
		CntCollect:           f.CntCollect,
		HeadImage:            f.HeadImage,
		SpaceGuid:            f.SpaceGuid,
		IsAllowCreateComment: f.IsAllowCreateComment,
		IsSiteTop:            f.IsSiteTop,
		IsRecommend:          f.IsRecommend,
		Format:               commonv1.FILE_FORMAT(f.FileFormat),
		EmojiList:            f.EmojiList,
		Content:              f.Content,
	}
}

func (f *FileCache) ToArticlePb(layout commonv1.SPC_LAYOUT) *articlev1.ArticleShow {
	return &articlev1.ArticleShow{
		Guid:                 f.Guid,
		Name:                 f.Name,
		Uid:                  f.CreatedUid,
		Nickname:             f.Nickname,
		Avatar:               f.Avatar,
		Ctime:                f.Ctime,
		CntComment:           f.CntComment,
		CntView:              f.CntView,
		CntCollect:           f.CntCollect,
		EmojiList:            f.EmojiList,
		HeadImage:            f.HeadImage,
		SpaceGuid:            f.SpaceGuid,
		IsAllowCreateComment: f.IsAllowCreateComment,
		IsSiteTop:            f.IsSiteTop,
		IsRecommend:          f.IsRecommend,
		Format:               commonv1.FILE_FORMAT(f.FileFormat),
	}
}

func (list FileCaches) ToQuestionPb() []*questionv1.QuestionShow {
	output := make([]*questionv1.QuestionShow, 0)
	for _, value := range list {
		output = append(output, value.ToQuestionPb())
	}
	return output
}

func (f *FileCache) ToQuestionPb() *questionv1.QuestionShow {
	return &questionv1.QuestionShow{
		Guid:                 f.Guid,
		Name:                 f.Name,
		Uid:                  f.CreatedUid,
		Nickname:             f.Nickname,
		Avatar:               f.Avatar,
		Ctime:                f.Ctime,
		List:                 f.EmojiList,
		SpaceGuid:            f.SpaceGuid,
		Content:              f.Content,
		IsAllowCreateComment: f.IsAllowCreateComment,
		CntComment:           f.CntComment,
		CntView:              f.CntView,
		CntCollect:           f.CntCollect,
		IsCollect:            f.IsCollect,
		Format:               commonv1.FILE_FORMAT(f.FileFormat),
	}
}

func (list FileCaches) ToAnswerPb() []*questionv1.AnswerShow {
	output := make([]*questionv1.AnswerShow, 0)
	for _, value := range list {
		output = append(output, value.ToAnswerPb())
	}
	return output
}

// ToQAPb question answer都有
func (list FileCaches) ToQAPb() []*questionv1.QAShow {
	output := make([]*questionv1.QAShow, 0)
	for _, value := range list {
		output = append(output, value.ToQAPb())
	}
	return output
}

func (f *FileCache) ToAnswerPb() *questionv1.AnswerShow {
	return &questionv1.AnswerShow{
		Guid:                 f.Guid,
		Name:                 f.Name,
		Uid:                  f.CreatedUid,
		Nickname:             f.Nickname,
		Avatar:               f.Avatar,
		Ctime:                f.Ctime,
		List:                 f.EmojiList,
		SpaceGuid:            f.SpaceGuid,
		IsAllowCreateComment: f.IsAllowCreateComment,
		CntComment:           f.CntComment,
		CntView:              f.CntView,
		CntCollect:           f.CntCollect,
		IsCollect:            f.IsCollect,
		Format:               commonv1.FILE_FORMAT(f.FileFormat),
	}
}

func (f *FileCache) ToQAPb() *questionv1.QAShow {
	return &questionv1.QAShow{
		Guid:                 f.Guid,
		ParentGuid:           f.ParentGuid,
		Name:                 f.Name,
		Uid:                  f.CreatedUid,
		Nickname:             f.Nickname,
		Avatar:               f.Avatar,
		Ctime:                f.Ctime,
		List:                 f.EmojiList,
		SpaceGuid:            f.SpaceGuid,
		IsAllowCreateComment: f.IsAllowCreateComment,
		CntComment:           f.CntComment,
		CntView:              f.CntView,
		CntCollect:           f.CntCollect,
		IsCollect:            f.IsCollect,
		Format:               commonv1.FILE_FORMAT(f.FileFormat),
		// Summary:              summary,
		// IsReadMore:           f.IsReadMore,
	}
}

func (f *FileCache) ToFilePb() *commonv1.FileInfo {
	return &commonv1.FileInfo{
		Guid:                 f.Guid,
		Name:                 f.Name,
		Ctime:                f.Ctime,
		Utime:                f.Utime,
		CreatedUid:           f.CreatedUid,
		UpdatedUid:           f.UpdatedUid,
		Type:                 commonv1.FILE_TYPE(f.FileType),
		Format:               commonv1.FILE_FORMAT(f.FileFormat),
		CntComment:           f.CntCollect,
		CntView:              int32(f.CntView),
		CntCollect:           int32(f.CntCollect),
		HeadImage:            f.HeadImage,
		IsAllowCreateComment: f.IsAllowCreateComment,
		ParentGuid:           f.ParentGuid,
		Size:                 f.Size,
		BizStatus:            commonv1.FILE_BIZSTS(f.BizStatus),
		ContentUrl:           f.ContentUrl,
	}
}
