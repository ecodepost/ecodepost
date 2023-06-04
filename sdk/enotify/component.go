package enotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"sync"

	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/user-svc/pkg/model/mysql"
	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const PackageName = "component.notify"

type (
	// Component letter 组件
	Component struct {
		sync.Mutex
		plugins map[commonv1.NOTIFY_CHANNEL]NotifyPlugin
		name    string
		db      *egorm.Component
	}

	// NotifyPlugin notify plugins
	NotifyPlugin interface {
		ChannelId() int                                        // 唯一 channelId register
		Destroy() error                                        // 销毁
		Init(*egorm.Component) error                           // 初始化
		Enable() bool                                          // 是否启用
		Send(req *SendRequest) (resp *SendResponse, err error) // send
	}

	// SendRequest 对业务 pb 不产生依赖
	SendRequest struct {
		Receiver     string                  // 接收者
		Sender       string                  // 发送者, 可以为空
		Vars         json.RawMessage         // 参数变量
		ExtraContent string                  // 业务方备注
		ExtraId      string                  // 业务方扩展id
		Tpl          *mysql.NotifyTpl        // 模板对象
		Ch           *mysql.NotifyTplChannel // tpl channel
		Sign         *mysql.NotifySign       // 签名
		TplData      []byte                  // 模板数据
		FinalContent string                  // 最终发送的内容
	}

	// SendResponse 对业务 pb 不产生依赖
	SendResponse struct {
		Code         int32  // code
		Reason       string // 详情
		FinalContent string // 最终发送内容
		ThirdMsgId   string // 三方消息id
	}
)

func newComponent(c *Container) *Component {
	return &Component{
		plugins: c.plugins,
		name:    c.name,
		db:      c.db,
	}
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Start() error {
	var err error

	fmt.Printf("registry--------------->"+"%+v\n", registry)
	for channelId, channel := range registry {
		if channel.Enable() {
			err = channel.Init(c.db)
			if err != nil {
				elog.Panic("start notify fail", elog.FieldErr(err))
			}
			c.plugins[channelId] = channel
			fmt.Printf("c.plugins--------------->"+"%+v\n", c.plugins)
		}
	}
	return err
}

func (c *Component) Stop() error {
	var err error
	if c.plugins != nil && len(c.plugins) > 0 {
		for k, v := range c.plugins {
			err = v.Destroy()
			if err != nil {
				elog.Error("destroy notify plugin error", elog.FieldName(cast.ToString(k)), elog.FieldErr(err))
			} else {
				elog.Info("destroy notify plugin", elog.FieldName(cast.ToString(k)))
			}
		}
	}
	return nil
}

func (c *Component) Send(msg mysql.Notify) (err error) {
	msgTpl := mysql.NotifyTpl{}

	// 查询模板
	err = c.db.Model(msgTpl).Where("`id` = ?", msg.MsgTmplId).First(&msgTpl).Error
	if err != nil {
		return err
	}

	// 查询发送渠道
	ch := mysql.NotifyTplChannel{}
	err = c.db.Model(ch).Where("channel_id = ? and tpl_id = ?", msg.Channel, msg.MsgTmplId).First(&ch).Error
	if err != nil {
		return err
	}

	// 如果是短信服务，需要查询SDK鉴权配置
	sign := mysql.NotifySign{}
	if msgTpl.ChannelType == 2 {
		err = c.db.Model(sign).Where("id = ?", ch.SignId).Find(&sign).Error
		if err != nil {
			return err
		}
	}

	// 获取Sender插件, 如果失败则更新MsgReceiver关联表状态
	msgRcvDB := c.db.Model(mysql.Notify{}).Where("`id` = ?", msg.ID)
	plugin, ok := c.plugins[msg.Channel]
	if !ok {
		elog.Error("plugin not exist", elog.FieldName(cast.ToString(msg.Channel)))
		return msgRcvDB.Updates(map[string]any{"error_log": "plugin not found", "status": uint8(commonv1.NOTIFY_STATUS_ERROR)}).Error
	}

	// 渲染消息模板, 如果失败则更新MsgReceiver关联表状态
	finalContent, err := c.renderTemplate(msgTpl.Content, msg.TplData)
	if err != nil {
		return msgRcvDB.Updates(map[string]any{"error_log": "render msgTpl failed. err=" + err.Error(), "status": uint8(commonv1.NOTIFY_STATUS_ERROR)}).Error
	}

	var receiver string

	switch msg.Channel {
	case commonv1.NOTIFY_CHANNEL_EMAIL_COMMON:
		receiver = msg.Email
	case commonv1.NOTIFY_CHANNEL_SMS_ALI:
		receiver = msg.Phone
	case commonv1.NOTIFY_CHANNEL_LETTER:
		receiver = cast.ToString(msg.Uid)
	default:
		return fmt.Errorf("unsupported channelType, %d", msg.Channel)
	}

	// 使用插件发送消息
	resp, err := plugin.Send(&SendRequest{
		Receiver:     receiver,
		Sender:       msg.Sender,
		ExtraContent: msg.ExtraContent,
		ExtraId:      msg.ExtraId,
		Tpl:          &msgTpl,
		Ch:           &ch,
		Sign:         &sign,
		TplData:      msg.TplData,
		FinalContent: finalContent,
		Vars:         json.RawMessage(msg.Vars),
	})
	if err != nil {
		ups := map[string]any{
			"final_content": resp.FinalContent,
			"msg_id":        resp.ThirdMsgId,
			"error_log":     resp.Reason,
			"status":        commonv1.NOTIFY_STATUS_ERROR,
		}
		return msgRcvDB.Updates(ups).Error
	}
	ups := map[string]any{
		"final_content": resp.FinalContent,
		"msg_id":        resp.ThirdMsgId,
		"error_log":     resp.Reason,
		"status":        commonv1.NOTIFY_STATUS_SUCCESS,
	}
	return msgRcvDB.Updates(ups).Error
}

func (c *Component) GetSenderPlugin(key commonv1.NOTIFY_CHANNEL) NotifyPlugin {
	return c.plugins[key]
}

// PluginStopFuncs 获取stop hook
func (c *Component) PluginStopFuncs() []func() error {
	arr := make([]func() error, 0, len(c.plugins))
	for _, p := range c.plugins {
		arr = append(arr, p.Destroy)
	}
	return arr
}

func (c *Component) renderTemplate(tplContent string, tplData []byte) (res string, err error) {
	globalVars := make(map[string]any)
	err = econf.UnmarshalKey("globalVariables", &globalVars)
	if err != nil {
		if errors.Is(err, econf.ErrInvalidKey) { // 如果没有配置, 忽略
			err = nil
		} else {
			return
		}
	}

	// 构造模板
	var tpl *template.Template
	tpl, err = template.New("tpl").Parse(tplContent)
	if err != nil {
		err = errors.Wrap(err, "parse template failed")
		return
	}

	buf := bytes.NewBuffer(nil)
	data := make(map[string]any)
	// 合并全局配置
	for k, v := range globalVars {
		data[k] = v
	}
	// 解析tplData
	if len(tplData) > 0 {
		err = json.Unmarshal(tplData, &data)
		if err != nil {
			err = errors.Wrap(err, "unmarshall tplData failed")
			return
		}
	}
	// 解析模板
	err = tpl.Execute(buf, data)
	if err != nil {
		err = errors.Wrap(err, "execute template failed")
		return
	}

	res = buf.String()
	return
}
