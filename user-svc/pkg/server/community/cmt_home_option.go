package community

import (
	"context"

	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"

	communityv1 "ecodepost/pb/community/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
)

// GetHomeOption 社区首页可选项信息
func (GrpcServer) GetHomeOption(ctx context.Context, req *communityv1.GetHomeOptionReq) (*communityv1.GetHomeOptionRes, error) {
	homeOption, err := mysql.GetSystemHomeOption(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("get home option fail, err: " + err.Error())
	}
	return homeOption.ToPb(), nil
}

// PutHomeOption 更新首页可选项信息
func (GrpcServer) PutHomeOption(ctx context.Context, req *communityv1.PutHomeOptionReq) (*communityv1.PutHomeOptionRes, error) {
	db := invoker.Db.WithContext(ctx).Begin()
	err := mysql.PutSystemHomeOption(db, req)
	if err != nil {
		db.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("put home option fail, err: " + err.Error())
	}
	db.Commit()
	return &communityv1.PutHomeOptionRes{}, nil
}
