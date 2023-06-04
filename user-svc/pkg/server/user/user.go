package user

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"ecodepost/user-svc/pkg/server/stat"
	"ecodepost/user-svc/pkg/service"

	errcodev1 "ecodepost/pb/errcode/v1"
	statv1 "ecodepost/pb/stat/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/ego-component/egorm"
	"github.com/google/uuid"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
)

type GrpcServer struct{}

var Svc userv1.UserServer = &GrpcServer{}

// ListPage 用户列表(带分页)
func (GrpcServer) ListPage(ctx context.Context, req *userv1.ListPageReq) (*userv1.ListPageRes, error) {
	list, err := mysql.UserListPage(invoker.Db.WithContext(ctx), req.GetPagination())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList fail, err: " + err.Error())
	}

	return &userv1.ListPageRes{
		List:       list.ToBasePb(),
		Pagination: req.GetPagination(),
	}, nil
}

func (GrpcServer) Search(ctx context.Context, req *userv1.SearchReq) (*userv1.SearchRes, error) {
	list, err := mysql.UserSearch(invoker.Db.WithContext(ctx), req.GetNickname())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("SearchSpaceMember Fail, err: " + err.Error())
	}

	return &userv1.SearchRes{
		List: list.ToBasePb(),
	}, nil
}

// LoginInfoByName 根据名称获取登入信息
func (GrpcServer) LoginInfoByName(ctx context.Context, req *userv1.LoginInfoByNameReq) (*userv1.LoginInfoByNameRes, error) {
	var info mysql.User
	err := invoker.Db.WithContext(ctx).Select("uid,password,nickname,name,avatar,status").Where("name = ?", req.GetName()).Find(&info).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("LoginInfoByPhone fail, err: " + err.Error())
	}

	return &userv1.LoginInfoByNameRes{User: &userv1.UserInfo{
		Uid:      info.Uid,
		Password: info.Password,
		Nickname: info.Nickname,
		Avatar:   info.Avatar,
		Name:     info.Name,
		Status:   int32(info.Status),
	}}, nil
}

// AfterUserCreate 执行创建后其他业务操作
func AfterUserCreate(ctx context.Context, u *mysql.User) {
	// 默认收藏夹
	_, err := stat.GrpcServer{}.CollectionGroupCreate(ctx, &statv1.CollectionGroupCreateReq{
		Uid:   u.Uid,
		Title: "默认收藏夹",
		Desc:  "收藏你想要的知识",
	})
	if err != nil {
		elog.Warn("CollectionGroupCreate fail", zap.Error(err))
	}

}

// Create 创建用户
func (GrpcServer) Create(ctx context.Context, req *userv1.CreateReq) (*userv1.CreateRes, error) {
	// 0 <= 随机数<= 30, 创建用户随机一个头像, 比如：https://cdn.gocn.vip/ava/24.png
	cdn := econf.GetString("user-svc.cdn")
	ava := cdn + "/ava/" + strconv.Itoa(rand.Intn(30)) + ".png"
	uObj, err := uuid.NewUUID()
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("uuid fail, err: " + err.Error())
	}
	var pwd string
	if req.Password != nil {
		pwd, err = service.Authorize.Hash(req.GetPassword())
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("authorize hash fail, err: " + err.Error())
		}
	}

	data := &mysql.User{
		Name:       uObj.String(),
		Password:   pwd,
		Phone:      req.GetPhone(),
		RegisterIp: req.GetRegisterIp(),
		Status:     1,
		Avatar:     ava,
	}

	tx := invoker.Db.WithContext(ctx).Begin()
	// WechatUnionId为空，表示用户用手机号直接注册；否则表示用户通过微信注册，然后关联某个手机号
	if req.Phone != "" {
		uid, err := mysql.UidByPhone(tx, req.Phone)
		if err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("UidByPhone fail, err: " + err.Error())
		}
		// 当用户用手机号直接注册时，需要查询手机号对应的用户是否已存在，如果已存在则需要报错
		if req.WechatUnionId == "" && uid != 0 {
			tx.Rollback()
			return nil, errcodev1.ErrProfilePhoneExist().WithMessage("user phone already exists, phone:" + req.Phone)
		}
		// 当用户通过微信注册时，如果试图进行绑定的手机号已关联某个用户，则给data.Uid赋值，同时后面对data.Uid做判断，避免重新创建新用户
		if req.WechatUnionId != "" && uid != 0 {
			data.Uid = uid
		}
	}

	// 只有data.Uid为空时，才创建新用户
	if data.Uid == 0 {
		if err := mysql.UserCreate(tx, data); err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("Create fail, err: " + err.Error())
		}
		// 更新昵称
		nick, _ := invoker.UserGuid.EncodeRandomInt64(data.Uid)
		data.Nickname = "u_" + nick
		if req.Name == "" {
			// 如果name为空，则使用nickname作为name
			data.Name = data.Nickname
		}
		if err := mysql.UserUpdateNicknameAndName(tx, data.Uid, data.Nickname, data.Name); err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("Create fail2, err: " + err.Error())
		}
	}

	// 以下业务逻辑中，data.Uid必定不为空
	// 表示用户通过微信注册，此时需要更新user_open表
	if req.WechatUnionId != "" {
		if err := mysql.UpdateUserOpenUid(tx, data.Uid, req.WechatUnionId); err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("Create fail2, err: " + err.Error())
		}
	}
	tx.Commit()
	AfterUserCreate(ctx, data)
	return &userv1.CreateRes{Uid: data.Uid}, nil
}

