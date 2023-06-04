package mysql

type NotifySign struct {
	BaseModel
	ChannelId   int    `gorm:"column:channel_id;default:0;NOT NULL" json:"channel_id"` // 通道id
	ThirdSignId string `gorm:"column:third_sign_id;NOT NULL" json:"third_sign_id"`     // 第三方签名id
	Content     string `gorm:"column:content;NOT NULL" json:"content"`                 // 内容
}

func (t NotifySign) TableName() string {
	return "notify_sign"
}
