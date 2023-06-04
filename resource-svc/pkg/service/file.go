package service

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service/cache"
	"ecodepost/resource-svc/pkg/utils/x"

	commonv1 "ecodepost/pb/common/v1"
	// trackv1 "ecodepost/pb/track/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/ego-component/egorm"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type file struct{}

func InitFile() *file {
	return &file{}
}

func (*file) create(ctx context.Context, tx *gorm.DB, data *mysql.File) (err error) {
	err = mysql.FileCreate(tx, data)
	if err != nil {
		return fmt.Errorf("file create fail, err: %w", err)
	}
	// _, err = invoker.GrpcTrack.Total(ctx, &trackv1.TotalReq{
	// 	Event:     commonv1.TRACK_TOTAL_EVENT_CREATE_FILE,
	// 	Tid:       etrace.ExtractTraceID(ctx),
	// 	Uid:       data.CreatedBy,
	// 	FileGuid:  data.Guid,
	// 	CmtGuid:   data.CmtGuid,
	// 	SpaceGuid: data.SpaceGuid,
	// 	Ts:        time.Now().UnixMilli(),
	// })
	// if err != nil {
	// 	return fmt.Errorf("file create fail2, err: %w", err)
	// }
	return nil
}

func (*file) Delete(ctx context.Context, db *gorm.DB, uid int64, spaceGuid, fileGuid string) (err error) {
	err = mysql.FileRemove(db, uid, spaceGuid, fileGuid)
	if err != nil {
		return fmt.Errorf("file delete fail, err: %w", err)
	}
	// _, err = invoker.GrpcTrack.Total(ctx, &trackv1.TotalReq{
	// 	Event:     commonv1.TRACK_TOTAL_EVENT_DELETE_FILE,
	// 	Tid:       etrace.ExtractTraceID(ctx),
	// 	Uid:       uid,
	// 	FileGuid:  fileGuid,
	// 	CmtGuid:   cmtGuid,
	// 	SpaceGuid: spaceGuid,
	// 	Ts:        time.Now().UnixMilli(),
	// })
	// if err != nil {
	// 	return fmt.Errorf("delete create fail2, err: %w", err)
	// }
	return nil
}

// GetContentByCreator 创作者获取内容
// 适用于文档、问答
// func (*file) GetContentByCreator(ctx context.Context, cmtGuid, guid string) (str []byte, err error) {
//	fileInfo, err := mysql.FileContentInfoByGuid(invoker.Db.WithContext(ctx), cmtGuid, guid)
//	if err != nil {
//		err = fmt.Errorf("GetContentByCreator fail, err: %w", err)
//		return
//	}
//	if fileInfo.Guid == "" {
//		err = fmt.Errorf("GetContentByCreator fail2, file not exist")
//		return
//	}
//
//	str, err = invoker.AliOss.GetObject(fileInfo.ContentKey)
//	if err != nil {
//		err = fmt.Errorf("GetContentByCreator fail4, err: %w", err)
//		return
//	}
//	return
// }

// GetContentUrl 普通人获取内容
// 适用于文档、问答
//func (*file) GetContentUrl(ctx context.Context, guid string) (str string, err error) {
//	fileInfo, err := mysql.FileContentInfoByGuid(invoker.Db.WithContext(ctx), guid)
//	if err != nil {
//		err = fmt.Errorf("GetContentUrl fail, err: %w", err)
//		return
//	}
//	if fileInfo.Guid == "" {
//		err = fmt.Errorf("GetContentUrl fail2, file not exist")
//		return
//	}
//
//	str, err = invoker.AliOss.CdnAuthURL(fileInfo.ContentKey, fileInfo.Hash)
//	if err != nil {
//		err = fmt.Errorf("GetContentUrl fail3, err: %w", err)
//		return
//	}
//	return
//}

// HomeRecommendList 根据分页条件查询list
// todo 公开的space guid in 一把
// func (f *file) HomeRecommendList(ctx context.Context, db *gorm.DB, cmtGuid string) (respList mysql.FileCaches, err error) {
//	conds := egorm.Conds{}
//	conds["space_top_time"] = 0
//	conds["cmt_guid"] = cmtGuid
//	conds["file_type"] = int(commonv1.FILE_TYPE_DOCUMENT)
//	reqList := &commonv1.Pagination{
//		PageSize: 10,
//	}
//	listSort := commonv1.CMN_SORT_RECOMMEND_SCORE
//	return f.ListPageCache(ctx, db, conds, reqList, listSort)
// }

