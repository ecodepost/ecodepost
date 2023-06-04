package service

import (
	"math/rand"
	"time"

	"ecodepost/user-svc/pkg/model/mysql"

	errcodev1 "ecodepost/pb/errcode/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// // ThirdLoginWechat 微信第三方登录的逻辑
// func ThirdLoginWechat(tx *gorm.DB, clientIp string, openId string, userOpen *mysql.UserOpenLogin) (resp *userv1.LoginThirdResponse, err error) {
// 	resp = &userv1.LoginThirdResponse{}
// 	mysqlUserOpen, err := userOpen.GetUserOpenWechat(openId)
// 	if err != nil {
// 		return
// 	}
//
// 	// 已经存在用户
// 	if mysqlUserOpen.Id > 0 {
// 		// 当数据库的union id为空，微信里的unionId不为空，那么就更新数据
// 		if mysqlUserOpen.UnionId == "" && userOpen.UnionId != "" {
// 			mysqlUserOpen.UnionId = userOpen.UnionId
// 			err = mysqlUserOpen.UpdateUnionId()
// 		}
// 		resp.Nickname = mysqlUserOpen.Nickname
// 		resp.Uid = int64(mysqlUserOpen.Uid)
// 		resp.Avatar = mysqlUserOpen.Avatar
// 		return
// 	}
//
// 	// id为0的时候，表示不存在该用户，需要创建该用户
// 	uniqueName := userOpen.UniqueNickname()
//
// 	// 如果他的union id不为空，那么就可以直接看下有没有对应的union id的用户
// 	if userOpen.UnionId != "" {
// 		var unionUser mysql.UserOpenWechat
// 		err = invoker.Db.Select("id,uid,union_id,nickname,avatar").Where("union_id=?", userOpen.UnionId).Find(&unionUser).Error
// 		// 系统错误
// 		if err != nil {
// 			err = fmt.Errorf("获取unionId失败，err: %w", err)
// 			return
// 		}
//
// 		// 如果存在用户
// 		if unionUser.Id > 0 {
// 			userOpen.UpdateOpenId()
// 			resp.Uid = int64(unionUser.Uid)
// 			resp.Nickname = unionUser.Nickname
// 			resp.Avatar = unionUser.Avatar
// 			return
// 		}
// 	}
// 	// 如果不存在用户，就到下一步，创建该用户
// 	resp, err = thirdLoginCreateUser(tx, userOpen, clientIp, uniqueName)
// 	return
// }
//
// // ThirdLoginGithub Github第三方登录的逻辑
// func ThirdLoginGithub(tx *gorm.DB, clientIp string, targetId int32, userOpen *mysql.UserOpenLogin) (resp *userv1.LoginThirdResponse, err error) {
// 	resp = &userv1.LoginThirdResponse{}
// 	mysqlUserOpen, err := userOpen.GetGithubUser(targetId)
// 	if err != nil {
// 		return
// 	}
//
// 	// 已经存在用户
// 	if mysqlUserOpen.Uid > 0 {
// 		resp.Nickname = mysqlUserOpen.Nickname
// 		resp.Uid = int64(mysqlUserOpen.Uid)
// 		resp.Avatar = mysqlUserOpen.Avatar
// 		return
// 	}
// 	// Uid为0的时候，表示不存在该用户，需要创建该用户
// 	uniqueName := userOpen.UniqueNickname()
// 	// 如果不存在，那么就创建这个用户user open，并且提示需要绑定
// 	resp, err = thirdLoginCreateUser(tx, userOpen, clientIp, uniqueName)
// 	return
// }

// func thirdLoginCreateUser(tx *gorm.DB, userOpen *mysql.UserOpenLogin, clientIp string, uniqueName string) (resp *userv1.LoginThirdResponse, err error) {
// 	resp = &userv1.LoginThirdResponse{}
// 	// 如果不存在，那么就创建这个用户user open，并且提示需要绑定
// 	err = userOpen.Create(tx)
// 	if err != nil {
// 		invoker.Logger.Warnf("create userOpen create error, maybe nickname conflicts, cause: %s", err.Error())
// 		// 如果昵称冲突了，我们用加了open id的唯一名字
// 		userOpen.Name = uniqueName
// 		// 再次尝试
// 		err = userOpen.Create(tx)
// 		if err != nil {
// 			return
// 		}
// 	}
//
// 	u := &mysql.User{
// 		Nickname:   userOpen.Nickname,
// 		Avatar:     userOpen.Avatar,
// 		RegisterIp: clientIp,
// 		State:      1,
// 		Phone:      userOpen.Telephone,
// 		Ctime:      time.Now().Unix(),
// 		Utime:      time.Now().Unix(),
// 	}
// 	err = mysql.UserCreate(tx, u)
// 	if err != nil {
// 		invoker.Logger.Warnf("create mdMembers create error, maybe nickname conflicts, err: %s", err.Error())
// 		// 可能是名字冲突，我们再次尝试
// 		u.Nickname = uniqueName
// 		err = mysql.UserCreate(tx, u)
// 		if err != nil {
// 			return
// 		}
// 	}
//
// 	err = tx.Model(mysql.UserOpen{}).Where("id = ?", userOpen.Id).Updates(map[string]interface{}{
// 		"uid": u.Uid,
// 	}).Error
// 	if err != nil {
// 		return
// 	}
// 	resp.Uid = u.Uid
// 	resp.Nickname = u.Nickname
// 	resp.Avatar = u.Avatar
// 	return
// }

// UserOpenLoginWechat 微信第三方登录的逻辑
func UserOpenLoginWechat(tx *gorm.DB, genre int, openId string, userOpen *mysql.UserOpen) (res *userv1.LoginUserOpenRes, err error) {
	res = &userv1.LoginUserOpenRes{}
	// unionId为空表示参数错误
	if userOpen.UnionId == "" {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("union id is empty")
	}

	mysqlUserOpen, err := userOpen.GetMysqlUserOpen(genre, openId, userOpen.UnionId)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetMysqlUserOpen fail, err:" + err.Error())
	}

	// 1. 如果已经存在userOpen，则返回
	// 外部通过uid是否为0判断是否创建了user表（是否绑定了手机号）
	if mysqlUserOpen.ID > 0 {
		res.Nickname = mysqlUserOpen.Nickname
		res.Uid = int64(mysqlUserOpen.Uid)
		res.Avatar = mysqlUserOpen.Avatar
		res.Phone = mysqlUserOpen.Telephone
		return
	}

	// 2. 如果不存在userOpen，则创建后返回, 此时返回的uid一定为0（用户一定未绑定手机号）
	userOpen.Extra = "{}" // TODO
	if err = userOpen.Create(tx); err != nil {
		elog.Warn("create userOpen create error, maybe nickname conflicts, ", zap.Error(err))
		// 如果昵称冲突了，我们用加了open id的唯一名字
		userOpen.Name = userOpen.UniqueNickname(genre)
		// 再次尝试
		if err = userOpen.Create(tx); err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("userOpen.Create fail, err:" + err.Error())
		}
	}

	// 创建user，更新user_open
	// u := &mysql.User{
	// 	Nickname: userOpen.Nickname,
	// 	Avatar:   userOpen.Avatar,
	// 	Phone:    userOpen.Telephone,
	// 	Ctime:    time.Now().Unix(),
	// 	Utime:    time.Now().Unix(),
	// 	// State:    1,
	// }
	// if err = mysql.UserCreate(tx, u); err != nil {
	// 	elog.Warn("create mdMembers create error, maybe nickname conflicts", zap.Error(err))
	// 	// 可能是名字冲突，我们再次尝试
	// 	u.Nickname = uniqueName
	// 	if err = mysql.UserCreate(tx, u); err != nil {
	// 		return nil, errcodev1.ErrDbError().WithMessage("UserCreate fail, err:" + err.Error())
	// 	}
	// }

	// if err = tx.Model(mysql.UserOpen{}).Where("id = ?", userOpen.ID).Updates(map[string]any{"uid": u.Uid}).Error; err != nil {
	// 	return nil, errcodev1.ErrDbError().WithMessage("update uid fail, err:" + err.Error())
	// }
	// resp.Uid = u.Uid
	// resp.Nickname = u.Nickname
	// resp.Avatar = u.Avatar
	return res, nil
}