// List 根据用户id列表，获取对应用户列表数据，以切片格式记录
func (GrpcServer) List(ctx context.Context, req *userv1.ListReq) (*userv1.ListRes, error) {
	users, err := mysql.UserList(invoker.Db.WithContext(ctx), req.UidList)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("List fail, err: " + err.Error())
	}

	ul := make([]*userv1.UserInfo, 0, len(users))
	for _, u := range users {
		ul = append(ul, &userv1.UserInfo{
			Nickname: u.Nickname,
			Email:    u.Email,
			Avatar:   u.Avatar,
			Uid:      u.Uid,
			Name:     u.Name,
		})
	}
	return &userv1.ListRes{UserList: ul}, nil
}

// Map 根据用户id列表，获取对应用户列表数据，以Map kv 记录
func (GrpcServer) Map(ctx context.Context, req *userv1.MapReq) (*userv1.MapRes, error) {
	uMap, err := mysql.UserMap(invoker.Db.WithContext(ctx), req.UidList)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Map Fail, err: " + err.Error())
	}
	return &userv1.MapRes{UserMap: uMap}, nil
}

// Info 获取用户信息
func (GrpcServer) Info(ctx context.Context, req *userv1.InfoReq) (*userv1.InfoRes, error) {
	if req.GetUid() == 0 && req.GetName() == "" {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("uid and username can't all be empty")
	}
	conds := map[string]any{}
	if req.GetUid() != 0 {
		conds["uid"] = req.GetUid()
	}
	if req.GetName() != "" {
		conds["name"] = req.GetName()
	}
	sql, binds := egorm.BuildQuery(conds)
	var info mysql.User
	err := invoker.Db.WithContext(ctx).Select("uid,nickname,name,email,avatar,status").Where(sql, binds...).Find(&info).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("User info fail, err: " + err.Error())
	}

	return &userv1.InfoRes{User: &userv1.UserInfo{
		Nickname: info.Nickname,
		Email:    info.Email,
		Avatar:   info.Avatar,
		Uid:      info.Uid,
		Name:     info.Name,
	}}, nil
}

// OauthInfo 用户 oauth 信息
func (GrpcServer) OauthInfo(ctx context.Context, req *userv1.OauthInfoReq) (*userv1.OauthInfoRes, error) {
	conds := map[string]any{
		"uid": req.GetUid(),
	}
	sql, binds := egorm.BuildQuery(conds)
	var info mysql.User
	err := invoker.Db.WithContext(ctx).Select("uid,nickname,name,email,avatar,identify_status,cmt_identify_status,apply_max_cmt_cnt").Where(sql, binds...).Find(&info).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Info fail, err: " + err.Error())
	}

	return &userv1.OauthInfoRes{
		Nickname: info.Nickname,
		Email:    info.Email,
		Avatar:   info.Avatar,
		Name:     info.Name,
	}, nil
}

