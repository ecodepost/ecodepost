package mysql

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"

	commonv1 "ecodepost/pb/common/v1"
	spacev1 "ecodepost/pb/space/v1"

	"github.com/gotomicro/ekit/slice"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Space struct {
	Id                    int64               `json:"id" gorm:"not null;primary_key;auto_increment"`
	SpaceGroupGuid        string              `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:spaceGroupGuid"`
	Guid                  string              `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	Name                  string              `gorm:"not null;default:'';comment:空间名称"`
	Ctime                 int64               `gorm:"not null;default:0;comment:创建时间"`
	Utime                 int64               `gorm:"not null;default:0;comment:更新时间"`
	Dtime                 int64               `gorm:"not null;default:0;comment:删除时间"`
	CreatedBy             int64               `gorm:"not null;default:0;comment:创建人"`
	UpdatedBy             int64               `gorm:"not null;default:0;comment:更新人"`
	DeletedBy             int64               `gorm:"not null;default:0;comment:删除人"`
	IconType              commonv1.FILE_IT    `gorm:"not null;type:smallint;default:1;comment:图标类型"`
	Icon                  string              `gorm:"not null;type:varchar(255);default:'';comment:图标"`
	Sort                  int64               `gorm:"not null;type:bigint;default:0;comment:排序值"`
	Visibility            commonv1.CMN_VISBL  `gorm:"not null;default:0;comment:可见级别"`
	Type                  commonv1.CMN_APP    `gorm:"not null;type:int;default:0;comment:空间类型"`
	Layout                commonv1.SPC_LAYOUT `gorm:"not null;type:int;default:0;comment:布局类型"`
	Access                commonv1.SPC_ACS    `gorm:"not null;type:int;default:0;comment:访问类型"`
	IsAllowReadMemberList bool                `gorm:"not null;default:false;comment:是否开启空间成员读取用户列表"` // 如果打开，属于这个空间下的用户，可以看到用户列表
	ChargeType            commonv1.SPC_CT     `gorm:"not null;default:0;comment:收费类型"`
	OriginPrice           int64               `gorm:"not null;default:0;comment:原价"`
	Price                 int64               `gorm:"not null;default:0;comment:现价"`
	Desc                  string              `gorm:"not null;type:varchar(1024);default:'';comment:空间介绍和说明"`
	HeadImage             string              `gorm:"not null;type:varchar(191);default:'';comment:头图"`
	Cover                 string              `gorm:"not null;type:varchar(191);default:'';comment:封面"`
	Link                  string              `gorm:"not null;type:varchar(255);default:'';comment:链接"`
	OptionList            SpaceOptions        `gorm:"-"`
}

func (Space) TableName() string {
	return "space"
}

func (s *Space) ToPb() *commonv1.SpaceInfo {
	return &commonv1.SpaceInfo{
		Guid:         s.Guid,
		Name:         s.Name,
		IconType:     s.IconType,
		Icon:         s.Icon,
		SpaceType:    s.Type,
		SpaceLayout:  s.Layout,
		Visibility:   s.Visibility,
		ChargeType:   s.ChargeType,
		OriginPrice:  s.OriginPrice,
		Price:        s.Price,
		Desc:         s.Desc,
		HeadImage:    s.HeadImage,
		Cover:        s.Cover,
		Access:       s.Access,
		SpaceOptions: s.OptionList.ToPb(s.Type),
	}
}

func (s *Space) ToTreePb() *spacev1.TreeSpace {
	return &spacev1.TreeSpace{
		Guid:           s.Guid,
		Name:           s.Name,
		SpaceGroupGuid: s.SpaceGroupGuid,
		Icon:           s.Icon,
		SpaceType:      s.Type,
		SpaceLayout:    s.Layout,
		Visibility:     s.Visibility,
		SpaceOptions:   s.OptionList.ToPb(s.Type),
		ChargeType:     s.ChargeType,
		OriginPrice:    s.OriginPrice,
		Price:          s.Price,
		Desc:           s.Desc,
		HeadImage:      s.HeadImage,
		Cover:          s.Cover,
		Access:         s.Access,
		Link:           s.Link,
	}
}

func (s *Space) ToPbWithCnt(ctx context.Context) *commonv1.SpaceInfo {
	cnt, _ := SpaceMemberCnt(invoker.Db.WithContext(ctx), s.Guid)
	res := s.ToPb()
	res.MemberCnt = cnt
	return res
}

type Spaces []*Space

func (list Spaces) ToGuids() []string {
	return slice.Map(list, func(idx int, e *Space) string {
		return e.Guid
	})
}

func (list Spaces) FindByGuid(guid string) *Space {
	return list.Find(func(e *Space) bool {
		return e.Guid == guid
	})
}

func (list Spaces) Find(fn func(e *Space) bool) *Space {
	for _, spaceInfo := range list {
		if fn(spaceInfo) {
			return spaceInfo
		}
	}
	return nil
}

func (list Spaces) Guids() []string {
	return lo.Map(list, func(s *Space, i int) string { return s.Guid })
}

// func (list Spaces) WithTopics(ctx context.Context) Spaces {
// 	topics, err := BatchListTopicFromCache(ctx, list.Guids())
// 	if err != nil {
// 		elog.Warn("BatchListTopicFromCache fail")
// 		return list
// 	}
// 	return lo.Map(list, func(s *Space, i int) *Space {
// 		s.Topics = topics[s.Guid]
// 		return s
// 	})
// }

func (list Spaces) ToMap() map[string]Space {
	output := make(map[string]Space)
	for _, value := range list {
		output[value.Guid] = *value
	}
	return output
}

func (list Spaces) ToPb() []*commonv1.SpaceInfo {
	output := make([]*commonv1.SpaceInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPb())
	}
	return output
}

func (list Spaces) ToPbWithCnt(ctx context.Context) []*commonv1.SpaceInfo {
	output := make([]*commonv1.SpaceInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPbWithCnt(ctx))
	}
	return output
}

type SpaceGroup struct {
	Id                    int64              `json:"id" gorm:"not null;primary_key;auto_increment"`
	Guid                  string             `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	Name                  string             `gorm:"not null; default:''; comment:团队名称"`
	Ctime                 int64              `gorm:"not null;default:0;comment:创建时间"`
	Utime                 int64              `gorm:"not null;default:0;comment:更新时间"`
	Dtime                 int64              `gorm:"not null;default:0;comment:删除时间"`
	CreatedBy             int64              `gorm:"created_by; not null; default:0; comment:创建人"`
	UpdatedBy             int64              `gorm:"updated_by; not null; default:0; comment:更新人"`
	DeletedBy             int64              `gorm:"deleted_by; not null; default:0; comment:删除人"`
	IconType              commonv1.FILE_IT   `gorm:"type:smallint; not null; default:1; comment:图标类型"`
	Icon                  string             `gorm:"type:varchar(255); not null; default:''; comment:图标"`
	Sort                  int64              `gorm:"type:bigint; not null; default:0; comment:排序值"`
	List                  Spaces             `gorm:"-"`
	Visibility            commonv1.CMN_VISBL `gorm:"not null; default:0; comment:可见级别"`
	IsAllowReadMemberList bool               `gorm:"not null; default:false; comment:是否开启成员读取用户列表"` // 如果打开，属于这个分组下的用户，可以看到用户列表
}