// FileListPage 根据分页条件查询list
func (f *file) FileListPage(ctx context.Context, db *gorm.DB, spaceGuid string, reqList *commonv1.Pagination, listSort commonv1.CMN_FILE_SORT) (respList mysql.FileCaches, err error) {
	conds := egorm.Conds{}
	conds["space_guid"] = spaceGuid
	conds["space_top_time"] = 0
	// todo 后面要根据space类型渲染不同的file list
	conds["file_type"] = egorm.Cond{
		Op:  "!=",
		Val: commonv1.FILE_TYPE_FILE.Number(),
	}
	conds["parent_guid"] = ""
	return f.ListPageCache(ctx, db, conds, reqList, listSort)
}

// FileListPageByParent 根据分页条件查询list
func (f *file) FileListPageByParent(ctx context.Context, db *gorm.DB, parentGuid string, reqList *commonv1.Pagination, listSort commonv1.CMN_FILE_SORT) (respList mysql.FileCaches, err error) {
	conds := egorm.Conds{}
	conds["parent_guid"] = parentGuid
	return f.ListPageCache(ctx, db, conds, reqList, listSort)
}

// QuestionListPage 根据分页条件查询list
func (f *file) QuestionListPage(ctx context.Context, db *gorm.DB, spaceGuid string, reqList *commonv1.Pagination, listSort commonv1.CMN_FILE_SORT) (respList mysql.FileCaches, err error) {
	conds := egorm.Conds{}
	conds["space_guid"] = spaceGuid
	conds["parent_guid"] = ""
	conds["file_type"] = int(commonv1.FILE_TYPE_QUESTION)
	return f.ListPageCache(ctx, db, conds, reqList, listSort)
}

// AnswerListPage 根据分页条件查询list
func (f *file) AnswerListPage(ctx context.Context, db *gorm.DB, parentGuid string, reqList *commonv1.Pagination, listSort commonv1.CMN_FILE_SORT) (respList mysql.FileCaches, err error) {
	conds := egorm.Conds{}
	conds["parent_guid"] = parentGuid
	conds["file_type"] = int(commonv1.FILE_TYPE_QUESTION)
	return f.ListPageCache(ctx, db, conds, reqList, listSort)
}

// PublicUserArticleListPage 根据创建者查询文章列表
func (f *file) PublicUserArticleListPage(ctx context.Context, db *gorm.DB, createdBy int64, reqList *commonv1.Pagination) (respList mysql.FileCaches, err error) {
	conds := egorm.Conds{}
	conds["created_by"] = createdBy
	conds["file_type"] = int(commonv1.FILE_TYPE_DOCUMENT)
	return f.ListPageCache(ctx, db, conds, reqList, commonv1.CMN_SORT_CREATE_TIME)
}

// PublicUserQAListPage 根据创建者查询文章列表
func (f *file) PublicUserQAListPage(ctx context.Context, db *gorm.DB, createdBy int64, reqList *commonv1.Pagination) (respList mysql.FileCaches, err error) {
	conds := egorm.Conds{}
	conds["created_by"] = createdBy
	conds["file_type"] = int(commonv1.FILE_TYPE_QUESTION)
	return f.ListPageCache(ctx, db, conds, reqList, commonv1.CMN_SORT_CREATE_TIME)
}

func (f *file) PublicDriveListPage(ctx context.Context, db *gorm.DB, cmtGuid, spaceGuid, parentGuid string, types []commonv1.FILE_TYPE, formats []commonv1.FILE_FORMAT,
	reqList *commonv1.Pagination, listSort commonv1.CMN_FILE_SORT) (respList mysql.FileCaches, err error) {

	conds := egorm.Conds{}
	conds["space_guid"] = spaceGuid
	conds["space_top_time"] = 0
	conds["cmt_guid"] = cmtGuid
	conds["parent_guid"] = parentGuid
	if len(types) != 0 {
		conds["file_type"] = x.Es2I32s(types)
	}
	if len(formats) != 0 {
		conds["file_format"] = x.Es2I32s(types)
	}
	return f.ListPageCache(ctx, db, conds, reqList, listSort)
}

