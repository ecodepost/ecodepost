package pmsaction

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	pmsv1 "ecodepost/pb/pms/v1"
	"ecodepost/resource-svc/pkg/model/mysql"

	"gorm.io/gorm"
)

// ActionOption 用户根据不同类型，可以看到可配置的权限列表
// 每种类型，都会有不同的权限属性，权限属性在 pms action表中
// 举例，超级管理员  可以控制什么
// 那么pms action type 是 1
// 然后将对应的action，id，放入到该表中
// todo 权限太尼玛复杂了，space种类不同，权限细分也不一样，在加一个子类别，碰到space类型，需要在查子类别
type ActionOption struct {
	Action     commonv1.PMS_ACT
	ActionType commonv1.PMS_ACT_TYPE
	SpaceType  commonv1.CMN_APP
}

type ActionOptions []ActionOption

func (list ActionOptions) ToActions() PmsActions {
	output := make(PmsActions, 0)
	for _, value := range list {
		pmsAction, err := GetAction(value.Action)
		if err != nil {
			continue
		}
		output = append(output, pmsAction)
	}
	return output
}

var (
	ActionOptionList = []ActionOption{
		{
			Action:     commonv1.PMS_COMMUNITY_DATA,
			ActionType: commonv1.PMS_ACT_TYPE_ROLE,
			SpaceType:  commonv1.CMN_APP_INVALID,
		},
		{
			Action:     commonv1.PMS_SPACE_GROUP_CREATE,
			ActionType: commonv1.PMS_ACT_TYPE_ROLE,
			SpaceType:  commonv1.CMN_APP_INVALID,
		},
		{
			Action:     commonv1.PMS_SPACE_GROUP_UPDATE,
			ActionType: commonv1.PMS_ACT_TYPE_ROLE,
			SpaceType:  commonv1.CMN_APP_INVALID,
		},
		{
			Action:     commonv1.PMS_SPACE_GROUP_DELETE,
			ActionType: commonv1.PMS_ACT_TYPE_ROLE,
			SpaceType:  commonv1.CMN_APP_INVALID,
		},
		{
			Action:     commonv1.PMS_SPACE_CREATE,
			ActionType: commonv1.PMS_ACT_TYPE_ROLE,
			SpaceType:  commonv1.CMN_APP_INVALID,
		},
		{
			Action:     commonv1.PMS_SPACE_UPDATE,
			ActionType: commonv1.PMS_ACT_TYPE_ROLE,
			SpaceType:  commonv1.CMN_APP_INVALID,
		},
		{
			Action:     commonv1.PMS_SPACE_DELETE,
			ActionType: commonv1.PMS_ACT_TYPE_ROLE,
			SpaceType:  commonv1.CMN_APP_INVALID,
		},
		{
			Action:     commonv1.PMS_FILE_CREATE,
			ActionType: commonv1.PMS_ACT_TYPE_SPACE,
			SpaceType:  commonv1.CMN_APP_ARTICLE,
		},
		{
			Action:     commonv1.PMS_FILE_DELETE,
			ActionType: commonv1.PMS_ACT_TYPE_SPACE,
			SpaceType:  commonv1.CMN_APP_ARTICLE,
		},
		{
			Action:     commonv1.PMS_FILE_DELETE_COMMENT,
			ActionType: commonv1.PMS_ACT_TYPE_SPACE,
			SpaceType:  commonv1.CMN_APP_ARTICLE,
		},
		{
			Action:     commonv1.PMS_FILE_SET_RECOMMEND,
			ActionType: commonv1.PMS_ACT_TYPE_SPACE,
			SpaceType:  commonv1.CMN_APP_ARTICLE,
		},
		{
			Action:     commonv1.PMS_FILE_SET_COMMENT,
			ActionType: commonv1.PMS_ACT_TYPE_SPACE,
			SpaceType:  commonv1.CMN_APP_ARTICLE,
		},
		{
			Action:     commonv1.PMS_FILE_SET_SITE_TOP,
			ActionType: commonv1.PMS_ACT_TYPE_SPACE,
			SpaceType:  commonv1.CMN_APP_ARTICLE,
		},
	}
)

// ListActionOptionByType 根据type类型获得，action 列表
// 这个适用于space group、role type。不适用于space
func ListActionOptionByType(actionType commonv1.PMS_ACT_TYPE) (list PmsActions, err error) {
	if actionType != commonv1.PMS_ACT_TYPE_ROLE && actionType != commonv1.PMS_ACT_TYPE_SPACE_GROUP {
		return nil, fmt.Errorf("dont support action type: " + actionType.String())
	}
	optionList := make(ActionOptions, 0)
	for _, value := range ActionOptionList {
		if value.ActionType == actionType {
			optionList = append(optionList, value)
		}
	}
	return optionList.ToActions(), nil
}

// ListActionOptionBySpaceType  根据type类型获得，action 列表
// 这个只适用于space，因为只有space的时候，才存在space type
func ListActionOptionBySpaceType(spaceType commonv1.CMN_APP) (list PmsActions, err error) {
	optionList := make(ActionOptions, 0)
	for _, value := range ActionOptionList {
		if value.ActionType == commonv1.PMS_ACT_TYPE_SPACE && value.SpaceType == spaceType {
			optionList = append(optionList, value)
		}
	}
	return optionList.ToActions(), nil
}