func (SpaceGroup) TableName() string {
	return "space_group"
}

type SpaceGroups []*SpaceGroup

func (list SpaceGroups) ToTreePb() []*spacev1.TreeSpaceGroup {
	output := make([]*spacev1.TreeSpaceGroup, 0)
	for _, value := range list {
		output = append(output, value.ToTreePb())
	}
	return output
}

func (s *SpaceGroup) ToPb() *spacev1.SpaceGroupInfo {
	return &spacev1.SpaceGroupInfo{
		Guid:       s.Guid,
		Name:       s.Name,
		Visibility: s.Visibility,
	}
}

func (s *SpaceGroup) ToTreePb() *spacev1.TreeSpaceGroup {
	return &spacev1.TreeSpaceGroup{
		Guid:       s.Guid,
		Name:       s.Name,
		Visibility: s.Visibility,
	}
}

func (list SpaceGroups) FindByGuid(guid string) *SpaceGroup {
	return list.Find(func(e *SpaceGroup) bool {
		return e.Guid == guid
	})
}
func (list SpaceGroups) Find(fn func(e *SpaceGroup) bool) *SpaceGroup {
	for _, spaceInfo := range list {
		if fn(spaceInfo) {
			return spaceInfo
		}
	}

	return nil
}

func (list SpaceGroups) ToMap() map[string]SpaceGroup {
	output := make(map[string]SpaceGroup)
	for _, value := range list {
		output[value.Guid] = *value
	}
	return output
}

