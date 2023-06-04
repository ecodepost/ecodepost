package enotify

import (
	"github.com/google/uuid"
)

type parentPlugin struct {
	channelId int // reg channelId
}

func (d *parentPlugin) ChannelId() int {
	return d.channelId
}

// 按需使用 如第三方不产生则自己产生 MsgId
func (d *parentPlugin) genMsgId() string {
	return uuid.New().String()
}
