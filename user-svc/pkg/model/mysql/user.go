package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/samber/lo"
	"gorm.io/gorm"
	sdelete "gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID    int64             `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Ctime int64             `gorm:"bigint;autoCreateTime;comment:创建时间" json:"ctime"`
	Utime int64             `gorm:"bigint;autoUpdateTime;comment:更新时间" json:"utime"`
	Dtime sdelete.DeletedAt `gorm:"bigint;comment:删除时间" json:"dtime"`
}

type User struct {
	Uid           int64             `gorm:"primary_key;auto_increment;comment:主键ID" json:"uid"` // 主键ID
	Ctime         int64             `gorm:"bigint;autoCreateTime;comment:创建时间" json:"ctime"`
	Utime         int64             `gorm:"bigint;autoUpdateTime;comment:更新时间" json:"utime"`
	Dtime         sdelete.DeletedAt `gorm:"bigint;comment:删除时间" json:"dtime"`
	Name          string            `gorm:"type:varchar(255);not null;index:idx_name,unique;comment:名称" json:"name"`
	Nickname      string            `gorm:"type:varchar(255);not null;comment:昵称" json:"nickname"`
	NamePinyin    string            `gorm:"type:varchar(255);not null;default:'';comment:拼音" json:"namePinyin"`
	Avatar        string            `gorm:"not null;default:'';comment:头像" json:"avatar"`
	Password      string            `gorm:"type:text;comment:密码" json:"-"`                      // 密码
	Balance       int               `gorm:"not null;default:0;comment:账户余额" json:"balance"`     // 账户余额
	Level         int               `gorm:"not null;default:0;comment:等级" json:"level"`         // 等级
	Intro         string            `gorm:"not null;default:'';comment:个性签名" json:"intro"`      // 个性签名
	Email         string            `gorm:"not null;default:'';comment:邮箱" json:"email"`        // 邮箱
	EmailStatus   commonv1.USER_EBS `gorm:"not null;default:0;comment:邮箱状态" json:"emailStatus"` // 邮箱
	Phone         string            `gorm:"not null;default:'';comment:电话" json:"phone"`
	RegisterIp    string            `gorm:"not null;default:'';comment:注册IP"`
	LastLoginTime int64             `gorm:"not null;default:0;comment:最后一次登陆时间" json:"lastLoginTime"`
	LastLoginIP   string            `gorm:"not null;default:'';comment:最后一次登陆IP" json:"lastLoginIp"` // 最后一次登陆IP
	Status        int               `gorm:"not null;default:0;comment:状态"`                           // 0 正常，1 被封禁
	Sex           int8              `gorm:"not null;default:0;comment:性别"`                           // 0 有问题， 1 男 ，2 女， 3未知
	Birthday      int64             `gorm:"not null;default:0;comment:生日"`
	Position      string            `gorm:"not null；default:'';comment:职位介绍" json:"position"`
	LastVisitCmt  string            `gorm:"not null;default:''"`
	ActiveTime    int64             `gorm:"bigint;not null; default:0;comment:最后一次活跃时间" json:"activeTime"`
	IsSuperAdmin  bool              `gorm:"not null;default:0;"`
}

func (User) TableName() string {
	return "user"
}

type Users []*User

func (u Users) ToBasePb() []*commonv1.UserBaseInfo {
	return lo.Map(u, func(item *User, index int) *commonv1.UserBaseInfo {
		return &commonv1.UserBaseInfo{
			Uid:      item.Uid,
			Nickname: item.Nickname,
			Avatar:   item.Avatar,
		}
	})
}

func (u Users) ToMap() map[int64]*User {
	res := make(map[int64]*User)
	for _, user := range u {
		res[user.Uid] = user
	}
	return res
}

func (u *User) ToPb() *commonv1.UserInfo {
	user := &commonv1.UserInfo{
		Uid:             u.Uid,
		Ctime:           u.Ctime,
		Email:           u.Email,
		Mobile:          u.Phone,
		Avatar:          u.Avatar,
		NamePinyin:      u.NamePinyin,
		EmailBindStatus: u.EmailStatus,
		Nickname:        u.Nickname,
	}
	return user
}

func UserList(db *gorm.DB, uids []int64) (Users, error) {
	var users Users
	err := db.Model(&User{}).Where("uid in (?)", uids).Find(&users).Error
	return users, err
}