// SpaceGetInfoByGuid 创建一条记录
func SpaceGetInfoByGuid(db *gorm.DB, field string, guid string) (info Space, err error) {
	err = db.Select(field).Where("guid = ?", guid).Find(&info).Error
	if err != nil {
		err = fmt.Errorf("SpaceGetInfoByGuid fail, err: %w", err)
		return
	}
	return
}

// SpaceGetInfoByInGuids 创建一条记录
func SpaceGetInfoByInGuids(db *gorm.DB, field string, guids []string) (list Spaces, err error) {
	err = db.Select(field).Where("guid in (?)", guids).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceGetInfoByInGuids fail, err: %w", err)
	}
	return
}

// SpaceGroupGetInfoByInGuids 创建一条记录
func SpaceGroupGetInfoByInGuids(db *gorm.DB, field string, guids []string) (list SpaceGroups, err error) {
	err = db.Select(field).Where("guid in (?)", guids).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceGroupGetInfoByInGuids fail, err: %w", err)
	}
	return
}

func TemplateSpaceList(db *gorm.DB, cmtGuid string) (spaceGroupList SpaceGroups, err error) {
	// 查到所有的space group list
	err = db.Where("cmt_guid = ? and dtime = ?", cmtGuid, 0).Order("`sort` asc").Find(&spaceGroupList).Error
	if err != nil {
		return nil, fmt.Errorf("space tree get space group list fail, err: %w", err)
	}
	var spaceList Spaces
	err = db.Where("cmt_guid = ? and dtime = ?", cmtGuid, 0).Order("`sort` asc").Find(&spaceList).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceListByUser fail, err: %w", err)
	}

	groupGuids := make([]string, 0)
	for _, value := range spaceGroupList {
		groupGuids = append(groupGuids, value.Guid)
	}

	// 没有数据
	if len(groupGuids) == 0 {
		return spaceGroupList, nil
	}

	spaceGroupMapToSpaceArr := make(map[string]Spaces)
	for _, value := range spaceList {
		arr, flag := spaceGroupMapToSpaceArr[value.SpaceGroupGuid]
		if !flag {
			spaceGroupMapToSpaceArr[value.SpaceGroupGuid] = make(Spaces, 0)
		}
		spaceGroupMapToSpaceArr[value.SpaceGroupGuid] = append(arr, value)
	}
	for _, value := range spaceGroupList {
		arr, flag := spaceGroupMapToSpaceArr[value.Guid]
		if flag {
			value.List = arr
		}
	}
	return
}

// SpaceAndGroupList 根据用户，获得space list， group list信息
func SpaceAndGroupList(db *gorm.DB, uid int64) (newSpaceList Spaces, spaceGroupList SpaceGroups, err error) {
	err = db.Where("dtime = ?", 0).Order("`sort` asc").Find(&spaceGroupList).Error
	spaceList, err := SpaceListByUser(db, uid)
	if err != nil {
		return nil, nil, fmt.Errorf("SpaceAndGroupList get space list fail, err: %w", err)
	}
	// 获取全部space的space option
	spaceGuids := lo.Map(spaceList, func(t *Space, i int) string {
		return t.Guid
	})
	spaceOptions, err := BatchGetSpaceOptionList(db, spaceGuids)
	if err != nil {
		return nil, nil, fmt.Errorf("SpaceAndGroupList get option list fail, err: %w", err)
	}

	newSpaceList = lo.Map(spaceList, func(oldSpace *Space, i int) *Space {
		mySpaceOptions := make(SpaceOptions, 0)
		for _, value := range spaceOptions {
			if value.Guid == oldSpace.Guid {
				mySpaceOptions = append(mySpaceOptions, value)
			}
		}
		oldSpace.OptionList = mySpaceOptions
		return oldSpace
	})
	return newSpaceList, spaceGroupList, nil

}

