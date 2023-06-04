package spaceoption

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
)

// SpaceOption 用户根据不同类型，可以看到可配置的权限列表
type SpaceOption struct {
	SpaceType commonv1.CMN_APP
	Option    commonv1.SPC_OPTION
	Type      commonv1.SPC_OPTION_TYPE
	Name      string
}

var (
	SpaceTypeList = []commonv1.CMN_APP{
		commonv1.CMN_APP_ARTICLE,
		commonv1.CMN_APP_QA,
		commonv1.CMN_APP_COLUMN,
		commonv1.CMN_APP_LINK,
	}

	List = []SpaceOption{
		{
			SpaceType: commonv1.CMN_APP_ARTICLE,
			Option:    commonv1.SPC_OPTION_FILE_DEFAULT_SORT,
			Type:      commonv1.SPC_OPTION_TYPE_SELECT,
			Name:      "排序",
		},
		{
			SpaceType: commonv1.CMN_APP_ARTICLE,
			Option:    commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_FILE,
			Type:      commonv1.SPC_OPTION_TYPE_SWITCH,
			Name:      "允许成员发表帖子",
		},
		{
			SpaceType: commonv1.CMN_APP_ARTICLE,
			Option:    commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_COMMENT,
			Type:      commonv1.SPC_OPTION_TYPE_SWITCH,
			Name:      "允许成员评论帖子",
		},
		{
			SpaceType: commonv1.CMN_APP_ARTICLE,
			Option:    commonv1.SPC_OPTION_SITE_TOP_FILE_SHOW_ALL,
			Type:      commonv1.SPC_OPTION_TYPE_SWITCH,
			Name:      "允许空间置顶帖子内容全部展开",
		},
		{
			SpaceType: commonv1.CMN_APP_QA,
			Option:    commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_FILE,
			Type:      commonv1.SPC_OPTION_TYPE_SELECT,
			Name:      "排序",
		},
		{
			SpaceType: commonv1.CMN_APP_QA,
			Option:    commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_COMMENT,
			Type:      commonv1.SPC_OPTION_TYPE_SWITCH,
			Name:      "允许成员发表问题",
		},
		{
			SpaceType: commonv1.CMN_APP_QA,
			Option:    commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_FILE,
			Type:      commonv1.SPC_OPTION_TYPE_SWITCH,
			Name:      "允许成员评论问题",
		},
		{
			SpaceType: commonv1.CMN_APP_COLUMN,
			Option:    commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_COMMENT,
			Type:      commonv1.SPC_OPTION_TYPE_SWITCH,
			Name:      "允许成员评论",
		},
	}

	spaceTypeMap = make(map[commonv1.CMN_APP][]SpaceOption)
)

func init() {
	for _, value := range SpaceTypeList {
		spaceTypeMap[value] = make([]SpaceOption, 0)
	}
	for _, value := range List {
		list, flag := spaceTypeMap[value.SpaceType]
		if !flag {
			continue
		}
		list = append(list, value)
		spaceTypeMap[value.SpaceType] = list
	}
}

func GetListBySpaceType(spcType commonv1.CMN_APP) ([]SpaceOption, error) {
	list, flag := spaceTypeMap[spcType]
	if !flag {
		return nil, fmt.Errorf("not exist type: " + spcType.String())
	}
	return list, nil
}
