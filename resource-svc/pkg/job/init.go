package job

import (
	"github.com/gotomicro/ego/task/ejob"
)

// RunInitData 初始化数据
// 所有权限               FULL_ACCESS
// 社区概况               COMMUNITY_INFO
// 管理员-更新日志         COMMUNITY_LOGGER
// 管理员-超级管理员       COMMUNITY_ADMIN
// 管理员-分组/空间管理员   SPACE_MANAGER
// 数据报表               COMMUNITY_DATA
// 邀请用户加入社区        COMMUNITY_INVITE
// 创建空间分组            SPACE_CREATE
// 空间设置               SPACE_SET
func RunInitData(ctx ejob.Context) error {
	// tx := invoker.Db.WithContext(ctx.Ctx)
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "FULL_ACCESS",
	//	Desc:      "所有权限",
	//	EditionId: commonv1.CMT_EID_PROFESSIONAL,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "COMMUNITY_DATA",
	//	Title:     "查看社区数据",
	//	Desc:      "允许成员查看社区数据分析",
	//	EditionId: commonv1.CMT_EID_PROFESSIONAL,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_GROUP_CREATE",
	//	Title:     "创建空间分组",
	//	Desc:      "允许成员创建新的空间分组",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_CREATE",
	//	Title:     "创建空间",
	//	Desc:      "允许成员创建新的空间",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_ARTICLE_CREATE",
	//	Title:     "创建帖子",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_ARTICLE_CREATE_COMMENT",
	//	Title:     "发表评论",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_ARTICLE_CREATE_COMMENT",
	//	Title:     "发表评论",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })

	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_SET",
	//	Title:     "设置空间",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_ACTIVITY_PUBLISH",
	//	Title:     "发布活动",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_IM_AT_ALL_MEMBER",
	//	Title:     "@全员",
	//	EditionId: commonv1.CMT_EID_PROFESSIONAL,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_IM_SEND_MSG",
	//	Title:     "发送消息",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_QUESTION_PUBLISH",
	//	Title:     "发布问题",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "COMMUNITY_LOGGER",
	//	Title:     "允许查看管理员日志",
	//	EditionId: commonv1.CMT_EID_PROFESSIONAL,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "SPACE_MANAGER",
	//	Title:     "设置为空间分组管理员",
	//	EditionId: commonv1.CMT_EID_PROFESSIONAL,
	// })
	// mysql.PmsActionCreate(tx, &mysql.PmsAction{
	//	Name:      "COMMUNITY_INVITE",
	//	Desc:      "邀请用户加入社区",
	//	EditionId: commonv1.CMT_EID_BASIC,
	// })

	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "FULL_ACCESS",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_SUPER_ADMIN,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "COMMUNITY_DATA",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_SUPER_ADMIN,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_GROUP_CREATE",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_SUPER_ADMIN,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_CREATE",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_SUPER_ADMIN,
	// })
	//
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "COMMUNITY_DATA",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_GROUP_CREATE",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_CREATE",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "COMMUNITY_INVITE",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_ARTICLE_CREATE",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_ARTICLE,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_ARTICLE_CREATE",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_ARTICLE,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_ARTICLE_CREATE_COMMENT",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_ARTICLE,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_SET",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_ARTICLE,
	// })
	//
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_ACTIVITY_PUBLISH",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_ACTIVITY,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_SET",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_ACTIVITY,
	// })
	//
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_IM_AT_ALL_MEMBER",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_IM,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_IM_SEND_MSG",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_IM,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_SET",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_IM,
	// })
	//
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_QUESTION_PUBLISH",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_QUESTION,
	// })
	// mysql.PmsActionOptionCreate(tx, &mysql.PmsActionOption{
	//	ActionName: "SPACE_SET",
	//	Type:       commonv1.PmsActionType_PMS_ACTION_TYPE_CUSTOM_SPACE,
	//	SpaceType:  commonv1.SPC_TYPE_QUESTION,
	// })

	return nil
}