func (f *file) PublicDriveInfo(ctx context.Context, db *gorm.DB, cmtGuid, spaceGuid, guid string, tp commonv1.FILE_TYPE) (file *mysql.FileCache, err error) {
	return f.Info(ctx, db, guid)
}

func (f *file) emojiList(db *gorm.DB, guids []string) (emojiListMap map[string][]*commonv1.EmojiInfo, err error) {
	mysqlEmojiList, err := mysql.FileEmojiList(db, guids)
	if err != nil {
		err = fmt.Errorf("emoji list fail, err: %w", err)
		return
	}
	emojiListMap = make(map[string][]*commonv1.EmojiInfo)
	// 初始化所有的emojiList
	// for _, guid := range guids {
	//	emojiListMap[guid] = mysql.EmojiList()
	// }

	for _, emojiInfo := range mysqlEmojiList {
		emojiTmpList := make([]*commonv1.EmojiInfo, 0)
		// todo 反射
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    1,
			Emoji: mysql.GetOneEmoji(1).GetEmoji(),
			Cnt:   emojiInfo.V1,
		})
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    2,
			Emoji: mysql.GetOneEmoji(2).GetEmoji(),
			Cnt:   emojiInfo.V2,
		})
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    3,
			Emoji: mysql.GetOneEmoji(3).GetEmoji(),
			Cnt:   emojiInfo.V3,
		})
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    4,
			Emoji: mysql.GetOneEmoji(4).GetEmoji(),
			Cnt:   emojiInfo.V4,
		})
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    5,
			Emoji: mysql.GetOneEmoji(5).GetEmoji(),
			Cnt:   emojiInfo.V5,
		})
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    6,
			Emoji: mysql.GetOneEmoji(6).GetEmoji(),
			Cnt:   emojiInfo.V6,
		})
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    7,
			Emoji: mysql.GetOneEmoji(7).GetEmoji(),
			Cnt:   emojiInfo.V7,
		})
		emojiTmpList = append(emojiTmpList, &commonv1.EmojiInfo{
			Id:    8,
			Emoji: mysql.GetOneEmoji(8).GetEmoji(),
			Cnt:   emojiInfo.V8,
		})
		emojiListMap[emojiInfo.Guid] = emojiTmpList
	}
	return
}

func (f *file) emojiInfo(db *gorm.DB, guid string) (emojiList []*commonv1.EmojiInfo, err error) {
	emojiInfo, err := mysql.FileEmojiInfo(db, guid)
	if err != nil {
		err = fmt.Errorf("emoji list fail, err: %w", err)
		return
	}

	emojiList = make([]*commonv1.EmojiInfo, 0)
	// todo 反射
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    1,
		Emoji: mysql.GetOneEmoji(1).GetEmoji(),
		Cnt:   emojiInfo.V1,
	})
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    2,
		Emoji: mysql.GetOneEmoji(2).GetEmoji(),
		Cnt:   emojiInfo.V2,
	})
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    3,
		Emoji: mysql.GetOneEmoji(3).GetEmoji(),
		Cnt:   emojiInfo.V3,
	})
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    4,
		Emoji: mysql.GetOneEmoji(4).GetEmoji(),
		Cnt:   emojiInfo.V4,
	})
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    5,
		Emoji: mysql.GetOneEmoji(5).GetEmoji(),
		Cnt:   emojiInfo.V5,
	})
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    6,
		Emoji: mysql.GetOneEmoji(6).GetEmoji(),
		Cnt:   emojiInfo.V6,
	})
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    7,
		Emoji: mysql.GetOneEmoji(7).GetEmoji(),
		Cnt:   emojiInfo.V7,
	})
	emojiList = append(emojiList, &commonv1.EmojiInfo{
		Id:    8,
		Emoji: mysql.GetOneEmoji(8).GetEmoji(),
		Cnt:   emojiInfo.V8,
	})
	return
}

