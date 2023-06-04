package enotify

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/user-svc/pkg/model/mysql"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AliConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
}

type aliPlugin struct {
	db     *egorm.Component
	config *AliConfig
	client *dysmsapi20170525.Client
	parentPlugin
	sync.Mutex
}

func init() {
	Register(commonv1.NOTIFY_CHANNEL_SMS_ALI, &aliPlugin{})
}

func (a *aliPlugin) Enable() bool {
	return econf.Get("notify.sms.ali") != nil
}
func (a *aliPlugin) Init(db *egorm.Component) error {
	var err error
	if econf.Get("notify.sms.ali") == nil {
		return nil
	}
	err = econf.UnmarshalKey("notify.sms.ali", &a.config)
	if err != nil {
		return err
	}
	config := &openapi.Config{
		// aliyun AccessKey ID
		AccessKeyId: &a.config.AccessKeyId,
		// aliyun AccessKey Secret
		AccessKeySecret: &a.config.AccessKeySecret,
		// endpoint
		Endpoint: &a.config.Endpoint,
	}
	a.client, err = dysmsapi20170525.NewClient(config)
	if err != nil {
		return err
	}
	return nil
}

func (a *aliPlugin) Destroy() error {
	a.Lock()
	defer a.Unlock()
	a.config = nil
	a.client = nil

	return nil
}

func (a *aliPlugin) Send(req *SendRequest) (resp *SendResponse, err error) {
	resp = &SendResponse{}
	var tplParam *string
	// TODO 使用varSms替换

	if req.Vars != nil {
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(req.Vars)
		if err != nil {
			resp.Code = int32(500)
			resp.Reason = err.Error()
			return nil, fmt.Errorf("json marshal fail, err: %w", err)
		}
		tplParam = tea.String(string(jsonBytes))
	}
	resp.FinalContent = fmt.Sprintf("【%s】%s", req.Sign.Content, req.Tpl.Content)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(req.Receiver),
		SignName:      tea.String(req.Sign.Content),
		TemplateCode:  tea.String(req.Ch.ThirdTplId),
		TemplateParam: tplParam,
		OutId:         tea.String(req.ExtraId),
	}

	result, err := a.client.SendSms(sendSmsRequest)
	if err != nil {
		resp.Code = int32(500)
		resp.Reason = err.Error()
		return nil, fmt.Errorf("json marshal fail, err: %w", err)
	}

	if !strings.EqualFold(strings.ToLower(*result.Body.Code), "ok") {
		resp.Code = int32(500)
	}

	if result.Body.Message != nil {
		resp.Reason = *result.Body.Message
	}
	if result.Body.BizId != nil {
		resp.ThirdMsgId = *result.Body.BizId
	}
	return
}

func (a *aliPlugin) DoJob(ctx context.Context) (err error) {
	a.CheckTplStatus()
	// a.CheckSendStatus()
	return nil
}

var contentRegex = regexp.MustCompile("(#\\w+?#)")

func (a *aliPlugin) selfToThird(content string) (result string) {
	result = content
	submatch := contentRegex.FindAllStringSubmatch(content, -1)
	for i := range submatch {
		if submatch[i] != nil && len(submatch[i]) > 0 {
			s := submatch[i][0]
			k := strings.ReplaceAll(s, "#", "")
			result = strings.ReplaceAll(result, s, fmt.Sprintf("${%s}", k))
		}
	}
	return
}

func (a *aliPlugin) CheckTplStatus() {
	tpls := make([]mysql.NotifyTplChannel, 0)
	err := a.db.Model(mysql.NotifyTplChannel{}).Where("`ch_status` = ? and channel_id = ?", int(commonv1.NOTIFY_CS_CHECKING), a.channelId).Find(&tpls).Error
	if err != nil {
		elog.Error("query error", elog.FieldErr(err))
		return
	}

	if len(tpls) == 0 {
		return
	}

	for i := range tpls {
		tpl := tpls[i]
		querySmsTemplateRequest := &dysmsapi20170525.QuerySmsTemplateRequest{
			TemplateCode: tea.String(tpl.ThirdTplId),
		}
		result, err := a.client.QuerySmsTemplate(querySmsTemplateRequest)
		if err != nil {
			elog.Error("ali QuerySmsTemplate error", elog.FieldErr(err))
			continue
		}
		// 0：审核中。
		// 1：审核通过。
		// 2：审核失败，请在返回参数Reason中查看审核失败原因
		ups := make(map[string]interface{}, 0)

		if *result.Body.TemplateStatus == 1 {
			ups["ch_status"] = commonv1.NOTIFY_CS_SUCCESS
		} else if *result.Body.TemplateStatus == 2 {
			ups["ch_status"] = commonv1.NOTIFY_CS_FAIL
		}
		elog.Info("res = ", elog.FieldValueAny(result))
		err = a.db.Model(mysql.NotifyTplChannel{}).Where("`id` = ?", tpl.ID).Updates(ups).Error
		if err != nil {
			elog.Error("ali 更新模板状态错误", elog.FieldErr(err))
			continue
		}
	}
	return
}

