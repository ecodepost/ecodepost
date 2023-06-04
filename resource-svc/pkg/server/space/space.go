package space

import (
	"context"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	"github.com/samber/lo"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	spacev1 "ecodepost/pb/space/v1"
)

func (GrpcServer) CreateSpace(ctx context.Context, req *spacev1.CreateSpaceReq) (*spacev1.CreateSpaceRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_CREATE, req.GetOperateUid(), ""); err != nil {
		return nil, err
	}

	var data *commonv1.SpaceInfo
	data, err := service.Space.Create(ctx, req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("create space fail, err: " + err.Error())
	}
	return &spacev1.CreateSpaceRes{Info: data}, nil
}

func (GrpcServer) UpdateSpace(ctx context.Context, req *spacev1.UpdateSpaceReq) (resp *spacev1.UpdateSpaceRes, err error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_UPDATE, req.GetOperateUid(), ""); err != nil {
		return nil, err
	}

	// 获取原有的space信息
	oldSpaceInfo, err := mysql.GetSpaceInfoByGuid(invoker.Db.WithContext(ctx), "*", req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("GetSpaceInfoByGuid fail:" + err.Error())
	}

	ups := make(map[string]any, 0)
	if req.Name != nil {
		ups["name"] = req.GetName()
	}
	if req.Icon != nil {
		ups["icon"] = req.GetIcon()
		ups["icon_type"] = req.GetIconType().Number()
	}
	if req.SpaceGroupGuid != nil {
		ups["space_group_guid"] = req.GetSpaceGroupGuid()
	}

	if req.OriginPrice != nil {
		ups["origin_price"] = req.GetOriginPrice()
	}
	if req.Price != nil {
		ups["price"] = req.GetPrice()
	}
	if req.Desc != nil {
		ups["desc"] = req.GetDesc()
	}
	if req.HeadImage != nil {
		ups["head_image"] = req.GetHeadImage()
	}
	if req.Cover != nil {
		ups["cover"] = req.GetCover()
	}
	if req.SpaceType != nil {
		ups["space_type"] = req.GetSpaceType().Number()
	}
	if req.SpaceLayout != nil {
		ups["layout"] = req.GetSpaceLayout().Number()
	}

	// visibility 优先级 大于 charge type 大于 access
	if req.Access != nil {
		// 如果是需要付费的类型，则设置access为
		ups["access"] = req.GetAccess().Number()
	}
	// visibility 优先级 大于 charge type 大于 access
	if req.Link != nil {
		// 如果是需要付费的类型，则设置access为
		ups["link"] = req.GetLink()
	}

	// 因为付费的变更，和access不在一个地方
	// 所以当变为付费后，access要变成付费可进入类型
	// 这里需要强制变更，该access要用于展示，同时限制权限
	if req.ChargeType != nil {
		ups["charge_type"] = int32(req.GetChargeType())
		if lo.Contains([]commonv1.SPC_CT{commonv1.SPC_CT_BUYOUT, commonv1.SPC_CT_MEMBERSHIP}, req.GetChargeType()) {
			ups["access"] = commonv1.SPC_ACS_USER_PAY
		}
		if req.GetChargeType() == commonv1.SPC_CT_FREE {
			ups["access"] = commonv1.SPC_ACS_OPEN
		}
	}

	// todo，用户开启私密时候，需要检测有没有开启付费功能，如果有开启要提示，开启私密将无法付费
	// 私密状态不能付费，需要提示用户
	if req.Visibility != nil {
		// 如果原space info不是draft，但是用户强行改成draft是不允许的，说明是脚本运行的
		if oldSpaceInfo.Visibility != commonv1.CMN_VISBL_DRAFT && req.GetVisibility() == commonv1.CMN_VISBL_DRAFT {
			return nil, errcodev1.ErrInvalidArgument().WithMessage("visibility param fail")
		}
		ups["visibility"] = req.GetVisibility().Number()
		if req.GetVisibility() == commonv1.CMN_VISBL_SECRET {
			ups["access"] = commonv1.SPC_ACS_DENY_ALL
		}
	}
	// todo 更改access，要判断visibility，charge type类型，不允许随便更改
	// if req.Status != nil {
	// 	ups["status"] = req.GetStatus().Number()
	// }

	tx := invoker.Db.WithContext(ctx).Begin()
	if err := service.Space.Update(tx, req.GetOperateUid(), req.GetSpaceGuid(), ups); err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("UpdateSpaceAttr fail, err: " + err.Error())
	}

	// 更新option
	// option分到了不同的tab页面，所以不是一次性拿到所有的权限数据信息
	// 如果没有就插入，有的话就更新数据
	// todo 会导致有些老得权限数据下限，数据脏在数据库里，这种情况很少
	nowTime := time.Now().Unix()
	// options := make([]mysql.SpaceOption, 0)
	for _, value := range req.GetSpaceOptions() {
		optionInfo := mysql.SpaceOption{}
		err = tx.Where("guid = ? and option_name = ?", req.GetSpaceGuid(), value.GetSpaceOptionId().String()).Find(&optionInfo).Error
		if err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("UpdateSpaceAttr fail2, err: " + err.Error())
		}
		// 创建
		if optionInfo.Id == 0 {
			err = tx.Create(&mysql.SpaceOption{
				Guid:        req.GetSpaceGuid(),
				OptionName:  value.GetSpaceOptionId().String(),
				OptionValue: value.GetValue(),
				CreatedBy:   req.GetOperateUid(),
				Ctime:       nowTime,
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, errcodev1.ErrDbError().WithMessage("UpdateSpaceAttr fail3, err: " + err.Error())
			}
			continue
		}
		// 更新
		err = tx.Model(mysql.SpaceOption{}).Where("guid = ? and option_name = ?", req.GetSpaceGuid(), value.GetSpaceOptionId().String()).Update("option_value", value.GetValue()).Error
		if err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("UpdateSpaceAttr fail4, err: " + err.Error())
		}
	}
	tx.Commit()
	return &spacev1.UpdateSpaceRes{}, nil
}

func (GrpcServer) DeleteSpace(ctx context.Context, req *spacev1.DeleteSpaceReq) (*spacev1.DeleteSpaceRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_DELETE, req.GetOperateUid(), ""); err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	if err := service.Space.Delete(invoker.Db.WithContext(ctx), req.GetOperateUid(), req.GetSpaceGuid()); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("delete space fail, err: " + err.Error())
	}
	return &spacev1.DeleteSpaceRes{}, nil
}

func (GrpcServer) ChangeSpaceSort(ctx context.Context, req *spacev1.ChangeSpaceSortReq) (*spacev1.ChangeSpaceSortRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_UPDATE, req.GetOperateUid(), ""); err != nil {
		return nil, err
	}
	if err := service.Space.ChangeSpaceSort(ctx, req.GetSpaceGuid(), req.TargetSpaceGuid, req.DropPosition, req.ParentSpaceGroupGuid); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("change space sort fail, err: " + err.Error())
	}
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_SPACE_CHANGE_SORT,
		Group:      commonv1.LOG_GROUP_SPACE,
		OperateUid: req.GetOperateUid(),
		SpaceGuid:  req.GetSpaceGuid(),
	})
	return &spacev1.ChangeSpaceSortRes{}, nil
}
