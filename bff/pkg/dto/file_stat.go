package dto

import (
	commonv1 "ecodepost/pb/common/v1"
)

type MyEmojiInfo struct {
	Guid string                `json:"guid"`
	List []*commonv1.EmojiInfo `json:"list"`
}

type MyCollectInfo struct {
	Guid      string `json:"guid"`
	IsCollect bool   `json:"isCollect"`
}
