package sso

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	errcodev1 "ecodepost/pb/errcode/v1"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/types"
	"ecodepost/user-svc/pkg/server/user"
	"ecodepost/user-svc/pkg/service"
	"github.com/gotomicro/ego/core/econf"
	"github.com/spf13/cast"

	ssov1 "ecodepost/pb/sso/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/ego-component/eoauth2/server"
)

func (GrpcServer) Login(ctx context.Context, req *ssov1.LoginReq) (res *ssov1.LoginRes, err error) {
	if req.ClientId == "" {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("client id is empty")
	}
	if req.ClientSecret == "" {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("client secret is empty")
	}

	// 直接登录也需要code模式，因为access，需要code。redirect uri没啥用只做校验
	ar := invoker.SsoServer.HandleAuthorizeRequest(ctx, server.AuthorizeRequestParam{
		ClientId:     req.ClientId,
		ResponseType: string(server.CODE),
		RedirectUri:  req.RedirectUri,
	})

	if ar.IsError() {
		err = fmt.Errorf("HandleAuthorizeRequest fail, err: " + ar.GetOutput("error").(string))
		return
	}

	opts := []server.AuthorizeRequestOption{
		server.WithAuthorizeRequestAuthorized(true),
		server.WithAuthorizeSsoUid(req.Uid),
		server.WithAuthorizeSsoPlatform("web"),
		server.WithAuthorizeSsoClientIP(req.ClientIp),
		server.WithAuthorizeSsoUA(req.UserAgent),
	}
	if req.GetParentToken() != "" {
		opts = append(opts, server.WithAuthorizeSsoParentToken(req.GetParentToken()))
	}
	err = ar.Build(opts...)
	if err != nil {
		err = fmt.Errorf("HandleAuthorizeRequest build fail, err: " + err.Error())
		return
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", req.ClientId, req.ClientSecret)))
	tokenAr := invoker.SsoServer.HandleAccessRequest(ctx, server.ParamAccessRequest{
		Method:    "POST",
		GrantType: string(server.AUTHORIZATION_CODE),
		AccessRequestParam: server.AccessRequestParam{
			Code: cast.ToString(ar.GetOutput("code")),
			ClientAuthParam: server.ClientAuthParam{
				Authorization: "Basic " + basicAuth,
			},
		},
	})
	err = tokenAr.Build(
		server.WithAccessRequestAuthorized(true),
		server.WithAccessAuthUA(req.UserAgent),
		server.WithAccessAuthClientIP(req.ClientIp),
	)
	if err != nil {
		err = fmt.Errorf("HandleAccessRequest build fail, err: " + err.Error())
		return
	}

	_, err = user.Svc.Info(ctx, &userv1.InfoReq{
		Uid: req.Uid,
	})
	if err != nil {
		return nil, fmt.Errorf("用户不存在，%s", err.Error())
	}
	t := time.Now().Unix()
	_, err = user.Svc.Update(ctx, &userv1.UpdateReq{
		Uid:           req.Uid,
		LastLoginIp:   &req.ClientIp,
		LastLoginTime: &t,
	})
	if err != nil {
		return nil, fmt.Errorf("user.LoginByUid 更新用户状态失败，err: %w", err)
	}

	return &ssov1.LoginRes{
		Parent: &ssov1.Token{
			Domain:    econf.GetString("oauth.tokenDomain"),
			Token:     ar.GetParentToken().Token,
			ExpiresIn: ar.GetParentToken().ExpiresIn,
			AuthAt:    ar.GetParentToken().AuthAt,
		},
		Sub: &ssov1.Token{
			Domain:    econf.GetString("oauth.tokenDomain"),
			Token:     tokenAr.GetOutput("access_token").(string),
			ExpiresIn: tokenAr.GetOutput("expires_in").(int64),
		},
	}, nil
}

func (GrpcServer) Logout(ctx context.Context, req *ssov1.LogoutReq) (res *ssov1.LogoutRes, err error) {
	err = invoker.TokenComponent.RemoveParentToken(ctx, req.GetParentToken())
	if err != nil {
		return
	}
	return &ssov1.LogoutRes{
		Domain: econf.GetString("oauth.tokenDomain"),
	}, nil
}

func (GrpcServer) Verify(ctx context.Context, req *ssov1.VerifyReq) (res *ssov1.VerifyRes, err error) {
	err = service.Authorize.Verify(req.Uid, req.PasswordHash, req.Password)
	if err != nil {
		return
	}
	return &ssov1.VerifyRes{}, nil
}

func (GrpcServer) GetUserByParentToken(ctx context.Context, req *ssov1.GetUserByParentTokenReq) (res *ssov1.GetUserByParentTokenRes, err error) {
	uid, err := invoker.TokenComponent.GetUidByParentToken(ctx, req.GetParentToken())
	if err != nil {
		return nil, err
	}
	info, err := user.Svc.Info(ctx, &userv1.InfoReq{
		Uid: uid,
	})
	if err != nil {
		return nil, err
	}
	return &ssov1.GetUserByParentTokenRes{
		Uid:      info.User.Uid,
		Nickname: info.User.Nickname,
		Avatar:   info.User.Avatar,
	}, nil
}

// GetToken 根据Code码，获取Access的Token信息
func (GrpcServer) GetToken(ctx context.Context, req *ssov1.GetTokenReq) (resp *ssov1.GetTokenRes, err error) {
	ar := invoker.SsoServer.HandleAccessRequest(ctx, server.ParamAccessRequest{
		Method:    "POST",
		GrantType: string(server.AUTHORIZATION_CODE),
		AccessRequestParam: server.AccessRequestParam{
			Code: req.Code,
			ClientAuthParam: server.ClientAuthParam{
				Authorization: req.Authorization,
			},
			RedirectUri: req.RedirectUri,
		},
	})
	err = ar.Build(
		server.WithAccessRequestAuthorized(true),
		server.WithAccessAuthUA(req.ClientUA),
		server.WithAccessAuthClientIP(req.ClientIP),
	)
	if err != nil {
		return nil, fmt.Errorf("GetAccess error, %w", err)
	}
	resp = &ssov1.GetTokenRes{
		Token:     ar.GetOutput("access_token").(string),
		ExpiresIn: ar.GetOutput("expires_in").(int64),
	}
	return
}

func (GrpcServer) RefreshToken(ctx context.Context, req *ssov1.RefreshTokenReq) (resp *ssov1.RefreshTokenRes, err error) {
	// todo 这里不是每次请求都刷新，根据过期时间，自动判断去做刷新。强制刷新接口后面也可以提供
	// 访问该接口的时间，通常大于1/2过期时间，我们才会触发refresh token操作
	// 例如token需要14天过期，那么这里会判断时间到7天之后，才会触发token
	// 换token操作，需要注意并发问题
	tokenInfo, err := invoker.TokenComponent.GetAPI().GetAllBySubToken(ctx, req.Code)
	if err != nil {
		return
	}

	if tokenInfo.TTL >= econf.GetInt64("oauth.subTokenRefreshTime") {
		return &ssov1.RefreshTokenRes{
			Token: req.Code,
		}, nil
	}

	ar := invoker.SsoServer.HandleAccessRequest(ctx, server.ParamAccessRequest{
		Method:    "POST",
		GrantType: string(server.REFRESH_TOKEN),
		AccessRequestParam: server.AccessRequestParam{
			Code: req.Code,
			ClientAuthParam: server.ClientAuthParam{
				Authorization: req.Authorization,
			},
		},
	})

	err = ar.Build(
		server.WithAccessRequestAuthorized(true),
		server.WithAccessAuthUA(req.ClientUA),
		server.WithAccessAuthClientIP(req.ClientIP),
	)
	if err != nil {
		return nil, fmt.Errorf("GetAccess error, %w", err)
	}
	resp = &ssov1.RefreshTokenRes{
		Token:     ar.GetOutput("access_token").(string),
		ExpiresIn: ar.GetOutput("expires_in").(int64),
	}
	return
}

func (GrpcServer) RemoveToken(ctx context.Context, req *ssov1.RemoveTokenReq) (resp *ssov1.RemoveTokenRes, err error) {
	resp = &ssov1.RemoveTokenRes{}
	err = invoker.TokenComponent.RemoveAllAccess(ctx, req.Token)
	return
}

// GetUserByToken 根据Token信息，获取用户数据
func (GrpcServer) GetUserByToken(ctx context.Context, req *ssov1.GetUserByTokenReq) (resp *ssov1.GetUserByTokenRes, err error) {
	uid, err := invoker.TokenComponent.GetUidByToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	info, err := user.Svc.Info(ctx, &userv1.InfoReq{Uid: uid})
	if err != nil {
		return nil, err
	}
	resp = &ssov1.GetUserByTokenRes{
		Uid:      info.User.Uid,
		Nickname: info.User.Nickname,
		Avatar:   info.User.Avatar,
		Email:    info.User.Email,
		Name:     info.User.Name,
	}
	return
}

func (GrpcServer) GetAccessByPToken(ctx context.Context, req *ssov1.GetAccessByPTokenReq) (*ssov1.GetAccessByPTokenRes, error) {
	if _, ok := types.PlatformMap[req.Platform]; !ok {
		return nil, fmt.Errorf("unknown platform")
	}

	uid, _ := invoker.TokenComponent.GetUidByParentToken(ctx, req.PToken)
	codeAr := invoker.SsoServer.HandleAuthorizeRequest(ctx, server.AuthorizeRequestParam{
		ClientId:     req.ClientId,
		ResponseType: "code",
	})

	err := codeAr.Build(
		server.WithAuthorizeRequestAuthorized(true),
		server.WithAuthorizeSsoParentToken(req.PToken),
		server.WithAuthorizeSsoUid(uid),
		server.WithAuthorizeSsoPlatform(types.PlatformMap[req.Platform]),
		server.WithAuthorizeSsoClientIP(req.ClientIp),
		server.WithAuthorizeSsoUA(req.ClientUa),
	)
	if err != nil {
		return nil, err
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", req.ClientId, req.ClientSecret)))
	tokenAr := invoker.SsoServer.HandleAccessRequest(ctx, server.ParamAccessRequest{
		Method:    "POST",
		GrantType: string(server.AUTHORIZATION_CODE),
		AccessRequestParam: server.AccessRequestParam{
			Code: cast.ToString(codeAr.GetOutput("code")),
			ClientAuthParam: server.ClientAuthParam{
				Authorization: "Basic " + basicAuth,
			},
		},
	})
	err = tokenAr.Build(
		server.WithAccessRequestAuthorized(true),
		server.WithAccessAuthUA(req.ClientUa),
		server.WithAccessAuthClientIP(req.ClientIp),
	)
	if err != nil {
		return nil, err
	}

	return &ssov1.GetAccessByPTokenRes{
		Token:     tokenAr.GetOutput("access_token").(string),
		ExpiresIn: tokenAr.GetOutput("expires_in").(int64),
	}, nil
}

func (GrpcServer) ResetPassword(ctx context.Context, req *ssov1.ResetPasswordReq) (resp *ssov1.ResetPasswordRes, err error) {
	resp = &ssov1.ResetPasswordRes{}
	_, err = user.Svc.Update(ctx, &userv1.UpdateReq{
		Uid:      req.GetUid(),
		Password: &req.Password,
	})
	return
}
