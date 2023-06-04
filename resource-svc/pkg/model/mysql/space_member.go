package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/ego-component/egorm"
	"gorm.io/gorm"
)

// SpaceGroupMember 空间角色
type SpaceGroupMember struct {
	Id        int64  `json:"id" gorm:"not null;primary_key;auto_increment"`
	Ctime     int64  `gorm:"not null;default:0;comment:创建时间"`
	Utime     int64  `gorm:"not null;default:0;comment:更新时间"`
	Uid       int64  `gorm:"type:bigint; not null; default:0;unique_index:idx_uid_guid;"`
	Nickname  string `gorm:"type:varchar(191); not null; comment:昵称"`
	Guid      string `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index:idx_uid_guid;comment:唯一标识"`
	CreatedBy int64  `gorm:"type:bigint; not null; default:0;"`
	UpdatedBy int64  `gorm:"type:bigint; not null; default:0;"`
}

func (SpaceGroupMember) TableName() string {
	return "space_group_member"
}

type SpaceGroupMembers []SpaceGroupMember

func (list SpaceGroupMembers) ToPb(userMap map[int64]*userv1.UserInfo) []*commonv1.MemberRole {
	output := make([]*commonv1.MemberRole, 0)
	for _, value := range list {
		output = append(output, value.ToPb(userMap))
	}
	return output
}

func (member SpaceGroupMember) ToPb(userMap map[int64]*userv1.UserInfo) *commonv1.MemberRole {
	var (
		nickname string
		avatar   string
	)
	if userMap != nil {
		userInfo, flag := userMap[member.Uid]
		if flag {
			nickname = userInfo.GetNickname()
			avatar = userInfo.GetAvatar()
		}
	}
	return &commonv1.MemberRole{
		Uid:      member.Uid,
		Nickname: nickname,
		Avatar:   avatar,
	}
}

func (list SpaceGroupMembers) ToMemberPb(userMap map[int64]*userv1.UserInfo, createdUid int64, superAdminMembers PmsSuperAdminMembers) []*commonv1.MemberRole {
	output := make([]*commonv1.MemberRole, 0)
	for _, value := range list {
		output = append(output, value.ToMemberPb(userMap, createdUid, superAdminMembers))
	}
	return output
}

func (member SpaceGroupMember) ToMemberPb(userMap map[int64]*userv1.UserInfo, createdUid int64, superAdminMembers PmsSuperAdminMembers) *commonv1.MemberRole {
	var isSuperAdmin bool
	// todo 理论上有个多个role id，后面在搞吧
	for _, value := range superAdminMembers {
		if member.Uid == value.Uid {
			isSuperAdmin = true
		}
	}
	return toMemberRole(isSuperAdmin, member.Uid, member.Ctime, userMap, createdUid)
}

func SpaceGroupMemberList(db *gorm.DB, spaceGroupGuid string, uids []int64, page *commonv1.Pagination) (respList SpaceGroupMembers, err error) {
	if page.PageSize == 0 || page.PageSize > 200 {
		page.PageSize = 200
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	conds := egorm.Conds{
		"guid": spaceGroupGuid,
	}
	if len(uids) != 0 {
		conds["uid"] = uids
	}
	sql, binds := egorm.BuildQuery(conds)
	query := db.Model(SpaceGroupMember{}).Where(sql, binds...)
	query.Count(&page.Total)
	err = query.Order("id desc ").Offset(int((page.CurrentPage - 1) * page.PageSize)).Limit(int(page.PageSize)).Find(&respList).Error
	return
}

func SpaceGroupMemberCnt(db *gorm.DB, spaceGroupGuid string) (cnt int64, err error) {
	err = db.Model(SpaceGroupMember{}).Where("guid = ?", spaceGroupGuid).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("SpaceGroupMemberCnt fail, err: %w", err)
		return
	}
	return
}

// SpaceGroupMemberCreate 创建一条记录
func SpaceGroupMemberCreate(db *gorm.DB, create *SpaceGroupMember) (err error) {
	err = db.Create(create).Error
	return
}

// SpaceGroupMemberBatchCreate 创建一条记录
func SpaceGroupMemberBatchCreate(db *gorm.DB, list SpaceGroupMembers) (err error) {
	err = db.CreateInBatches(list, len(list)).Error
	return
}

// SpaceGroupMemberBatchDelete 创建一条记录
func SpaceGroupMemberBatchDelete(db *gorm.DB, spaceGroupGuid string, uids []int64) (err error) {
	err = db.Where("guid = ? and uid in (?)", spaceGroupGuid, uids).Delete(&SpaceGroupMember{}).Error
	return
}

func SpaceGroupMemberSearch(db *gorm.DB, spaceGroupGuid, keyword string) (respList SpaceGroupMembers, err error) {
	err = db.Where("guid = ?  and nickname like ?", spaceGroupGuid, "%"+keyword+"%").Limit(10).Find(&respList).Error
	return
}

// SpaceGroupMemberInfo 创建一条记录
func SpaceGroupMemberInfo(db *gorm.DB, field, spaceGuid string, uid int64) (info SpaceGroupMember, err error) {
	err = db.Select(field).Where("guid = ? and  uid = ?", spaceGuid, uid).Find(&info).Error
	return
}

// GetSpaceGroupMemberByBatch 创建一条记录
func GetSpaceGroupMemberByBatch(db *gorm.DB, spaceGuid string, uids []int64) (list SpaceGroupMembers, err error) {
	err = db.Select("uid").Where("guid = ? and  uid in (?)", spaceGuid, uids).Find(&list).Error
	return
}