// ListActionByRoleId 根据role id，获取角色的权限列表enable情况
func ListActionByRoleId(rolePolicies mysql.PmsPolicies) (list []*commonv1.PmsItem, err error) {
	pmsActions, err := ListActionOptionByType(commonv1.PMS_ACT_TYPE_ROLE)
	if err != nil {
		err = fmt.Errorf("ListAction fail3, err: %w", err)
		return
	}

	policyMap := make(map[string]mysql.PmsPolicy)
	for _, value := range rolePolicies {
		policyMap[value.ActionName] = value
	}
	list = pmsActions.ToPb(policyMap)
	return
}

// ListSpaceGroupActionByRoleId 根据role id，获取space group的权限列表enable情况
func ListSpaceGroupActionByRoleId(db *gorm.DB, roleId int64, spaceGroupPolicies mysql.PmsPolicies) (list []*pmsv1.SpaceGroupPmsItem, err error) {
	spaceGroupActions, err := ListActionOptionByType(commonv1.PMS_ACT_TYPE_SPACE_GROUP)
	if err != nil {
		err = fmt.Errorf("ListAction fail3, err: %w", err)
		return
	}

	spaceRoleList, err := mysql.PmsRoleSpaceList(db, roleId, commonv1.CMN_GUID_SPACE_GROUP)
	if err != nil {
		err = fmt.Errorf("ListAction fail3, err: %w", err)
		return
	}
	guids := spaceRoleList.ToGuids()
	// guids := make([]string, 0)
	list = make([]*pmsv1.SpaceGroupPmsItem, 0)
	spaceGroupPoliciesMap := make(map[string]map[string]mysql.PmsPolicy)
	for _, spaceGroupPmsPolicy := range spaceGroupPolicies {
		info, flag := spaceGroupPoliciesMap[spaceGroupPmsPolicy.ResourceGuid]
		if !flag {
			info = make(map[string]mysql.PmsPolicy)
		}
		info[spaceGroupPmsPolicy.ActionName] = spaceGroupPmsPolicy
		spaceGroupPoliciesMap[spaceGroupPmsPolicy.ResourceGuid] = info
	}

	spaceGroupInfoList, err := mysql.SpaceGroupGetInfoByInGuids(db, "guid,name", guids)
	if err != nil {
		err = fmt.Errorf("ListSpaceGroupActionByRoleId fail, err: %w", err)
		return
	}

	for _, value := range spaceGroupInfoList {
		info, flag := spaceGroupPoliciesMap[value.Guid]
		if flag {
			list = append(list, &pmsv1.SpaceGroupPmsItem{
				Guid: value.Guid,
				Name: value.Name,
				List: spaceGroupActions.ToPb(info),
			})
		} else {
			list = append(list, &pmsv1.SpaceGroupPmsItem{
				Guid: value.Guid,
				Name: value.Name,
				List: spaceGroupActions.ToPb(map[string]mysql.PmsPolicy{}),
			})
		}
	}

	return
}

// ListSpaceActionByRoleId 根据role id，获取space的权限列表enable情况
func ListSpaceActionByRoleId(db *gorm.DB, roleId int64, spacePolicies mysql.PmsPolicies) (list []*pmsv1.SpacePmsItem, err error) {
	spaceRoleList, err := mysql.PmsRoleSpaceList(db, roleId, commonv1.CMN_GUID_SPACE)
	if err != nil {
		err = fmt.Errorf("ListAction fail3, err: %w", err)
		return
	}
	guids := spaceRoleList.ToGuids()
	list = make([]*pmsv1.SpacePmsItem, 0)
	spacePoliciesMap := make(map[string]map[string]mysql.PmsPolicy)
	for _, spacePmsPolicy := range spacePolicies {
		info, flag := spacePoliciesMap[spacePmsPolicy.ResourceGuid]
		if !flag {
			info = make(map[string]mysql.PmsPolicy)
		}
		info[spacePmsPolicy.ActionName] = spacePmsPolicy
		spacePoliciesMap[spacePmsPolicy.ResourceGuid] = info
	}

	spaceInfoList, err := mysql.SpaceGetInfoByInGuids(db, "guid,name,type", guids)
	if err != nil {
		err = fmt.Errorf("ListSpaceGroupActionByRoleId fail, err: %w", err)
		return
	}
	spaceInfoMap := spaceInfoList.ToMap()

	for _, value := range spaceInfoList {
		spaceActions, err := ListActionOptionBySpaceType(spaceInfoMap[value.Guid].Type)
		if err != nil {
			return nil, err
		}

		info, flag := spacePoliciesMap[value.Guid]
		if flag {
			list = append(list, &pmsv1.SpacePmsItem{
				Guid: value.Guid,
				Name: value.Name,
				List: spaceActions.ToPb(info),
			})
		} else {
			list = append(list, &pmsv1.SpacePmsItem{
				Guid: value.Guid,
				Name: value.Name,
				List: spaceActions.ToPb(map[string]mysql.PmsPolicy{}),
			})
		}
	}
	return
}
