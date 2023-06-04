package column

import (
	"context"
	"fmt"

	userv1 "ecodepost/pb/user/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"

	columnv1 "ecodepost/pb/column/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
)

type GrpcServer struct{}

var _ columnv1.ColumnServer = (*GrpcServer)(nil)

func (GrpcServer) CreateSpaceInfo(ctx context.Context, req *columnv1.CreateSpaceInfoReq) (*columnv1.CreateSpaceInfoRes, error) {
	err := invoker.Db.WithContext(ctx).Create(&mysql.SpaceColumn{
		SpaceGuid: req.GetSpaceGuid(),
		AuthorUid: req.GetAuthorUid(),
	}).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("update space info fail, err: " + err.Error())
	}
	return &columnv1.CreateSpaceInfoRes{}, nil
}

func (GrpcServer) UpdateSpaceInfo(ctx context.Context, req *columnv1.UpdateSpaceInfoReq) (*columnv1.UpdateSpaceInfoRes, error) {
	err := invoker.Db.WithContext(ctx).Model(mysql.SpaceColumn{}).Where("space_guid =?", req.GetSpaceGuid()).Updates(map[string]any{
		"author_uid": req.GetAuthorUid(),
	}).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("update space info fail, err: " + err.Error())
	}
	return &columnv1.UpdateSpaceInfoRes{}, nil
}

func (GrpcServer) GetSpaceInfo(ctx context.Context, req *columnv1.GetSpaceInfoReq) (*columnv1.GetSpaceInfoRes, error) {
	var columnInfo mysql.SpaceColumn

	err := invoker.Db.WithContext(ctx).Model(mysql.SpaceColumn{}).Where("space_guid =?", req.GetSpaceGuid()).Find(&columnInfo).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("get space column info fail, err: " + err.Error())
	}

	if columnInfo.AuthorUid == 0 {
		return &columnv1.GetSpaceInfoRes{
			AuthorUid:      0,
			AuthorNickname: "",
			AuthorAvatar:   "",
		}, nil
	}

	userInfo, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{
		Uid: columnInfo.AuthorUid,
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("get space column info userinfo fail, err: " + err.Error())
	}

	return &columnv1.GetSpaceInfoRes{
		AuthorUid:      userInfo.GetUser().Uid,
		AuthorNickname: userInfo.GetUser().Nickname,
		AuthorAvatar:   userInfo.GetUser().Avatar,
	}, nil
}

func (GrpcServer) ChangeSort(ctx context.Context, req *columnv1.ChangeSortReq) (*columnv1.ChangeSortRes, error) {
	err := service.File.ChangeSort(ctx, req.GetFileGuid(), req.TargetFileGuid, req.DropPosition, req.GetParentFileGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("change space sort fail, err: " + err.Error())
	}
	return &columnv1.ChangeSortRes{}, nil
}

func (GrpcServer) ListFile(ctx context.Context, req *columnv1.ListFileReq) (*columnv1.ListFileRes, error) {
	fis, err := service.File.ListFileInfos(ctx, req.SpaceGuid, req.Uid)
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("ListFileInfos fail, err:" + err.Error())
	}
	fmt.Printf("fis--------------->"+"%+v\n", fis)
	return &columnv1.ListFileRes{Files: fis}, nil
}
