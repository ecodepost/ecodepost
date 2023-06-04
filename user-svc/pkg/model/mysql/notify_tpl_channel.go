package mysql

type NotifyTplChannel struct {
	BaseModel
	ChannelId  int    `gorm:"column:channel_id;default:0;NOT NULL; comment:'mns 通道id'" json:"channel_id"`         // mns 通道id
	TplId      int    `gorm:"column:tpl_id;default:0;NOT NULL; comment:'mns模板id'" json:"tpl_id"`                  // mns模板id
	ThirdTplId string `gorm:"column:third_tpl_id;NOT NULL; comment:'第三方模板id'" json:"third_tpl_id"`                // 第三方模板id
	SignId     int    `gorm:"column:sign_id;default:0;NOT NULL; comment:'mns签名id 如果某通道签名在内容中则为0'" json:"sign_id"` // mns签名id 如果某通道签名在内容中则为0
	Weight     int    `gorm:"column:weight;default:0;NOT NULL; comment:'权重 正整数'" json:"weight"`                   // 权重 正整数
	ChStatus   int    `gorm:"column:ch_status;default:0;NOT NULL; comment:'审核状态'" json:"ch_status"`               // 审核状态
	Reason     string `gorm:"column:reason; type:text; comment:'原因(第三方接口返回原因)'" json:"reason"`                    // 原因(第三方接口返回原因)
}

func (t NotifyTplChannel) TableName() string {
	return "notify_tpl_channel"
}
