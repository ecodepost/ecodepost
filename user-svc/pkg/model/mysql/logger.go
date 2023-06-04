package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	"github.com/ego-component/egorm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Logger struct {
	Id              int64              `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Event           commonv1.LOG_EVENT `gorm:"not null;default:0;comment:版本号"`
	Group           commonv1.LOG_GROUP `gorm:"not null;default:0;comment:版本号"`
	TargetUid       int64              `gorm:"not null;default:0;comment:目标用户uid"`
	OperateUid      int64              `gorm:"not null;default:0;comment:操作用户uid"`
	SpaceGuid       string             `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT ''; index:idx_space_guid; comment:空间guid"`
	SpaceGroupGuid  string             `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT ''; index:idx_space_group_guid; comment:空间分组guid"`
	FileGuid        string             `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT ''; index:idx_file_guid; comment:文件guid"`
	Metadata        datatypes.JSON     `gorm:"type:json;comment:元信息"`
	Ctime           int64              `gorm:"bigint;autoCreateTime;comment:创建时间" json:"ctime"`
	EventName       string             `gorm:"-"`
	GroupName       string             `gorm:"-"`
	Message         string             `gorm:"-"`
	OperateNickname string             `gorm:"-"`
	OperateAvatar   string             `gorm:"-"`
	TargetNickname  string             `gorm:"-"`
	TargetAvatar    string             `gorm:"-"`
}

func (Logger) TableName() string {
	return "logger"
}

type Loggers []*Logger

func (list Loggers) ToAllUids() []int64 {
	output := make([]int64, 0)
	for _, value := range list {
		if value.OperateUid > 0 {
			output = append(output, value.OperateUid)
		}
		if value.TargetUid == 0 {
			output = append(output, value.TargetUid)
		}
	}
	return output
}

func (list Loggers) ToPb() []*loggerv1.LogInfo {
	output := make([]*loggerv1.LogInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPb())
	}
	return output
}

func (value Logger) ToPb() *loggerv1.LogInfo {
	return &loggerv1.LogInfo{
		Id:            value.Id,
		EventName:     value.EventName,
		GroupName:     value.GroupName,
		Message:       value.Message,
		OperateUid:    value.OperateUid,
		OperateName:   value.OperateNickname,
		OperateAvatar: value.OperateAvatar,
		TargetUid:     value.TargetUid,
		TargetName:    value.TargetNickname,
		TargetAvatar:  value.TargetAvatar,
		Ctime:         value.Ctime,
	}
}

func LoggerCreate(db *gorm.DB, logger *Logger) (err error) {
	err = db.Create(&logger).Error
	if err != nil {
		return fmt.Errorf("LoggerCreate fail, err: %w", err)
	}
	return nil
}

func LoggerCreateInBatches(db *gorm.DB, loggers []Logger) (err error) {
	err = db.CreateInBatches(loggers, len(loggers)).Error
	if err != nil {
		return fmt.Errorf("LoggerCreateInBatches fail, err: %w", err)
	}
	return nil
}

func LoggerListPage(db *gorm.DB, conds egorm.Conds, reqList *commonv1.Pagination) (respList Loggers, err error) {
	if reqList.PageSize == 0 {
		reqList.PageSize = 20
	}
	if reqList.CurrentPage == 0 {
		reqList.CurrentPage = 1
	}
	if reqList.Sort == "" {
		reqList.Sort = "id desc"
	}
	sql, binds := egorm.BuildQuery(conds)

	listDb := db.Model(Logger{}).Where(sql, binds...)
	err = listDb.Count(&reqList.Total).Error
	if err != nil {
		return nil, fmt.Errorf("audit list page fail,err: %w", err)
	}
	err = listDb.Order(reqList.Sort).Offset(int(reqList.CurrentPage-1) * int(reqList.PageSize)).Limit(int(reqList.PageSize)).Find(&respList).Error
	if err != nil {
		return nil, fmt.Errorf("audit list page fail2,err: %w", err)
	}
	return
}
