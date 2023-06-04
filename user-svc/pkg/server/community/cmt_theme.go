package community

import (
	"context"
	"encoding/json"

	communityv1 "ecodepost/pb/community/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
)

// GetTheme 获取社区主题
func (GrpcServer) GetTheme(ctx context.Context, req *communityv1.GetThemeReq) (*communityv1.GetThemeRes, error) {
	themeInfo, err := mysql.GetSystemTheme(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("get theme fail1, err: " + err.Error())
	}

	customColorByte, _ := json.Marshal(themeInfo.CustomColor)
	return &communityv1.GetThemeRes{
		IsCustom:          themeInfo.IsCustom,
		ThemeName:         themeInfo.ThemeName,
		CustomColor:       string(customColorByte),
		DefaultAppearance: themeInfo.DefaultAppearance,
	}, nil
}

// SetTheme 设置社区主题
func (GrpcServer) SetTheme(ctx context.Context, req *communityv1.SetThemeReq) (*communityv1.SetThemeRes, error) {
	db := invoker.Db.WithContext(ctx).Begin()
	err := mysql.PutSystemTheme(db, req)
	if err != nil {
		db.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("SetTheme fail, err: " + err.Error())
	}
	db.Commit()
	return &communityv1.SetThemeRes{}, nil
}
