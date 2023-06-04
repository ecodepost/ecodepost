package stat

import (
	"context"
	"time"

	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"

	errcodev1 "ecodepost/pb/errcode/v1"
	statv1 "ecodepost/pb/stat/v1"
	"github.com/ego-component/egorm"
)

type GrpcServer struct{}

var _ statv1.StatServer = &GrpcServer{}

// MyCollectionListByFileGuids 根据文件GUIDS查看收藏列表
func (GrpcServer) MyCollectionListByFileGuids(ctx context.Context, req *statv1.MyCollectionListByFileGuidsReq) (*statv1.MyCollectionListByFileGuidsRes, error) {
	res, err := mysql.CollectionListByFileGuids(invoker.Db.WithContext(ctx), req.GetUid(), req.GetFileGuids())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("fail ,err: " + err.Error())
	}
	return &statv1.MyCollectionListByFileGuidsRes{
		List: res.ToPb(),
	}, nil
}

// CollectionGroupCreate 创建收藏列表
func (GrpcServer) CollectionGroupCreate(ctx context.Context, req *statv1.CollectionGroupCreateReq) (*statv1.CollectionGroupCreateRes, error) {
	data := mysql.UserCollectionGroup{
		Uid:   req.GetUid(),
		Title: req.GetTitle(),
		Desc:  req.GetDesc(),
	}
	if err := mysql.UserCollectionGroupCreate(invoker.Db.WithContext(ctx), &data); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CollectionGroupCreate fail, err: " + err.Error())
	}
	return &statv1.CollectionGroupCreateRes{Id: data.ID}, nil
}

// CollectionGroupList 获取收藏列表
func (GrpcServer) CollectionGroupList(ctx context.Context, req *statv1.CollectionGroupListReq) (*statv1.CollectionGroupListRes, error) {
	res, err := mysql.UserCollectionGroupList(invoker.Db.WithContext(ctx), req.GetUid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CollectionGroupList fail, err: " + err.Error())
	}
	return &statv1.CollectionGroupListRes{List: res.ToPb()}, nil
}

// CollectionGroupUpdate 更新收藏列表信息
func (GrpcServer) CollectionGroupUpdate(ctx context.Context, req *statv1.CollectionGroupUpdateReq) (*statv1.CollectionGroupUpdateRes, error) {
	ups := make(map[string]any)
	if req.Title != nil {
		ups["title"] = req.Title
	}
	if req.Desc != nil {
		ups["desc"] = req.Desc
	}
	conds := egorm.Conds{"uid": req.Uid, "id": req.Id}
	if err := mysql.UserCollectionGroupUpdateX(invoker.Db.WithContext(ctx), conds, ups); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CollectionGroupList fail, err: " + err.Error())
	}
	return &statv1.CollectionGroupUpdateRes{}, nil
}

// CollectionGroupDelete 删除收藏列表
func (GrpcServer) CollectionGroupDelete(ctx context.Context, req *statv1.CollectionGroupDeleteReq) (*statv1.CollectionGroupDeleteRes, error) {
	tx := invoker.Db.WithContext(ctx).Begin()
	if err := mysql.CollectionGroupDelete(tx, req.Uid, req.Id); err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("CollectionGroupDelete fail, err: " + err.Error())
	}
	tx.Commit()

	return &statv1.CollectionGroupDeleteRes{}, nil
}

// CollectionCreate 创建收藏，收藏到某个收藏组
func (GrpcServer) CollectionCreate(ctx context.Context, req *statv1.CollectionCreateReq) (*statv1.CollectionCreateRes, error) {
	list := make(mysql.UserCollections, 0)
	var cnt int64
	nowTime := time.Now().Unix()
	for _, groupId := range req.GetCollectionGroupIds() {
		collectionInfo, err := mysql.GetCollectionInfo(invoker.Db.WithContext(ctx), req.GetUid(), req.GetBizGuid(), req.GetBizType(), groupId)
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("CollectionCreate fail2, err: " + err.Error())
		}
		if collectionInfo.ID > 0 {
			continue
		}
		cnt = cnt + 1
		list = append(list, &mysql.UserCollection{
			Ctime:   nowTime,
			Uid:     req.GetUid(),
			BizGuid: req.GetBizGuid(),
			BizType: req.GetBizType(),
			GroupId: groupId,
		})
	}

	tx := invoker.Db.WithContext(ctx).Begin()
	if err := mysql.CollectionCreateToGroup(tx, list); err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("CollectionCreate fail, err: " + err.Error())
	}
	tx.Commit()

	ids := make([]int64, 0, len(list))
	for _, v := range list {
		ids = append(ids, v.ID)
	}

	return &statv1.CollectionCreateRes{
		Ids:   ids,
		Delta: cnt,
	}, nil
}

// CollectionDelete 从某个收藏组删除某个收藏
func (GrpcServer) CollectionDelete(ctx context.Context, req *statv1.CollectionDeleteReq) (*statv1.CollectionDeleteRes, error) {
	cnt, err := mysql.CollectionDeleteCount(invoker.Db.WithContext(ctx), req.GetUid(), req.GetBizGuid(), req.GetBizType(), req.GetCollectionGroupIds())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CollectionDeleteFromGroup fail1, err: " + err.Error())
	}

	tx := invoker.Db.WithContext(ctx).Begin()
	err = mysql.CollectionDeleteFromGroup(tx, req.GetUid(), req.GetBizGuid(), req.GetBizType(), req.GetCollectionGroupIds())
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("CollectionDeleteFromGroup fail2, err: " + err.Error())
	}
	tx.Commit()
	return &statv1.CollectionDeleteRes{
		Delta: cnt,
	}, nil
}

// CollectionList 某个收藏组的收藏列表
func (GrpcServer) CollectionList(ctx context.Context, req *statv1.CollectionListReq) (*statv1.CollectionListRes, error) {
	ret, err := mysql.CollectionListByGroup(invoker.Db.WithContext(ctx), req.Uid, req.CollectionGroupId, req.Pagination)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CollectionListByGroup fail, err: " + err.Error())
	}
	return &statv1.CollectionListRes{List: ret.ToPb(), Pagination: nil}, nil
}

// IsCollection 是否收藏
func (GrpcServer) IsCollection(ctx context.Context, req *statv1.IsCollectionReq) (*statv1.IsCollectionRes, error) {
	cnt, err := mysql.CollectionCount(invoker.Db.WithContext(ctx), req.GetUid(), req.GetBizGuid(), req.GetBizType())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CollectionDeleteFromGroup fail1, err: " + err.Error())
	}

	isCollect := 0
	if cnt > 0 {
		isCollect = 1
	}

	return &statv1.IsCollectionRes{
		IsCollect: int32(isCollect),
	}, nil

}
