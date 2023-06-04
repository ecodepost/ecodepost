package service

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	loggerv1 "ecodepost/pb/logger/v1"
)

type logger struct {
}

// LoggerEvent 这个数据没多少，写死在代码里，不放数据库了
type LoggerEvent struct {
	Id     commonv1.LOG_EVENT
	CNName string // 显示在外的名称
	ENName string // 显示在外的名称
}

type LoggerGroup struct {
	Id     commonv1.LOG_GROUP
	CNName string // 显示在外的名称
	ENName string // 显示在外的名称
}

var (
	LoggerEventList = []LoggerEvent{
		{
			Id:     commonv1.LOG_EVENT_SPACE_CREATE,
			CNName: "创建空间",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_UPDATE,
			CNName: "更新空间",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_DELETE,
			CNName: "删除空间",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_CHANGE_SORT,
			CNName: "调整空间顺序",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_GROUP_CREATE,
			CNName: "创建空间分组",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_GROUP_CREATE,
			CNName: "创建空间分组",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_GROUP_UPDATE,
			CNName: "更新空间分组",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_GROUP_DELETE,
			CNName: "删除空间分组",
		},
		{
			Id:     commonv1.LOG_EVENT_SPACE_GROUP_CHANGE_SORT,
			CNName: "调整空间分组顺序",
		},
		{
			Id:     commonv1.LOG_EVENT_PMS_SUPER_ADMIN_CREATE_MEMBER,
			CNName: "添加超级管理员",
		},
		{
			Id:     commonv1.LOG_EVENT_PMS_SUPER_ADMIN_DELETE_MEMBER,
			CNName: "删除超级管理员",
		},
		{
			Id:     commonv1.LOG_EVENT_PMS_ROLE_CREATE,
			CNName: "创建权限分组",
		},
		{
			Id:     commonv1.LOG_EVENT_PMS_ROLE_UPDATE,
			CNName: "更新权限分组",
		},
		{
			Id:     commonv1.LOG_EVENT_PMS_ROLE_DELETE,
			CNName: "删除权限分组",
		},
		{
			Id:     commonv1.LOG_EVENT_PMS_ROLE_CREATE_MEMBER,
			CNName: "添加分组成员",
		},
		{
			Id:     commonv1.LOG_EVENT_PMS_ROLE_DELETE_MEMBER,
			CNName: "删除分组成员",
		},
	}
	LoggerGroupList = []LoggerGroup{
		{
			Id:     commonv1.LOG_GROUP_SPACE,
			CNName: "空间",
		},
		{
			Id:     commonv1.LOG_GROUP_SPACE_GROUP,
			CNName: "空间分组",
		},
		{
			Id:     commonv1.LOG_GROUP_PMS,
			CNName: "权限",
		},
	}
	LoggerEventCNPbList = make([]*loggerv1.EventInfo, 0)
	LoggerGroupCNPbList = make([]*loggerv1.GroupInfo, 0)
)

func init() {
	for _, value := range LoggerEventList {
		LoggerEventCNPbList = append(LoggerEventCNPbList, &loggerv1.EventInfo{
			Event: value.Id,
			Name:  value.CNName,
		})
	}
	for _, value := range LoggerGroupList {
		LoggerGroupCNPbList = append(LoggerGroupCNPbList, &loggerv1.GroupInfo{
			Group: value.Id,
			Name:  value.CNName,
		})
	}
}

func (*logger) GetEvent(eventId commonv1.LOG_EVENT) (LoggerEvent, error) {
	for _, value := range LoggerEventList {
		if value.Id == eventId {
			return value, nil
		}
	}
	return LoggerEvent{}, fmt.Errorf("not exist log event id")
}

func (*logger) GetGroup(groupId commonv1.LOG_GROUP) (LoggerGroup, error) {
	for _, value := range LoggerGroupList {
		if value.Id == groupId {
			return value, nil
		}
	}
	return LoggerGroup{}, fmt.Errorf("not exist log group id")
}