func UserMap(db *gorm.DB, uids []int64) (map[int64]*userv1.UserInfo, error) {
	var users Users
	err := db.Select("nickname,name,avatar,uid,position,active_time").Model(&User{}).Where("uid in (?)", uids).Find(&users).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Map Fail, err: " + err.Error())
	}
	// 响应map数据
	uMap := make(map[int64]*userv1.UserInfo)
	for _, u := range users {
		uMap[u.Uid] = &userv1.UserInfo{
			Name:       u.Name,
			Nickname:   u.Nickname,
			Avatar:     u.Avatar,
			Uid:        u.Uid,
			Position:   u.Position,
			ActiveTime: u.ActiveTime,
		}
	}
	return uMap, err
}

// UserMapByAdminView 管理员看的用户数据，有用户的登录信息
func UserMapByAdminView(db *gorm.DB, uids []int64) (map[int64]*User, error) {
	var users Users
	err := db.Select("uid,nickname,email,avatar,last_login_time").Model(&User{}).Where("uid in (?)", uids).Find(&users).Error
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Map Fail, err: " + err.Error())
	}
	// 响应map数据
	uMap := make(map[int64]*User)
	for _, u := range users {
		uMap[u.Uid] = u
	}
	return uMap, err
}

// UserCreate insert a new regular User into database and returns
func UserCreate(db *gorm.DB, data *User) (err error) {
	if err = db.Create(data).Error; err != nil {
		return fmt.Errorf("UserCreate, err:%w", err)
	}
	return
}

// UidByPhone 根据号码查询用户
func UidByPhone(db *gorm.DB, phone string) (int64, error) {
	var uid int64
	err := db.Model(&User{}).Select("uid").Where("phone = ?", phone).Find(&uid).Error
	return uid, err
}

func UserUpdateNicknameAndName(db *gorm.DB, uid int64, nickname string, name string) (err error) {
	ups := map[string]any{
		"nickname": nickname,
		"name":     name,
	}
	if err = db.Table("user").Where("uid = ?", uid).Updates(ups).Error; err != nil {
		err = fmt.Errorf("UserUpdateNicknameAndName fail, %w", err)
		return
	}
	return
}

// UserUpdatePhone 更新手机号
func UserUpdatePhone(db *gorm.DB, uid int64, phone string) (err error) {
	// 更新手机号
	if err = db.Model(User{}).Where("uid = ?", uid).Update("phone", phone).Error; err != nil {
		err = fmt.Errorf("UserUpdatePhone fail, %w", err)
		return
	}
	return
}

// UserUpdateEmail todo 会有并发问题，先这么处理
func UserUpdateEmail(db *gorm.DB, uid int64, email string) (err error) {

	// 更新手机号
	if err = db.Model(User{}).Where("uid = ?", uid).Update("email", email).Error; err != nil {
		err = fmt.Errorf("UserUpdateEmail fail2, %w", err)
		return
	}
	return
}

func UserUpdate(db *gorm.DB, uid int64, ups map[string]any) (err error) {
	if err = db.Table("user").Where("uid = ?", uid).Updates(ups).Error; err != nil {
		err = fmt.Errorf("user update err, %w", err)
		return
	}
	return
}

func UserSearch(db *gorm.DB, keyword string) (respList Users, err error) {
	err = db.Where("nickname like ?", "%"+keyword+"%").Limit(10).Find(&respList).Error
	return
}

func UserListPage(db *gorm.DB, page *commonv1.Pagination) (respList Users, err error) {
	if page.PageSize == 0 || page.PageSize > 200 {
		page.PageSize = 200
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}

	query := db.Model(User{})
	query.Count(&page.Total)
	err = query.Order("uid desc ").Offset(int((page.CurrentPage - 1) * page.PageSize)).Limit(int(page.PageSize)).Find(&respList).Error
	return
}

func GetUserById(db *gorm.DB, id int64) (*User, error) {
	u := &User{}
	resp := db.Where("uid = ?", id).Find(u)
	return u, resp.Error
}

func GetUserNicknameByUid(db *gorm.DB, uid int64) (string, error) {
	u := &User{}
	err := db.Select("nickname").Where("uid = ?", uid).Find(u).Error
	if err != nil {
		return "", fmt.Errorf("GetUserNicknameByUid fail, err: %w", err)
	}
	return u.Nickname, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	u := &User{}
	resp := db.Where("email = ?", email).Find(u)
	return u, resp.Error
}
