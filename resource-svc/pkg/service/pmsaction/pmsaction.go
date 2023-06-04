package pmsaction

import (
	"fmt"

	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

// PmsAction 展示所有的信息，不同版本对应不同的action
// 自定义role下面有单独的一些权限，例如查看社区数据，创建空间分组、创建空间等
// 自定义role里面有多个space group，多个space
// 每一个space group可以单独配置权限
// 每一个space也可以单独配置权限，这里每个space还有子类型，导致权限点还不一样，小心判断吧
type PmsAction struct {
	Action commonv1.PMS_ACT
	Title  string
	Desc   string
}

type PmsActions []PmsAction

func (list PmsActions) ToPb(policyMap map[string]mysql.PmsPolicy) []*commonv1.PmsItem {
	output := make([]*commonv1.PmsItem, 0)
	for _, value := range list {
		output = append(output, value.ToPb(policyMap))
	}
	return output
}

func (info PmsAction) ToPb(policyMap map[string]mysql.PmsPolicy) *commonv1.PmsItem {
	_, flag := policyMap[info.Action.String()]
	switchFlag := 0
	if flag {
		switchFlag = 1
	}
	return &commonv1.PmsItem{
		ActionName: info.Action.String(),
		Title:      info.Title,
		Desc:       info.Desc,
		Flag:       int32(switchFlag),
	}
}

var (
	ActionList = []PmsAction{
		{
			Action: commonv1.PMS_COMMUNITY_DATA,
			Title:  "查看社区数据",
			Desc:   "查看社区数据",
		},
		{
			Action: commonv1.PMS_COMMUNITY_LOGGER,
			Title:  "允许查看管理员日志",
			Desc:   "允许查看管理员日志",
		},
		{
			Action: commonv1.PMS_SPACE_GROUP_CREATE,
			Title:  "创建空间分组",
			Desc:   "创建空间分组",
		},
		{
			Action: commonv1.PMS_SPACE_GROUP_UPDATE,
			Title:  "设置空间分组",
			Desc:   "设置空间分组",
		},
		{
			Action: commonv1.PMS_SPACE_GROUP_DELETE,
			Title:  "删除空间分组",
			Desc:   "删除空间分组",
		},
		{
			Action: commonv1.PMS_SPACE_CREATE,
			Title:  "创建空间",
			Desc:   "创建空间",
		},
		{
			Action: commonv1.PMS_SPACE_UPDATE,
			Title:  "设置空间",
			Desc:   "设置空间",
		},
		{
			Action: commonv1.PMS_SPACE_DELETE,
			Title:  "删除空间",
			Desc:   "删除空间",
		},
		{
			Action: commonv1.PMS_FILE_CREATE,
			Title:  "创建帖子",
			Desc:   "创建帖子",
		},
		{
			Action: commonv1.PMS_FILE_SET_RECOMMEND,
			Title:  "设置推荐",
			Desc:   "设置推荐",
		},
		{
			Action: commonv1.PMS_FILE_SET_SITE_TOP,
			Title:  "设置置顶",
			Desc:   "设置置顶",
		},
		{
			Action: commonv1.PMS_FILE_SET_COMMENT,
			Title:  "设置评论开关",
			Desc:   "设置评论开关",
		},
		{
			Action: commonv1.PMS_FILE_DELETE,
			Title:  "删除帖子",
			Desc:   "删除帖子",
		},
		{
			Action: commonv1.PMS_FILE_CREATE,
			Title:  "发布问题",
			Desc:   "发布问题",
		},
		{
			Action: commonv1.PMS_FILE_DELETE,
			Title:  "删除问题",
			Desc:   "删除问题",
		},
		{
			Action: commonv1.PMS_FILE_CREATE,
			Title:  "发布活动",
			Desc:   "发布活动",
		},
		{
			Action: commonv1.PMS_FILE_DELETE,
			Title:  "删除活动",
			Desc:   "删除活动",
		},
	}
	actionMap = make(map[commonv1.PMS_ACT]PmsAction)
)

func init() {
	for _, value := range ActionList {
		actionMap[value.Action] = value
	}
}

func GetAction(action commonv1.PMS_ACT) (PmsAction, error) {
	value, flag := actionMap[action]
	if !flag {
		return PmsAction{}, fmt.Errorf("not exit action type")
	}
	return value, nil
}
