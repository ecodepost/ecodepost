package user

import (
	"context"
	"fmt"

	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"ecodepost/user-svc/pkg/service"
	"ecodepost/user-svc/pkg/util"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	userv1 "ecodepost/pb/user/v1"
)

// Update 更新用户信息
func (GrpcServer) Update(ctx context.Context, req *userv1.UpdateReq) (*userv1.UpdateRes, error) {
	var ups = util.S2a{}
	if req.Nickname != nil {
		ups["nickname"] = *req.Nickname
	}
	if req.Avatar != nil {
		ups["avatar"] = *req.Avatar
	}
	if req.Password != nil {
		var pwdHash string
		pwdHash, err := service.Authorize.Hash(req.GetPassword())
		if err != nil {
			return nil, fmt.Errorf("authorize hash err, %w", err)
		}
		ups["password"] = pwdHash
	}
	if req.LastLoginIp != nil {
		ups["last_login_ip"] = *req.LastLoginIp
	}
	if req.LastLoginTime != nil {
		ups["last_login_time"] = *req.LastLoginTime
	}
	if req.Sex != nil {
		ups["sex"] = *req.Sex
	}
	if req.Birthday != nil {
		ups["birthday"] = *req.Birthday
	}
	if req.Intro != nil {
		ups["intro"] = *req.Intro
	}
	if len(ups) != 0 {
		if err := mysql.UserUpdate(invoker.Db.WithContext(ctx), req.GetUid(), ups); err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("UpdatePassword fail, err: " + err.Error())
		}
	}
	return &userv1.UpdateRes{}, nil
}

// UpdatePhone 更新用户手机号
func (GrpcServer) UpdatePhone(ctx context.Context, req *userv1.UpdatePhoneReq) (*userv1.UpdatePhoneRes, error) {
	// todo 会有并发问题，先这么处理
	userInfo := mysql.User{}
	// 如果不是本用户，但手机号已经存在，那么无法更新手机号
	err := invoker.Db.WithContext(ctx).Select("uid").Where("phone = ? and uid != ?", req.GetPhone(), req.GetUid()).Find(&userInfo).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdatePhone fail, err: " + err.Error())
	}

	if userInfo.Uid > 0 {
		return nil, errcodev1.ErrProfilePhoneExist()
	}

	err = mysql.UserUpdatePhone(invoker.Db.WithContext(ctx), req.GetUid(), req.GetPhone())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdatePhone fail2, err: " + err.Error())
	}
	return &userv1.UpdatePhoneRes{}, nil
}

// UpdateEmail 更新用户邮箱
func (GrpcServer) UpdateEmail(ctx context.Context, req *userv1.UpdateEmailReq) (*userv1.UpdateEmailRes, error) {
	// todo 会有并发问题，先这么处理
	userInfo := mysql.User{}
	// 如果不是本用户，但手机号已经存在，那么无法更新手机号
	err := invoker.Db.WithContext(ctx).Select("uid").Where("email = ? and email_status = ? and uid != ?", req.GetEmail(), commonv1.USER_EBS_CONFIRMED, req.GetUid()).Find(&userInfo).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdateEmail fail, err: " + err.Error())
	}

	if userInfo.Uid > 0 {
		return nil, errcodev1.ErrProfileEmailExist()
	}

	err = mysql.UserUpdateEmail(invoker.Db.WithContext(ctx), req.GetUid(), req.GetEmail())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdateEmail fail2, err: " + err.Error())
	}
	return &userv1.UpdateEmailRes{}, nil
}