// ListFileGuidsFromDB 从数据中查询FileGuids
func (f *file) ListFileGuidsFromDB(ctx context.Context, db *gorm.DB, conds egorm.Conds, reqList *commonv1.Pagination, listSort commonv1.CMN_FILE_SORT) (mysql.FileGuids, error) {
	if reqList.PageSize == 0 || reqList.PageSize > 200 {
		reqList.PageSize = 20
	}
	if reqList.CurrentPage == 0 {
		reqList.CurrentPage = 1
	}

	fileGuids := make(mysql.FileGuids, 0)
	if reqList.Sort == "" {
		reqList.Sort = cast.ToString(commonv1.CMN_SORT_CREATE_TIME.Number())
	}
	sort := "id desc"
	switch listSort {
	case commonv1.CMN_SORT_CREATE_TIME:
		sort = "id desc"
		reqList.Sort = cast.ToString(commonv1.CMN_SORT_CREATE_TIME.Number())
	case commonv1.CMN_SORT_HOT_SCORE:
		sort = "hot_score desc"
		reqList.Sort = cast.ToString(commonv1.CMN_SORT_HOT_SCORE.Number())
	case commonv1.CMN_SORT_RECOMMEND_SCORE:
		sort = "recommend_score desc"
		reqList.Sort = cast.ToString(commonv1.CMN_SORT_RECOMMEND_SCORE.Number())
	}

	conds["dtime"] = 0
	sql, binds := egorm.BuildQuery(conds)
	db = db.Table("file").Select("guid").Where(sql, binds...)
	err := db.Count(&reqList.Total).Error
	if err != nil {
		return nil, err
	}
	err = db.Order(sort).Offset(int((reqList.CurrentPage - 1) * reqList.PageSize)).Limit(int(reqList.PageSize)).Find(&fileGuids).Error
	if err != nil {
		return nil, err
	}

	return fileGuids, nil
}

// ListPageCache 根据分页条件查询FileCaches list
func (f *file) ListPageCache(ctx context.Context, db *gorm.DB, conds egorm.Conds, reqList *commonv1.Pagination, listSort commonv1.CMN_FILE_SORT) (fileCaches mysql.FileCaches, err error) {
	// 查询FileGuids
	fileGuids, err := f.ListFileGuidsFromDB(ctx, db, conds, reqList, listSort)
	if err != nil {
		return nil, err
	}
	// 遍历fileGuids，获取guid
	guids := lo.Map(fileGuids, func(v mysql.FileGuid, _ int) string { return v.Guid })
	return f.ListCacheByGuids(ctx, db, guids)

	// emojiListMap, err := f.emojiList(db, guids)
	// if err != nil {
	//	return
	// }
	//
	// // 批量查询缓存
	// cacheMap, err := cache.FileCache.BatchGetInfo(ctx, guids)
	// if err != nil {
	//	return
	// }
	// uids := lo.MapToSlice(cacheMap, func(_ string, v *dto.FileCache) int64 { return v.CreatedUid })
	//
	// // 根据uids批量查询用户
	// users, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{UidList: uids})
	// // 不return，用户服务不影响，展示，记录log
	// if err != nil {
	//	return
	// }
	//
	// fileCaches = make(dto.FileCaches, 0)
	// for _, value := range fileGuids {
	//	cacheInfo, flag := cacheMap[value.Guid]
	//	if !flag {
	//		fileCaches = append(fileCaches, &dto.FileCache{Guid: value.Guid})
	//		continue
	//	}
	//	cacheInfo.EmojiList = emojiListMap[value.Guid]
	//	cacheInfo.Nickname = users.GetUserMap()[cacheInfo.CreatedUid].GetNickname()
	//	cacheInfo.Avatar = users.GetUserMap()[cacheInfo.CreatedUid].GetAvatar()
	//	cacheInfo.CntView = value.CntView
	//	cacheInfo.CntComment = value.CntComment
	//	cacheInfo.CntCollect = value.CntCollect
	//	cacheInfo.Size = value.Size
	//	cacheInfo.FileType = value.FileType
	//	fileCaches = append(fileCaches, cacheInfo)
	// }
	// return
}