func (a *aliPlugin) CheckSendStatus() {
	elog.Info("start check ali msg")
	infos := make([]mysql.Notify, 0)
	err := a.db.
		Where("`status` = ? and channel = ? and msg_id <> '' and final_content<>'mock'", commonv1.NOTIFY_STATUS_SENDING, a.channelId).
		Find(&infos).Error
	if err != nil {
		elog.Error("query error", elog.FieldErr(err))
		return
	}
	elog.Info("got sending msg ", zap.Int("length", len(infos)))
	for i := range infos {
		info := infos[i]
		sendDate := time.Unix(info.Ctime, 0).Format("20060102")
		querySendDetailsRequest := &dysmsapi20170525.QuerySendDetailsRequest{
			PhoneNumber: tea.String(info.Phone),
			BizId:       tea.String(info.MsgId),
			SendDate:    tea.String(sendDate),
			PageSize:    tea.Int64(50), // page相关参数先写死，使用 bizId 无需翻页
			CurrentPage: tea.Int64(1),
		}
		result, err := a.client.QuerySendDetails(querySendDetailsRequest)
		if err != nil {
			errorlog := "ali QuerySendDetails " + err.Error()
			elog.Error("ali QuerySendDetails fail", zap.Error(err))
			skipInfo(a.db, info, errorlog)
			continue
		}
		if !strings.EqualFold(strings.ToLower(*result.Body.Code), "ok") {
			errorlog := "ali QuerySendDetails error"
			elog.Error("ali QuerySendDetails fail", zap.Any("result", result), zap.String("msgId", info.MsgId))
			skipInfo(a.db, info, errorlog)
			continue
		}
		// 1：等待回执。
		// 2：发送失败。
		// 3：发送成功。
		if result.Body.SmsSendDetailDTOs.SmsSendDetailDTO == nil || len(result.Body.SmsSendDetailDTOs.SmsSendDetailDTO) == 0 {
			errorlog := "ali QuerySendDetails response empty"
			elog.Error("ali QuerySendDetails response empty fail", zap.Any("result", result), zap.String("msgId", info.MsgId))
			skipInfo(a.db, info, errorlog)
			continue
		}

		// 只使用单条发送 取1条
		dto := result.Body.SmsSendDetailDTOs.SmsSendDetailDTO[0]

		if dto.ReceiveDate == nil {
			skipInfo(a.db, info, "ali parse ReceiveDate error"+err.Error())
			continue
		}
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", *dto.ReceiveDate, time.Local)
		if err != nil {
			elog.Error("ali parse ReceiveDate fail", zap.Error(err), zap.String("msgId", info.MsgId))
			skipInfo(a.db, info, "ali parse ReceiveDate error"+err.Error())
			continue
		}
		ups := make(map[string]interface{}, 0)
		ups["utime"] = parsed.Unix()

		ups["error_log"] = *dto.ErrCode
		if *dto.ErrCode == "DELIVERED" || *dto.SendStatus == 3 {
			// 发送成功
			ups["status"] = commonv1.NOTIFY_STATUS_SUCCESS
		} else if *dto.SendStatus == 2 {
			ups["status"] = commonv1.NOTIFY_STATUS_ERROR
		} else if *dto.SendStatus == 1 {
			ups["status"] = commonv1.NOTIFY_STATUS_SENDING
		} else {
			ups["status"] = commonv1.NOTIFY_STATUS_INVALID
		}

		err = a.db.Model(mysql.Notify{}).
			Where("id = ?", info.ID).
			Updates(ups).Error

		if err != nil {
			elog.Errorf("update notify error msg_id=%s, channel_id=%s, ")
		} else {
			elog.Info("update notify success")
		}

	}
}

func skipInfo(db *gorm.DB, info mysql.Notify, errorlog string) {
	ts := time.Now().Unix()
	if ts-info.Ctime > 600 {
		ups := make(map[string]interface{})
		ups["msg_status"] = commonv1.NOTIFY_STATUS_ERROR
		ups["error_log"] = errorlog
		// 10 分钟取不到状态则默认失败
		db.Model(mysql.Notify{}).Where("id = ?", info.ID).Updates(ups)
	}
}
