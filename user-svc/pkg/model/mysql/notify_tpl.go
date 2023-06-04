package mysql

type NotifyTpl struct {
	BaseModel
	Title           string `gorm:"column:title; type:text; comment:'消息标题'" json:"title"`                         // 消息标题
	Content         string `gorm:"column:content; type:longtext; comment:'消息内容'" json:"content"`                 // 消息内容
	TplType         int    `gorm:"column:tpl_type;default:0;NOT NULL; comment:'模板类型'" json:"tpl_type"`           // 模板类型
	ChannelType     int    `gorm:"column:channel_type;default:0;NOT NULL; comment:'通道类型'" json:"channel_type"`   // 通道类型
	ChannelTypeName string `gorm:"column:channel_type_name;NOT NULL; comment:'通道类型名称'" json:"channel_type_name"` // 通道类型名称
}
type MsgTpls []*NotifyTpl

// TableName 设置表明
func (t NotifyTpl) TableName() string {
	return "notify_tpl"
}