// ProfileInfo 用户 profile 信息
func (GrpcServer) ProfileInfo(ctx context.Context, req *userv1.ProfileInfoReq) (*userv1.ProfileInfoRes, error) {
	if req.GetUid() == 0 && req.GetName() == "" {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("uid and username can't all be empty")
	}

	conds := map[string]any{}

	if req.GetUid() != 0 {
		conds["uid"] = req.GetUid()
	}

	if req.GetName() != "" {
		conds["name"] = req.GetName()
	}

	sql, binds := egorm.BuildQuery(conds)

	var info mysql.User

	err := invoker.Db.WithContext(ctx).Select("uid,nickname,name,email,avatar,identify_status,intro,sex,birthday,ctime").Where(sql, binds...).Find(&info).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Info fail, err: " + err.Error())
	}

	return &userv1.ProfileInfoRes{
		Uid:          info.Uid,
		Nickname:     info.Nickname,
		Email:        info.Email,
		Avatar:       info.Avatar,
		Intro:        info.Intro,
		Sex:          int32(info.Sex),
		Birthday:     info.Birthday,
		Name:         info.Name,
		RegisterTime: info.Ctime,
	}, nil
}

// InfoByPhone 根据手机号获取用户信息
func (GrpcServer) InfoByPhone(ctx context.Context, req *userv1.InfoByPhoneReq) (*userv1.InfoByPhoneRes, error) {
	var info mysql.User
	err := invoker.Db.WithContext(ctx).Select("uid,name,nickname,email,avatar,identify_status").Where("phone = ?", req.GetPhone()).Find(&info).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("InfoByPhone fail, err: " + err.Error())
	}

	return &userv1.InfoByPhoneRes{User: &userv1.UserInfo{
		Nickname: info.Nickname,
		Name:     info.Name,
		Email:    info.Email,
		Avatar:   info.Avatar,
		Uid:      info.Uid,
	}}, nil
}

// LoginInfo 登入信息
func (GrpcServer) LoginInfo(ctx context.Context, req *userv1.LoginInfoReq) (*userv1.LoginInfoRes, error) {
	var info mysql.User
	err := invoker.Db.WithContext(ctx).Select("uid,password,nickname,avatar").Where("uid = ?", req.GetUid()).Find(&info).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("LoginInfo fail, err: " + err.Error())
	}

	return &userv1.LoginInfoRes{
		Uid:      info.Uid,
		Password: info.Password,
		Nickname: info.Nickname,
		Avatar:   info.Avatar,
		Name:     info.Name,
	}, nil
}

// LoginInfoByPhone 根据手机号获取登入信息
func (GrpcServer) LoginInfoByPhone(ctx context.Context, req *userv1.LoginInfoByPhoneReq) (*userv1.LoginInfoByPhoneRes, error) {
	var info mysql.User
	err := invoker.Db.WithContext(ctx).Select("uid,password,nickname,name,avatar,status").Where("phone = ?", req.GetPhone()).Find(&info).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("LoginInfoByPhone fail, err: " + err.Error())
	}

	return &userv1.LoginInfoByPhoneRes{User: &userv1.UserInfo{
		Uid:      info.Uid,
		Password: info.Password,
		Nickname: info.Nickname,
		Avatar:   info.Avatar,
		Name:     info.Name,
		Status:   int32(info.Status),
	}}, nil
}

// LoginUserOpen 用户第三方登入
func (GrpcServer) LoginUserOpen(ctx context.Context, req *userv1.LoginUserOpenReq) (*userv1.LoginUserOpenRes, error) {
	now := time.Now().Unix()
	userOpen := &mysql.UserOpen{
		Ctime:     now,
		Utime:     now,
		UnionId:   req.UnionId,
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Sex:       int(req.Sex),
		Country:   req.Country,
		Province:  req.Province,
		City:      req.City,
		Telephone: req.Telephone,
		Name:      req.Nickname,
	}

	tx := invoker.Db.WithContext(ctx).Begin()
	res, err := service.UserOpenLoginWechat(tx, int(req.Genre), req.OpenId, userOpen)
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrInternal().WithMessage("UserOpenLoginWechat fail, err:" + err.Error())
	}
	tx.Commit()
	return res, nil
}
