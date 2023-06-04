package filecontent

import (
	"encoding/json"
	"fmt"
)

func GetPointerContent(content json.RawMessage) (pContent *string, err error) {
	if content == nil {
		return
	}
	strByte, err := content.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("GetPointerContent fail, err: %w", err)
		return
	}
	strstr := string(strByte)
	pContent = &strstr
	return
}

type SlateItem struct {
	Text string `json:"text"`
}

type SlateParent struct {
	Type     string      `json:"type"`
	Children []SlateItem `json:"children"`
}

// GetContentStr 获取内容
// 初始化数据，必须得是这样得数据，否则slate不识别
// strByte = []byte(`[{type: "p", children: [{text: ""}]}]`)
func GetContentStr(content json.RawMessage) (contentStr string, err error) {
	if content == nil {
		info := []SlateParent{{Type: "p", Children: []SlateItem{
			{""},
		}}}
		contentBytes, err := json.Marshal(info)
		if err != nil {
			return "", fmt.Errorf("get content str fail, err1: %w", err)
		}
		return string(contentBytes), nil

	}
	contentBytes, err := content.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("get content str fail, err2: %w", err)
	}
	return string(contentBytes), nil
}