func SpaceTree(db *gorm.DB, uid int64) (spaceGroupList SpaceGroups, err error) {
	err = db.Where(" dtime = ?", 0).Order("`sort` asc").Find(&spaceGroupList).Error
	spaceList, err := SpaceListByUser(db, uid)
	if err != nil {
		return nil, fmt.Errorf("space tree get space list fail, err: %w", err)
	}

	// 获取全部space的space option
	spaceGuids := lo.Map(spaceList, func(t *Space, i int) string {
		return t.Guid
	})
	spaceOptions, err := BatchGetSpaceOptionList(db, spaceGuids)
	if err != nil {
		return nil, fmt.Errorf("space tree get option list fail, err: %w", err)
	}

	newSpaceList := lo.Map(spaceList, func(oldSpace *Space, i int) *Space {
		mySpaceOptions := make(SpaceOptions, 0)
		for _, value := range spaceOptions {
			if value.Guid == oldSpace.Guid {
				mySpaceOptions = append(mySpaceOptions, value)
			}
		}
		oldSpace.OptionList = mySpaceOptions
		return oldSpace
	})

	spaceGroupMapToSpaceArr := make(map[string]Spaces)
	for _, value := range newSpaceList {
		arr, flag := spaceGroupMapToSpaceArr[value.SpaceGroupGuid]
		if !flag {
			spaceGroupMapToSpaceArr[value.SpaceGroupGuid] = make(Spaces, 0)
		}
		spaceGroupMapToSpaceArr[value.SpaceGroupGuid] = append(arr, value)
	}
	for _, value := range spaceGroupList {
		arr, flag := spaceGroupMapToSpaceArr[value.Guid]
		if flag {
			value.List = arr
		}
	}
	return
}

// SpaceListByUser 根据用户的权限，返回不同的空间列表
func SpaceListByUser(db *gorm.DB, uid int64) (spaceList Spaces, err error) {
	spaceList = make(Spaces, 0)
	spaceList, err = SpacePublicList(db)
	if err != nil {
		return nil, fmt.Errorf("SpaceListByUser fail, err: %w", err)
	}

	visibityArr := []int64{int64(commonv1.CMN_VISBL_DRAFT.Number()), int64(commonv1.CMN_VISBL_SECRET.Number())}
	spaceSecretList := make(Spaces, 0)
	if uid != 0 {
		// 用户自己创建的私有的可以看到
		err = db.Where("created_by = ? and dtime = ? and visibility in (?)", uid, 0, visibityArr).Order("`sort` asc").Find(&spaceSecretList).Error
		if err != nil {
			return nil, fmt.Errorf("SpaceListByUser fail2, err: %w", err)
		}
		spaceList = append(spaceList, spaceSecretList...)
	}
	return
}

// SpaceListByUserAndGroups 根据用户的权限，返回不同的空间列表
func SpaceListByUserAndGroups(db *gorm.DB, uid int64, groupGuids []string) (spaceList Spaces, err error) {
	spaceList = make(Spaces, 0)
	// 获取社区内internal和private可以看到的
	err = db.Where(" dtime = ? and space_group_guid in (?) and visibility = ?", 0, groupGuids, int32(commonv1.CMN_VISBL_INTERNAL.Number())).Order("`sort` asc").Find(&spaceList).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceListByUser fail, err: %w", err)
	}

	spaceSecretList := make(Spaces, 0)
	// 用户自己创建的私有的可以看到
	err = db.Where("created_by = ?  and dtime = ? and space_group_guid in (?) and visibility = ?", uid, 0, groupGuids, commonv1.CMN_VISBL_SECRET.Number()).Order("`sort` asc").Find(&spaceSecretList).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceListByUser fail2, err: %w", err)
	}
	spaceList = append(spaceList, spaceSecretList...)
	return
}

