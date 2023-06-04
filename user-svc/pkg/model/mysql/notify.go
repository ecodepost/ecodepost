package mysql

import (
	commonv1 "ecodepost/pb/common/v1"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Notify 消息和接收人关联表
type Notify struct {
	BaseModel
	ClientID       string                  `gorm:"client_id; not null; size:60; comment:'客户端id'" json:"client_id"`         // 客户端id
	MsgTmplId      int                     `gorm:"msg_tmpl_id; not null; comment:'消息模板id'" json:"msg_tmpl_id"`             // 消息模板id
	Email          string                  `gorm:"email; not null; size:60; comment:'email'" json:"email"`                 // 邮箱
	Phone          string                  `gorm:"phone; not null; size:60; comment:'phone'" json:"phone"`                 // 手机号
	Uid            int64                   `gorm:"uid; not null; size:60; comment:'uid'" json:"uid"`                       // 用户Uid
	Channel        commonv1.NOTIFY_CHANNEL `gorm:"channel; not null; comment:'消息类型'" json:"channel"`                       // 通道
	Status         commonv1.NOTIFY_STATUS  `gorm:"status; not null; comment:'消息状态'" json:"status"`                         // 消息状态
	MsgId          string                  `gorm:"msg_id; not null; size:128; comment:'消息结果ID'" json:"msg_id"`             // 消息id
	ErrorLog       string                  `gorm:"error_log; type:text; comment:'失败log'" json:"error_log"`                 // 错误原因
	RetryCount     int                     `gorm:"retry_count; not null; comment:'重试次数'" json:"retry_count"`               // 重试次数
	Title          string                  `gorm:"title; type:text; comment:'消息标题'" json:"title"`                          // 标题
	ExtraId        string                  `gorm:"extra_id; not null; size:32; comment:'业务层扩展id'" json:"extra_id"`         // 业务id
	ExtraContent   string                  `gorm:"extra_content; type:text; comment:'业务方自定义数据'" json:"extra_content"`      // 业务内容
	FinalContent   string                  `gorm:"final_content; type:text; comment:'消息体最终内容（按需存储）'" json:"final_content"` // 最后内容
	ErrorType      int                     `gorm:"error_type; not null; comment:'错误类型'" json:"error_type"`                 // 错误类型
	CCEmail        string                  `gorm:"cc_email; type:text; comment:'邮件抄送人，分号隔离'" json:"cc_email"`              //  抄送
	Sender         string                  `gorm:"not null" json:"sender"`
	Vars           datatypes.JSON          `gorm:"vars; comment:'变量数据'" json:"vars"`                                                     // 模板变量
	AttachEmailUrl string                  `gorm:"attach_email_url; type:text; comment:'邮件附件FDS地址 格式 ['', '']'" json:"attach_email_url"` // 额外参数, 站内信, 需配合模板使用, 短信, 需配合模板使用,额外参数, 邮件, 需配合模板使用
	TplData        datatypes.JSON          `gorm:"type:longtext; comment:模板数据" json:"tpl_data"`
}

type Notifies []*Notify

// TableName 设置表明
func (t Notify) TableName() string {
	return "notify"
}

// NotifyCreate 创建一条记录
func NotifyCreate(db *gorm.DB, data *Notify) (err error) {
	if err = db.Create(data).Error; err != nil {
		elog.Error("create notify error", zap.Error(err))
		return
	}
	return
}

// NotifyInitList 创建一条记录
func NotifyInitList(db *gorm.DB) (list Notifies, err error) {
	if err = db.Where("status = ?", commonv1.NOTIFY_STATUS_INIT).Find(&list).Error; err != nil {
		elog.Error("create NotifyInitList error", zap.Error(err))
		return
	}
	return
}
