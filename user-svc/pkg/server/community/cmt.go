package community

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	commonv1 "ecodepost/pb/common/v1"
	communityv1 "ecodepost/pb/community/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
)

type GrpcServer struct{}

var Svc communityv1.CommunityServer = &GrpcServer{}

// Home 社区首页信息
func (GrpcServer) Home(ctx context.Context, req *communityv1.HomeReq) (res *communityv1.HomeRes, err error) {
	cmtInfo, err := mysql.GetSystemBase(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("community info fail, err: " + err.Error())
	}

	res = &communityv1.HomeRes{
		Name:        cmtInfo.Name,
		Description: cmtInfo.Description,
		Logo:        cmtInfo.Logo,
		Access:      cmtInfo.Access,
	}
	//res.GongxinbuBeian = "沪ICP备2022029758号-1  "
	res.GongxinbuBeian = cmtInfo.GongxinbuBeian
	res.GonganbuBeian = cmtInfo.GonganBeian
	return res, nil
}

// Update 创建
func (GrpcServer) Update(ctx context.Context, req *communityv1.UpdateReq) (*communityv1.UpdateRes, error) {
	if req.Name != nil {
		name := strings.TrimSpace(req.GetName())
		if err := checkCommunityName(name); err != nil {
			return nil, errcodev1.ErrInvalidArgument().WithMessage("Apply invalid argument fail").WithMetadata(map[string]string{
				"cmtName": name,
			})
		}

	}
	db := invoker.Db.WithContext(ctx).Begin()
	err := mysql.PutSystemBase(db, req)
	if err != nil {
		db.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("Apply fail2, err: " + err.Error())
	}
	db.Commit()
	return &communityv1.UpdateRes{}, nil
}

// Info 查询指定社区
func (GrpcServer) Info(ctx context.Context, req *communityv1.InfoReq) (res *communityv1.InfoRes, err error) {
	cmtInfo, err := mysql.GetSystemBase(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("community info fail, err: " + err.Error())
	}
	return &communityv1.InfoRes{
		Community: &commonv1.CommunityInfo{
			Name:        cmtInfo.Name,
			Description: cmtInfo.Description,
			Logo:        cmtInfo.Logo,
			Access:      cmtInfo.Access,
		},
	}, nil
}

func checkCommunityName(communityName string) error {
	if communityName == "" {
		return fmt.Errorf("community name is empty")
	}
	if utf8.RuneCountInString(communityName) > 20 {
		return fmt.Errorf("community name is over size")
	}
	return nil
}
