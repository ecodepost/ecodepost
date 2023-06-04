package enotify

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"

	commonv1 "ecodepost/pb/common/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	"github.com/ego-component/egorm"
	"github.com/go-gomail/gomail"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
)

type emailConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Subject  string `toml:"subject"`
	From     string `toml:"from"`
	FromName string `toml:"fromName"`
	ToName   string `toml:"toName"`
}

type emailPlugin struct {
	config emailConfig
	state  atomic.Value
	parentPlugin
	gomail.SendCloser
}

func init() {
	Register(commonv1.NOTIFY_CHANNEL_EMAIL_COMMON, &emailPlugin{})
}

func (e *emailPlugin) Init(db *egorm.Component) (err error) {
	err = econf.UnmarshalKey("notify.email.common", &e.config)
	if err != nil {
		elog.Error("get email config error", elog.FieldErr(err))
		return err
	}
	var sender gomail.SendCloser
	dialer := gomail.NewDialer(e.config.Host, e.config.Port, e.config.Username, e.config.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if sender, err = dialer.Dial(); err != nil {
		elog.Error("dial smtp server failed, continue", elog.FieldValueAny(e.config), elog.FieldErr(err))
	}
	e.SendCloser = sender
	return err
}

func (e *emailPlugin) Enable() bool {
	return econf.Get("notify.email.common") != nil
}

func (e *emailPlugin) Destroy() (err error) {
	return
}

func (e *emailPlugin) Send(req *SendRequest) (resp *SendResponse, err error) {
	resp = &SendResponse{Code: 0}
	err = e.doSend(req)
	resp.FinalContent = req.FinalContent
	resp.ThirdMsgId = e.genMsgId()
	if err != nil {
		resp.Code = 500
		resp.Reason = err.Error()
		return nil, fmt.Errorf("send fail, err: %w", err)
	}
	return
}

func (w *emailPlugin) doSend(req *SendRequest) (err error) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", w.config.From, w.config.FromName)

	jsonData := notifyv1.Email{}
	err = json.Unmarshal(req.Vars, &jsonData)
	if err != nil {
		elog.Error("invalid json unmarshal faile", zap.Any("req", req), elog.FieldErr(err))
		return fmt.Errorf("do send fail, err: %w", err)
	}

	toName := jsonData.ToName
	if strings.TrimSpace(toName) == "" {
		toName = w.config.ToName
	}
	subject := jsonData.Subject
	if strings.TrimSpace(subject) == "" {
		subject = w.config.Subject
	}
	m.SetAddressHeader("To", req.Receiver, toName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", req.FinalContent)
	if err = gomail.Send(w.SendCloser, m); err != nil {
		elog.Error("send emailPlugin error", elog.FieldErr(err))
	}
	return
}
