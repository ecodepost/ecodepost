package ssoservice

import (
	"context"
	"fmt"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/util"
	commonv1 "ecodepost/pb/common/v1"
	ssov1 "ecodepost/pb/sso/v1"
	userv1 "ecodepost/pb/user/v1"
)

type LoginType string

const (
	BindPhoneType             LoginType = "bind_phone"
	RegisterLoginType         LoginType = "register"
	PasswordLoginType         LoginType = "password"
	SmsLoginType              LoginType = "sms"
	RetrievePasswordLoginType LoginType = "retrieve_password"
	ResetPasswordLoginType    LoginType = "reset_password"
	CheckLoginType            LoginType = "check_login"
	WechatType                LoginType = "login_wechat"
)

type user struct {
}

func InitUser() *user {
	return &user{}
}

// RequestValid todo后期多了可以换成多种类型，使用interface实现
func (u *user) RequestValid(ctx context.Context, req dto.OauthRequest, loginType LoginType) (showTips string, err error) {
	switch loginType {
	case RegisterLoginType:
		if req.Account == "" {
			showTips = "账号不能为空"
			err = fmt.Errorf(showTips)
			return
		}

		if !util.CheckMobile(req.Account) {
			showTips = "不是正确的手机号"
			err = fmt.Errorf(showTips)
			return
		}

		// todo 判断长度、只允许数字，字母？
		if req.Password == "" {
			showTips = "密码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		if req.Code == "" {
			showTips = "手机验证码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		var userInfo *userv1.InfoByPhoneRes
		userInfo, err = invoker.GrpcUser.InfoByPhone(ctx, &userv1.InfoByPhoneReq{Phone: req.Account})
		if err != nil {
			showTips = "系统错误"
			err = fmt.Errorf("user register get user info failed,err: %w", err)
			return
		}
		if userInfo.User.Uid > 0 {
			showTips = "该用户已经注册"
			err = fmt.Errorf("该用户已经注册")
			return
		}
	case RetrievePasswordLoginType:
		if req.Account == "" {
			showTips = "账号不能为空"
			err = fmt.Errorf(showTips)
			return
		}

		if !util.CheckMobile(req.Account) {
			showTips = "不是正确的手机号"
			err = fmt.Errorf(showTips)
			return
		}

		// todo 判断长度、只允许数字，字母？
		if req.Password == "" {
			showTips = "密码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		if req.Code == "" {
			showTips = "手机验证码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
	case ResetPasswordLoginType:
		// todo 判断长度、只允许数字，字母？
		if req.OldPassword == "" {
			showTips = "旧密码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		if req.Password == "" {
			showTips = "新密码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
	case PasswordLoginType:
		if req.Account == "" {
			showTips = "账号不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		if !util.CheckMobile(req.Account) {
			showTips = "不是正确的手机号"
			err = fmt.Errorf(showTips)
			return
		}
		// todo 判断长度、只允许数字，字母？
		if req.Password == "" {
			showTips = "密码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
	case SmsLoginType:
		if req.Account == "" {
			showTips = "账号不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		if !util.CheckMobile(req.Account) {
			showTips = "不是正确的手机号"
			err = fmt.Errorf(showTips)
			return
		}
		if req.Code == "" {
			showTips = "手机验证码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
	case BindPhoneType:
		if req.Account == "" {
			showTips = "账号不能为空"
			err = fmt.Errorf(showTips)
			return
		}

		if !util.CheckMobile(req.Account) {
			showTips = "不是正确的手机号"
			err = fmt.Errorf(showTips)
			return
		}

		// todo 判断长度、只允许数字，字母？
		if req.Password == "" {
			showTips = "密码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		if req.Code == "" {
			showTips = "手机验证码不能为空"
			err = fmt.Errorf(showTips)
			return
		}
		var userInfo *userv1.InfoByPhoneRes
		userInfo, err = invoker.GrpcUser.InfoByPhone(ctx, &userv1.InfoByPhoneReq{
			Phone: req.Account,
		})
		if err != nil {
			showTips = "系统错误"
			err = fmt.Errorf("user register get user info failed,err: %w", err)
			return
		}
		if userInfo.User.Uid > 0 {
			showTips = "该用户已经注册"
			err = fmt.Errorf("该用户已经注册")
			return
		}
	default:
		showTips = "不存在的登录类型"
		err = fmt.Errorf(showTips)
		return
	}
	return
}

// Verify 目前只有手机号
func (u *user) Verify(ctx context.Context, req dto.OauthRequest, loginType LoginType) (userInfo *userv1.LoginInfoByPhoneRes, showTips string, err error) {
	switch loginType {
	case RetrievePasswordLoginType:
		fallthrough
	case RegisterLoginType:
		fallthrough
	case ResetPasswordLoginType:
		fallthrough
	case PasswordLoginType:
		userInfo, err = invoker.GrpcUser.LoginInfoByPhone(ctx, &userv1.LoginInfoByPhoneReq{
			Phone: req.Account,
		})
		if err != nil {
			showTips = "系统错误"
			err = fmt.Errorf("user PasswordLoginType verify get user info failed,err: %w", err)
			return
		}

		if userInfo.User.GetUid() == 0 {
			showTips = "用户不存在"
			err = fmt.Errorf("用户不存在")
			return
		}

		if userInfo.User.GetPassword() == "" {
			showTips = "当前用户没有设置密码"
			err = fmt.Errorf("当前用户没有设置密码")
			return
		}

		if userInfo.User.GetStatus() == int32(commonv1.USER_STATUS_BAN) {
			showTips = "您的账户已被封禁，请联系管理员尝试解除封禁"
			err = fmt.Errorf("请联系管理员尝试解除封禁")
			return
		}

		_, err = invoker.GrpcSso.Verify(ctx, &ssov1.VerifyReq{
			Uid:          userInfo.User.GetUid(),
			Password:     req.Password,
			PasswordHash: userInfo.User.GetPassword(),
		})
		if err != nil {
			showTips = "验证密码失败"
			err = fmt.Errorf("user PasswordLoginType verify password failed, err: %w", err)
			return
		}
	case SmsLoginType:
		userInfo, err = invoker.GrpcUser.LoginInfoByPhone(ctx, &userv1.LoginInfoByPhoneReq{
			Phone: req.Account,
		})

		if err != nil {
			showTips = "系统错误"
			err = fmt.Errorf("user PasswordLoginType verify get user info failed,err: %w", err)
			return
		}

		if userInfo.User.GetUid() == 0 {
			showTips = "用户不存在"
			err = fmt.Errorf("用户不存在")
			return
		}

		if userInfo.User.GetStatus() == int32(commonv1.USER_STATUS_BAN) {
			showTips = "您的账户已被封禁，请联系管理员尝试解除封禁"
			err = fmt.Errorf("请联系管理员尝试解除封禁")
			return
		}

		var verifyShowTips string
		verifyShowTips, err = Code.VerifyCode(ctx, req.Account, req.Code, LoginCodeType)
		if err != nil {
			showTips = verifyShowTips
			err = fmt.Errorf("验证短信失败, err: %w", err)
			return
		}
	default:
		showTips = "不存在的登录类型"
		err = fmt.Errorf(showTips)
		return
	}
	return
}
