package file

import (
	"context"
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	filev1 "ecodepost/pb/file/v1"
	statv1 "ecodepost/pb/stat/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"

	"gorm.io/gorm"
)

type GrpcServer struct{}

func (s GrpcServer) OpenComment(ctx context.Context, req *filev1.OpenCommentReq) (*filev1.OpenCommentRes, error) {
	// TODO implement me
	panic("implement me")
}

func (s GrpcServer) CloseComment(ctx context.Context, req *filev1.CloseCommentReq) (*filev1.CloseCommentRes, error) {
	// TODO implement me
	panic("implement me")
}

var _ filev1.FileServer = &GrpcServer{}

// EmojiList https://www.qqxiuzi.cn/zh/emoji.html
func (GrpcServer) EmojiList(ctx context.Context, req *filev1.EmojiListReq) (*filev1.EmojiListRes, error) {
	return &filev1.EmojiListRes{
		List: mysql.EmojiList(),
	}, nil
}

// MyEmojiList 根据文章guid和用户信息，返回他的emoji
func (GrpcServer) MyEmojiList(ctx context.Context, req *filev1.MyEmojiListReq) (*filev1.MyEmojiListRes, error) {
	guidMap := make(map[string]*filev1.EmojiMap)

	list, err := mysql.MyEmojiList(invoker.Db.WithContext(ctx), req.GetGuids(), req.GetUid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("my emoji list, err: " + err.Error())
	}

	for _, value := range list {
		var (
			emojiTmpMap map[int32]*commonv1.EmojiInfo
			ok          bool
			emojiMap    *filev1.EmojiMap
		)

		emojiMap, ok = guidMap[value.Guid]
		if !ok {
			emojiMap = &filev1.EmojiMap{}
			emojiTmpMap = make(map[int32]*commonv1.EmojiInfo, 0)
		} else {
			emojiTmpMap = emojiMap.Map
		}
		emojiTmpMap[value.V] = value.ToPb()
		emojiMap.Map = emojiTmpMap
		guidMap[value.Guid] = emojiMap
	}

	return &filev1.MyEmojiListRes{
		Map: guidMap,
	}, nil
}

// CreateEmoji Emoji Create
func (GrpcServer) CreateEmoji(ctx context.Context, req *filev1.CreateEmojiReq) (*filev1.CreateEmojiRes, error) {
	spaceGuid, err := mysql.FileSpaceGuidByGuid(invoker.Db.WithContext(ctx), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateEmoji FileSpaceGuidByGuid fail, err: " + err.Error())
	}
	if spaceGuid == "" {
		return nil, errcodev1.ErrDbError().WithMessage("space guid empty")
	}

	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = mysql.EmojiIncrease(tx, req.GetUid(), req.GetGuid(), req.GetV())
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("CreateEmoji fail, err: " + err.Error())
		}
		err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetGuid(), "cnt_like", 1)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("create comment track total fail2, err: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// _, err = invoker.GrpcTrack.Total(ctx, &trackv1.TotalReq{
	// 	Event:     commonv1.TRACK_TOTAL_EVENT_CREATE_EMOJI,
	// 	Tid:       etrace.ExtractTraceID(ctx),
	// 	Uid:       req.GetUid(),
	// 	FileGuid:  req.GetGuid(),
	// 	CmtGuid:   req.GetCmtGuid(),
	// 	SpaceGuid: spaceGuid,
	// 	Ts:        time.Now().UnixMilli(),
	// })
	// if err != nil {
	// 	return nil, errcodev1.ErrInternal().WithMessage("track fail, err: " + err.Error())
	// }

	return &filev1.CreateEmojiRes{}, nil
}

