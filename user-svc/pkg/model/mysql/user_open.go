package mysql

import (
	"errors"

	"gorm.io/gorm"
	sdelete "gorm.io/plugin/soft_delete"

	"ecodepost/user-svc/pkg/constx"
	"ecodepost/user-svc/pkg/invoker"
)

type UserOpen struct {
	ID           int64             `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Ctime        int64             `gorm:"bigint;autoCreateTime;comment:创建时间" json:"ctime"` // 创建时间
	Utime        int64             `gorm:"bigint;autoUpdateTime;comment:更新时间" json:"utime"` // 更新时间
	Dtime        sdelete.DeletedAt `gorm:"bigint;comment:删除时间" json:"dtime"`                // 删除时间
	Genre        int               `gorm:"not null;" json:"genre"`                          // 类型： 1 wechat web，4 wechat miniprogram，5 github
	Uid          int               `gorm:"not null;" json:"uid"`                            // GoCN uid
	Extra        string            `gorm:"not null;type:json"`                              // 额外数据
	TargetId     int               `gorm:"nut null" json:"targetId"`                        // 如果是github，那么target id为github的uid信息
	H5OpenId     string            `gorm:"not null;index" json:"h5OpenId"`                  // openId (如果类型是微信 则代表公众号openId)
	WebOpenId    string            `gorm:"not null;index" json:"webOpenId"`                 // web_open_id (如果类型是微信 则代表webOpenId)
	AppOpenId    string            `gorm:"not null;" json:"appOpenId"`                      // app_open_id (如果类型是微信 则代表开放平台openId)
	MiniOpenId   string            `gorm:"not null;" json:"miniOpenId"`                     // mini_open_id (如果类型是微信 则代表小程序openId)
	UnionId      string            `gorm:"not null;" json:"unionId"`                        // unionId
	AccessToken  string            `gorm:"not null;" json:"accessToken"`                    // access_token
	ExpiresIn    int               `gorm:"not null;" json:"expiresIn"`                      // access_token过期时间
	RefreshToken string            `gorm:"not null;" json:"refreshToken"`                   // access_token过期可用该字段刷新用户access_token
	Scope        string            `gorm:"not null;" json:"scope"`                          // 应用授权作用域
	Nickname     string            `gorm:"not null;" json:"nickname"`                       // 用户来源平台的昵称
	Avatar       string            `gorm:"not null;" json:"avatar"`                         // 头像
	Sex          int               `gorm:"not null;" json:"sex"`                            // 性别[1男 2女]
	Country      string            `gorm:"not null;" json:"country"`                        // 国家
	Province     string            `gorm:"not null;" json:"province"`                       // 省份
	City         string            `gorm:"not null;" json:"city"`                           // 城市
	State        int               `gorm:"not null;" json:"state"`                          // 是否绑定主帐号[默认0否 1是]
	Telephone    string            `gorm:"not null;" json:"telephone"`                      // 电话
	Email        string            `gorm:"not null;" json:"email"`                          // 邮箱
	Name         string            `gorm:"not null;UNIQUE_INDEX" json:"name"`               // 用户在所有平台唯一昵称
	Intro        string            `gorm:"not null;" json:"intro"`
}

func (*UserOpen) TableName() string {
	return "user_open"
}

type UserOpenGithub struct {
	Uid      int    `gorm:"not null;" json:"uid"`      // GoCN uid
	Nickname string `gorm:"not null;" json:"nickname"` // 用户来源平台的昵称
	Avatar   string `gorm:"not null;" json:"avatar"`   // 头像
}

func (*UserOpenGithub) TableName() string {
	return "user_open"
}

func GetUserOpenByUid(db *gorm.DB, uid int) (*UserOpen, error) {
	uo := &UserOpen{}
	resp := db.Where("uid = ? ", uid).Find(uo)
	return uo, resp.Error
}

// GetMysqlUserOpen 是否存在该用户 并将open id赋值
func (u *UserOpen) GetMysqlUserOpen(genre int, openId string, unionId string) (resp *UserOpen, err error) {
	if genre == 0 {
		return nil, errors.New("genre cant 0")
	}
	if openId == "" {
		return nil, errors.New("openId cant empty")
	}

	err = invoker.Db.Where("union_id = ?", unionId).First(&resp).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	ups := make(map[string]any)
	switch genre {
	case constx.UserOpenWechatWeb:
		if resp.WebOpenId == "" {
			ups["web_open_id"] = openId
		}
		u.WebOpenId = openId
	case constx.UserOpenWechatH5:
		if resp.H5OpenId == "" {
			ups["h5_open_id"] = openId
		}
		u.H5OpenId = openId
	// case constx.UserOpenWechatMini:
	// 	if resp.MiniOpenId == "" {
	// 		ups["mini_open_id"] = openId
	// 	}
	// 	u.MiniOpenId = openId
	default:
		return nil, errors.New("not exist type wechat")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if len(ups) != 0 {
		if e := invoker.Db.Model(&UserOpen{}).Where("union_id = ?", unionId).Updates(ups).Error; e != nil {
			return nil, e
		}
	}
	return
}

func (u *UserOpen) UniqueNickname(genre int) (uniqueName string) {
	switch genre {
	case constx.UserOpenWechatWeb:
		uniqueName = u.Nickname + "_" + u.WebOpenId
	case constx.UserOpenWechatMini:
		uniqueName = u.Nickname + "_" + u.MiniOpenId
	}
	return
}

// UpdateOpenId 根据来源更新 open id
// func (u *UserOpen) UpdateOpenId(genre int) (err error) {
// 	if u.Genre == 0 {
// 		err = errors.New("genre cant 0")
// 		return
// 	}
//
// 	switch genre {
// 	case constx.UserOpenWechatWeb:
// 		if u.WebOpenId == "" {
// 			err = errors.New("mini_open_id cant empty")
// 			return
// 		}
// 		// 2.获取数据库信息
// 		err = invoker.Db.Where("union_id=?", u.UnionId).Updates(map[string]interface{}{
// 			"web_open_id": u.WebOpenId,
// 		}).Error
// 		if err != nil {
// 			return
// 		}
// 	case constx.UserOpenWechatMini:
// 		if u.MiniOpenId == "" {
// 			err = errors.New("mini_open_id cant empty")
// 			return
// 		}
// 		// 2.获取数据库信息
// 		err = invoker.Db.Where("union_id=?", u.UnionId).Updates(map[string]interface{}{
// 			"web_open_id": u.MiniOpenId,
// 		}).Error
// 		if err != nil {
// 			return
// 		}
// 	default:
// 		err = errors.New("not exist type wechat")
// 		return
// 	}
// 	return
// }

func UpdateUserOpenUid(db *gorm.DB, uid int64, unionId string) error {
	return invoker.Db.Model(&UserOpen{}).Where("union_id = ?", unionId).Update("uid", uid).Error
}

func (u *UserOpen) UpdateUnionId() error {
	return invoker.Db.Model(u).Where("id = ?", u.ID).Update("union_id", u.UnionId).Error
}

// Create 新增一条记
func (g *UserOpen) Create(db *gorm.DB) error {
	return db.Create(g).Error
}