// SpaceMember 空间角色
type SpaceMember struct {
	Id           int64  `json:"id" gorm:"not null;primary_key;auto_increment"`
	Ctime        int64  `gorm:"not null;default:0;comment:创建时间"`
	Utime        int64  `gorm:"not null;default:0;comment:更新时间"`
	Uid          int64  `gorm:"type:bigint; not null; default:0;unique_index:idx_uid_guid;"`
	Nickname     string `gorm:"type:varchar(191); not null; comment:昵称"`
	Guid         string `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index:idx_uid_guid;comment:唯一标识"`
	CreatedBy    int64  `gorm:"type:bigint; not null; default:0;"`
	UpdatedBy    int64  `gorm:"type:bigint; not null; default:0;"`
	IsSuperAdmin bool   `gorm:"-"`
}

func (SpaceMember) TableName() string {
	return "space_member"
}

type SpaceMembers []SpaceMember

func (list SpaceMembers) ToBasePb(userMap map[int64]*userv1.UserInfo) []*commonv1.MemberRole {
	output := make([]*commonv1.MemberRole, 0)
	for _, value := range list {
		output = append(output, value.ToBasePb(userMap))
	}
	return output
}

func (member SpaceMember) ToBasePb(userMap map[int64]*userv1.UserInfo) *commonv1.MemberRole {
	var (
		nickname string
		avatar   string
	)
	if userMap != nil {
		userInfo, flag := userMap[member.Uid]
		if flag {
			nickname = userInfo.GetNickname()
			avatar = userInfo.GetAvatar()
		}
	}
	return &commonv1.MemberRole{
		Uid:      member.Uid,
		Nickname: nickname,
		Avatar:   avatar,
	}
}

func (list SpaceMembers) ToMemberPb(userMap map[int64]*userv1.UserInfo, createdUid int64, superAdminMembers PmsSuperAdminMembers) []*commonv1.MemberRole {
	output := make([]*commonv1.MemberRole, 0)
	for _, value := range list {
		output = append(output, value.ToMemberPb(userMap, createdUid, superAdminMembers))
	}
	return output
}

func (member SpaceMember) ToMemberPb(userMap map[int64]*userv1.UserInfo, createdUid int64, superAdminMembers PmsSuperAdminMembers) *commonv1.MemberRole {
	// todo 理论上有个多个role id，后面在搞吧
	for _, value := range superAdminMembers {
		if member.Uid == value.Uid {
			member.IsSuperAdmin = true
		}
	}
	return toMemberRole(member.IsSuperAdmin, member.Uid, member.Ctime, userMap, createdUid)
}

func SpaceMemberList(db *gorm.DB, spaceGroupGuid string, uids []int64, page *commonv1.Pagination) (respList SpaceMembers, err error) {
	if page.PageSize == 0 || page.PageSize > 200 {
		page.PageSize = 200
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	conds := egorm.Conds{
		"guid": spaceGroupGuid,
	}
	if len(uids) != 0 {
		conds["uid"] = uids
	}
	sql, binds := egorm.BuildQuery(conds)
	query := db.Model(SpaceMember{}).Where(sql, binds...)
	query.Count(&page.Total)
	err = query.Order("id desc ").Offset(int((page.CurrentPage - 1) * page.PageSize)).Limit(int(page.PageSize)).Find(&respList).Error
	return
}

// SpaceMemberCreate 创建一条记录
func SpaceMemberCreate(db *gorm.DB, create *SpaceMember) (err error) {
	err = db.Create(create).Error
	return
}

// SpaceMemberBatchCreate 创建一条记录
func SpaceMemberBatchCreate(db *gorm.DB, list SpaceMembers) (err error) {
	err = db.CreateInBatches(list, len(list)).Error
	return
}

// SpaceMemberBatchDelete 创建一条记录
func SpaceMemberBatchDelete(db *gorm.DB, spaceGuid string, uids []int64) (err error) {
	err = db.Where("guid = ? and uid in (?)", spaceGuid, uids).Delete(&SpaceMember{}).Error
	return
}

func SpaceMemberSearch(db *gorm.DB, spaceGuid string, keyword string) (respList SpaceMembers, err error) {
	err = db.Where("guid = ? and nickname like ?", spaceGuid, "%"+keyword+"%").Limit(10).Find(&respList).Error
	return
}

func SpaceMemberCnt(db *gorm.DB, spaceGuid string) (cnt int64, err error) {
	err = db.Model(SpaceMember{}).Where("guid = ?", spaceGuid).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("SpaceMemberCnt fail, err: %w", err)
		return
	}
	return
}

// SpaceMemberInfo 创建一条记录
func SpaceMemberInfo(db *gorm.DB, field, spaceGuid string, uid int64) (info SpaceMember, err error) {
	err = db.Select(field).Where("guid = ? and uid = ?", spaceGuid, uid).Find(&info).Error
	return
}

// GetSpaceMemberByBatch 创建一条记录
func GetSpaceMemberByBatch(db *gorm.DB, spaceGuid string, uids []int64) (list SpaceMembers, err error) {
	err = db.Select("uid").Where("guid = ?  and uid in (?)", spaceGuid, uids).Find(&list).Error
	return
}

// GetSpaceMemberBySpaces 根据spaceGuids查询记录
func GetSpaceMemberBySpaces(db *gorm.DB, spaceGuids []string, uid int64) (list SpaceMembers, err error) {
	err = db.Select("guid,uid").Where("guid in (?)  and uid = ?", spaceGuids, uid).Find(&list).Error
	return
}

func GetSpaceMemberIDListByGroupID(tx *gorm.DB, spaceGuid string) ([]int64, error) {
	var groupMemberIDList []int64
	err := tx.Model(SpaceMember{}).Where("guid=?", spaceGuid).Pluck("uid", &groupMemberIDList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberIDList, nil
}