// DeleteEmoji Emoji Delete
func (GrpcServer) DeleteEmoji(ctx context.Context, req *filev1.DeleteEmojiReq) (*filev1.DeleteEmojiRes, error) {
	spaceGuid, err := mysql.FileSpaceGuidByGuid(invoker.Db.WithContext(ctx), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateEmoji FileSpaceGuidByGuid fail, err: " + err.Error())
	}
	if spaceGuid == "" {
		return nil, errcodev1.ErrDbError().WithMessage("space guid empty")
	}

	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = mysql.EmojiDecrease(tx, req.GetUid(), req.GetGuid(), req.GetV())
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("DeleteEmoji fail, err: " + err.Error())
		}
		err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetGuid(), "cnt_like", -1)
		if err != nil {
			return fmt.Errorf("create comment track total fail2, err: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// _, err = invoker.GrpcTrack.Total(ctx, &trackv1.TotalReq{
	// 	Event:     commonv1.TRACK_TOTAL_EVENT_DELETE_EMOJI,
	// 	Tid:       etrace.ExtractTraceID(ctx),
	// 	Uid:       req.GetUid(),
	// 	FileGuid:  req.GetGuid(),
	// 	CmtGuid:   req.GetCmtGuid(),
	// 	SpaceGuid: spaceGuid,
	// 	Ts:        time.Now().UnixMilli(),
	// })
	// if err != nil {
	// 	return nil, errcodev1.ErrInternal().WithMessage("track fail, err: " + err.Error())
	// }
	return &filev1.DeleteEmojiRes{}, nil
}

// CollectionCreate 收藏某个目标到某几个收藏夹
func (GrpcServer) CollectionCreate(ctx context.Context, req *filev1.CollectionCreateReq) (*filev1.CollectionCreateRes, error) {
	info, err := invoker.GrpcStat.CollectionCreate(ctx, &statv1.CollectionCreateReq{
		Uid:                req.GetUid(),
		CollectionGroupIds: req.GetCollectionGroupIds(),
		BizGuid:            req.GetBizGuid(),
		BizType:            req.GetBizType(),
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("change space sort fail, err: " + err.Error())
	}
	err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetBizGuid(), "cnt_collect", int(info.GetDelta()))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("change space sort fail, err: " + err.Error())
	}
	return &filev1.CollectionCreateRes{}, nil
}

// CollectionDelete 从几个收藏夹取消收藏某个目标
func (GrpcServer) CollectionDelete(ctx context.Context, req *filev1.CollectionDeleteReq) (*filev1.CollectionDeleteRes, error) {
	info, err := invoker.GrpcStat.CollectionDelete(ctx, &statv1.CollectionDeleteReq{
		Uid:                req.GetUid(),
		CollectionGroupIds: req.GetCollectionGroupIds(),
		BizGuid:            req.GetBizGuid(),
		BizType:            req.GetBizType(),
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("change space sort fail, err: " + err.Error())
	}
	err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetBizGuid(), "cnt_collect", -1*int(info.GetDelta()))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("change space sort fail, err: " + err.Error())
	}
	return &filev1.CollectionDeleteRes{}, nil
}

func (GrpcServer) UpdateFileSize(ctx context.Context, req *filev1.UpdateFileSizeReq) (*filev1.UpdateFileSizeRes, error) {
	err := mysql.FileUpdate(invoker.Db.WithContext(ctx), req.Guid, map[string]any{"size": req.Size, "status": int32(commonv1.FILE_STATUS_UPLOADED_SUCC)})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileUpdate fail, err:" + err.Error())
	}
	return &filev1.UpdateFileSizeRes{}, nil
}

func deleteDocument(ctx context.Context, guid string, uid int64) error {
	// 查询fileInfo
	fileInfo, err := mysql.FileInfoMustExistsEerror(invoker.Db.WithContext(ctx), "id,space_guid,size", guid)
	if err != nil {
		return err
	}

	// 删除file
	err = service.File.Delete(ctx, invoker.Db.WithContext(ctx), uid, fileInfo.SpaceGuid, guid)
	if err != nil {
		return errcodev1.ErrDbError().WithMessage("DeleteDocument fail, err: " + err.Error())
	}

	if err != nil {
		return errcodev1.ErrDbError().WithMessage("Delete fail, err: " + err.Error())
	}
	return nil
}