// SpacePublicList 公开访问的社区
func SpacePublicList(db *gorm.DB) (spaceList Spaces, err error) {
	spaceList = make(Spaces, 0)
	// 获取社区内internal可以看到的
	err = db.Where("dtime = ? and visibility = ?", 0, int32(commonv1.CMN_VISBL_INTERNAL.Number())).Order("`sort` asc").Find(&spaceList).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceListByUser fail, err: %w", err)
	}
	return
}

// SpacePublicListByGuids 公开访问的社区
func SpacePublicListByGuids(db *gorm.DB, guids []string) (spaceList Spaces, err error) {
	spaceList = make(Spaces, 0)
	// 获取社区内internal可以看到的
	err = db.Where("dtime = ? and visibility = ? and guid in (?)",
		0, int32(commonv1.CMN_VISBL_INTERNAL.Number()), guids).Order("`sort` asc").Find(&spaceList).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceListByUser fail, err: %w", err)
	}
	return
}

// ListSpaceGuidsByUser 根据用户的权限，返回不同的空间列表
//func ListSpaceGuidsByUser(db *gorm.DB, uid int64) (spaceList Spaces, err error) {
//	spaceList = make(Spaces, 0)
//	// 获取社区内internal和private可以看到的
//	err = db.Select("guid").Where("dtime = ? and visibility in (?)", 0, []int32{int32(commonv1.CMN_VISBL_INTERNAL.Number())}).Order("`sort` asc").Find(&spaceList).Error
//	if err != nil {
//		return nil, fmt.Errorf("SpaceListByUser fail, err: %w", err)
//	}
//
//	spaceSecretList := make(Spaces, 0)
//	// 用户自己创建的私有的可以看到
//	err = db.Select("guid").Where("created_by = ? and dtime = ? and visibility = ?", uid, 0, commonv1.CMN_VISBL_SECRET.Number()).Order("`sort` asc").Find(&spaceSecretList).Error
//	if err != nil {
//		return nil, fmt.Errorf("SpaceListByUser fail2, err: %w", err)
//	}
//	spaceList = append(spaceList, spaceSecretList...)
//	return
//}

// SpaceListBySpaceGroupGuid 返回一个分组下的空间列表
func SpaceListBySpaceGroupGuid(db *gorm.DB, groupGuid string) (spaceList Spaces, err error) {
	spaceList = make(Spaces, 0)
	// 获取社区内internal和private可以看到的
	err = db.Where("dtime = ? and space_group_guid = ?", 0, groupGuid).Find(&spaceList).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceListByUser fail, err: %w", err)
	}
	return
}

// SpaceGroupCreate 创建一条记录
func SpaceGroupCreate(db *gorm.DB, data *SpaceGroup) (err error) {
	if err = db.Create(data).Error; err != nil {
		return fmt.Errorf("SpaceGroupCreate fail, err: %w", err)
	}
	return
}

// SpaceGroupUpdateGuid 根据主键更新一条记录
func SpaceGroupUpdateGuid(db *gorm.DB, id int64, guid string) (err error) {
	if err = db.Model(SpaceGroup{}).Where("id = ?", id).Update("guid", guid).Error; err != nil {
		return fmt.Errorf("SpaceGroupUpdateGuid fail, err: %w", err)
	}
	return
}

// GetSpaceGroupInfo 获取信息
func GetSpaceGroupInfo(db *gorm.DB, field string, guid string) (info SpaceGroup, err error) {
	if err = db.Select(field).Where("guid = ?", guid).Find(&info).Error; err != nil {
		err = fmt.Errorf("GetSpaceGroupInfo fail, err: %w", err)
		return
	}
	return
}

