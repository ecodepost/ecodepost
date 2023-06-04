package enotify

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"ecodepost/pb/common/v1"
	"ecodepost/pb/notify/v1"
	"ecodepost/user-svc/pkg/model/mysql"
	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// NotifyLetter 站内信
type letter struct {
	db *egorm.Component
	parentPlugin
	sync.Mutex
}

func init() {
	Register(commonv1.NOTIFY_CHANNEL_LETTER, &letter{})
}

func (l *letter) Enable() bool {
	return econf.Get("notify.letter") != nil
}

func (l *letter) Init(db *egorm.Component) error {
	l.db = db
	return nil
}

// Send 根据 vars recvType、 link 设置kafka消息体中响应字段
func (l *letter) Send(req *SendRequest) (resp *SendResponse, err error) {
	uid := cast.ToInt64(req.Receiver)
	if uid == 0 {
		elog.Error("invalid sendRequest", zap.Any("req", req))
		return &SendResponse{Code: 500, Reason: "uid can't be empty"}, fmt.Errorf("invalid sendrequest, uid can't be empty")
	}

	jsonData := notifyv1.Letter{}
	err = json.Unmarshal(req.Vars, &jsonData)
	if err != nil {
		elog.Error("invalid json unmarshal fail", zap.Any("req", req), elog.FieldErr(err))
		return &SendResponse{Code: 500, Reason: "invalid json unmarshal fail"}, fmt.Errorf("invalid json unmarshal fail, err: %w", err)
	}

	createReq := CreateNotificationReq{
		Type:        commonv1.NTF_TYPE(jsonData.Type),
		TargetId:    jsonData.TargetId,
		ReceiverUid: uid,
		SenderId:    req.Sender,
		Link:        jsonData.Link,
		Meta:        jsonData.Meta,
	}
	notification, e := l.CreateLetter(context.Background(), &createReq)
	if e != nil {
		return &SendResponse{Code: 500, Reason: fmt.Sprintf("Marshal text error %s", e.Error())}, fmt.Errorf("marshal fail, err: %w", err)
	}
	return &SendResponse{
		Code:         0,
		FinalContent: "",
		ThirdMsgId:   cast.ToString(notification.ID),
	}, nil
}

type CreateNotificationReq struct {
	Type        commonv1.NTF_TYPE `json:"type"`
	TargetId    string            `json:"targetId"`
	Link        string            `json:"link"`
	Meta        []byte            `json:"meta"`
	ReceiverUid int64             `json:"receiverUid"`
	SenderId    string            `json:"senderId"`
	Ctime       int64             `json:"ctime"`
}

// CreateLetter 创建通知
func (l *letter) CreateLetter(ctx context.Context, req *CreateNotificationReq) (*mysql.NotifyLetter, error) {
	if req.ReceiverUid == 0 {
		return nil, fmt.Errorf("uid can't be empty")
	}
	// 创建 NotifyLetter
	notification := &mysql.NotifyLetter{
		Type:     req.Type,
		TargetId: req.TargetId,
		Link:     req.Link,
		Uid:      req.ReceiverUid,
		SenderId: req.SenderId,
	}
	if len(req.Meta) > 0 {
		notification.Meta = string(req.Meta)
	}
	// 入库
	if err := mysql.LetterCreate(l.db.WithContext(ctx), notification); err != nil {
		return nil, fmt.Errorf("mysql.LetterCreate fail, %w", err)
	}
	return notification, nil
}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func (l *letter) Destroy() (err error) {
	l.db = nil
	return err
}