// ListCacheByGuids 根据 guids 拿到cache
// 一定要注意guids，是有序的
func (f *file) ListCacheByGuids(ctx context.Context, db *gorm.DB, sortedGuids []string) (fileCaches mysql.FileCaches, err error) {
	if len(sortedGuids) == 0 {
		return
	}

	// 遍历fileGuids，获取guid
	emojiListMap, err := f.emojiList(db, sortedGuids)
	if err != nil {
		return
	}

	var mysqlFiles mysql.Files
	err = db.Select("guid,cnt_view,cnt_comment,cnt_collect,size,ip").Where("guid in (?)", sortedGuids).Find(&mysqlFiles).Error
	if err != nil {
		return
	}
	// 批量查询缓存
	cacheMap, err := cache.FileCache.BatchGetInfo(ctx, sortedGuids)
	if err != nil {
		return
	}
	uids := lo.MapToSlice(cacheMap, func(_ string, v *mysql.FileCache) int64 { return v.CreatedUid })

	// 根据uids批量查询用户
	users, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{UidList: uids})
	// 不return，用户服务不影响，展示，记录log
	if err != nil {
		return
	}

	fileTotalMap := lo.SliceToMap(mysqlFiles, func(t *mysql.File) (string, *mysql.File) { return t.Guid, t })
	fileCaches = lo.Map(sortedGuids, func(guid string, i int) *mysql.FileCache {
		cacheInfo, flag := cacheMap[guid]
		if !flag {
			return &mysql.FileCache{Guid: guid}
		}
		// 数据库没有，就用默认emoji list
		emojiList, flag := emojiListMap[guid]
		if flag {
			cacheInfo.EmojiList = emojiList
		} else {
			// todo 后续是跟着space变化
			cacheInfo.EmojiList = mysql.EmojiList()
		}
		cacheInfo.Nickname = users.GetUserMap()[cacheInfo.CreatedUid].GetNickname()
		cacheInfo.Avatar = users.GetUserMap()[cacheInfo.CreatedUid].GetAvatar()
		cacheInfo.Username = users.GetUserMap()[cacheInfo.CreatedUid].GetName()
		totalInfo, flag := fileTotalMap[guid]
		if !flag {
			return cacheInfo
		}
		cacheInfo.CntView = totalInfo.CntView
		cacheInfo.CntComment = totalInfo.CntComment
		cacheInfo.CntCollect = totalInfo.CntCollect
		cacheInfo.Size = totalInfo.Size
		cacheInfo.Ip = totalInfo.Ip
		return cacheInfo
	})
	return
}

// ListAllBySpaceGuid 根据分页条件查询list
func (f *file) ListAllBySpaceGuid(db *gorm.DB, spaceGuid string) (respList mysql.Files, err error) {
	db = db.Table("file").Select("guid,parent_guid,name,content_key,head_image,file_type").Where("space_guid = ? and  dtime = ?", spaceGuid, 0)
	err = db.Find(&respList).Error
	if err != nil {
		return
	}
	return
}

func (f *file) Info(ctx context.Context, db *gorm.DB, guid string) (cacheInfo *mysql.FileCache, err error) {
	emojiList, err := f.emojiInfo(db, guid)
	if err != nil {
		return
	}

	cacheInfo, err = cache.FileCache.GetInfo(ctx, guid)
	if err != nil {
		return
	}
	fmt.Printf("cacheInfo--------------->"+"%+v\n", cacheInfo)

	userInfo, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{Uid: cacheInfo.CreatedUid})
	if err != nil {
		return
	}
	cacheInfo.EmojiList = emojiList
	cacheInfo.Nickname = userInfo.User.GetNickname()
	cacheInfo.Avatar = userInfo.User.GetAvatar()
	return
}

func (f *file) ListFileInfos(ctx context.Context, spaceGuid string, uid int64) ([]*commonv1.FileInfo, error) {
	var edges []mysql.Edge
	fis := make([]*commonv1.FileInfo, 0)
	db := invoker.Db.WithContext(ctx)
	if err := db.Model(&mysql.Edge{}).Where("space_guid = ?", spaceGuid).Order("`sort` asc").Find(&edges).Error; err != nil {
		return nil, fmt.Errorf("select sectionsCnt fail, %w", err)
	}
	edgesMap := make(map[string]mysql.Edge)
	guids := lo.Map(edges, func(edge mysql.Edge, i int) string {
		edgesMap[edge.FileGuid] = edge
		return edge.FileGuid
	})
	// 批量查询缓存
	fcMap, err := cache.FileCache.BatchGetInfo(ctx, guids)
	if err != nil {
		return nil, fmt.Errorf("BatchGetInfo fail, %w", err)
	}
	for _, guid := range guids {
		fc, ok := fcMap[guid]
		if !ok {
			continue
		}
		fi := fc.ToFilePb()
		if edge, ok := edgesMap[guid]; ok {
			fi.Sort = edge.Sort
		}
		fis = append(fis, fi)
	}
	fmt.Printf("fis--------------->"+"%+v\n", fis)
	return fis, nil
}