// SpaceGroupGetId 创建一条记录
func SpaceGroupGetId(db *gorm.DB, guid string) (id int64, err error) {
	var info SpaceGroup
	err = db.Select("id").Where("guid = ? ", guid).Find(&info).Error
	if err != nil {
		return 0, fmt.Errorf("space group exist fail, err: %w", err)
	}
	if info.Id == 0 {
		return 0, fmt.Errorf("space create get space group fail")
	}
	id = info.Id

	return
}

// SpaceCreate 创建一条记录
func SpaceCreate(db *gorm.DB, data *Space) (err error) {
	if err = db.Create(data).Error; err != nil {
		return fmt.Errorf("SpaceCreate failed,err: %w", err)
	}
	return
}

// SpaceUpdateGuid 根据主键更新一条记录
func SpaceUpdateGuid(db *gorm.DB, id int64, guid string) (err error) {
	if err = db.Model(Space{}).Where("id = ?", id).Update("guid", guid).Error; err != nil {
		return fmt.Errorf("SpaceUpdateGuid fail, err: %w", err)
	}
	return
}

// GetSpaceInfo 获取信息
// TODO 改造成eerrors，减少代码
func GetSpaceInfo(db *gorm.DB, field string, guid string) (info Space, err error) {
	if err = db.Select(field).Where("guid = ?", guid).Find(&info).Error; err != nil {
		err = fmt.Errorf("GetSpaceSetInfo fail, err: %w", err)
		return
	}
	return
}

// GetSpaceInfoByGuid 获取信息
func GetSpaceInfoByGuid(db *gorm.DB, field string, guid string) (info Space, err error) {
	if err = db.Select(field).Where("guid = ?", guid).Find(&info).Error; err != nil {
		err = fmt.Errorf("GetSpaceSetInfo fail, err: %w", err)
		return
	}
	return
}

// SpaceGetId 创建一条记录
func SpaceGetId(db *gorm.DB, guid string) (id int64, err error) {
	var info Space
	err = db.Select("id").Where("guid = ?", guid).Find(&info).Error
	if err != nil {
		return 0, fmt.Errorf("space exist fail, err: %w", err)
	}
	id = info.Id
	return
}

// SpaceGroupSortGetInfoByIn 创建一条记录
func SpaceGroupSortGetInfoByIn(db *gorm.DB, guids []string) (list SpaceGroups, err error) {
	err = db.Select("guid,sort").Where("guid in (?)", guids).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceSortGetInfoByIn fail, err: %w", err)
	}
	return
}

// SpaceGetSpaceGroupGuid 创建一条记录
func SpaceGetSpaceGroupGuid(db *gorm.DB, cmtGuid string, guid string) (spaceGroupGuid string, err error) {
	var info Space
	err = db.Select("space_group_guid").Where("guid = ? and cmt_guid = ?", guid, cmtGuid).Find(&info).Error
	if err != nil {
		return "", fmt.Errorf("SpaceSortGetInfoByIn fail, err: %w", err)
	}
	spaceGroupGuid = info.SpaceGroupGuid
	return
}

// SpaceSortGetInfoByIn 创建一条记录
func SpaceSortGetInfoByIn(db *gorm.DB, guids []string) (list Spaces, err error) {
	err = db.Select("guid,space_group_guid,sort").Where("guid in (?) ", guids).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("SpaceSortGetInfoByIn fail, err: %w", err)
	}
	return
}

// SpaceGetSpaceType 查看space的频道
func SpaceGetSpaceType(db *gorm.DB, cmtGuid string, guid string) (spaceType commonv1.CMN_APP, err error) {
	var info Space
	err = db.Select("`type`").Where("guid = ? and cmt_guid = ?", guid, cmtGuid).Find(&info).Error
	if err != nil {
		return commonv1.CMN_APP_INVALID, fmt.Errorf("SpaceGetSpaceType fail, err: %w", err)
	}
	spaceType = info.Type
	return
}
