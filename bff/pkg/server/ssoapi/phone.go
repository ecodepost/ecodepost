package ssoapi

import (
	"fmt"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/service/ssoservice"
	"ecodepost/bff/pkg/util"
	notifyv1 "ecodepost/pb/notify/v1"
	"ecodepost/pb/user/v1"
	"github.com/gotomicro/ego/core/econf"
)

type SendPhoneCodeRequest struct {
	Phone string `json:"phone"`
}

func RegisterPhoneCode(c *bffcore.Context) {
	var req SendPhoneCodeRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err.Error())
		return
	}
	codeTTL, needSend := phoneCode(c, req.Phone, ssoservice.RegisterCodeType)
	if !needSend {
		return
	}
	c.JSONOK(phoneCodeRes{TTL: codeTTL})
}

func LoginPhoneCode(c *bffcore.Context) {
	var req SendPhoneCodeRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	codeTTL, needSend := phoneCode(c, req.Phone, ssoservice.LoginCodeType)
	if !needSend {
		return
	}
	c.JSONOK(phoneCodeRes{TTL: codeTTL})
}

func RetrievePhoneCode(c *bffcore.Context) {
	var req SendPhoneCodeRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	codeTTL, needSend := phoneCode(c, req.Phone, ssoservice.RetrievePasswordCodeType)
	if !needSend {
		return
	}
	c.JSONOK(phoneCodeRes{TTL: codeTTL})
}

type phoneCodeRes struct {
	TTL int64 `json:"ttl"`
}

func phoneCode(c *bffcore.Context, phone string, codeType ssoservice.CodeType) (codeTTL int64, needRes bool) {
	if !util.CheckMobile(phone) {
		c.JSONE(1, "手机号不正确", phone)
		return
	}

	// 注册的验证码需要校验，手机号用过没
	if codeType == ssoservice.RegisterCodeType {
		userInfo, err := invoker.GrpcUser.InfoByPhone(c.Request.Context(), &userv1.InfoByPhoneReq{Phone: phone})
		if err != nil {
			c.JSONE(1, "系统错误", fmt.Errorf("user register get user info failed,err: %w", err))
			return
		}
		if userInfo.User.GetUid() > 0 {
			c.JSONE(1, "该用户已经注册", fmt.Errorf("该用户已经注册"))
			return
		}
	}

	//var msgCode string
	var isSent bool
	var err error
	var msgCode string
	msgCode, codeTTL, isSent, err = ssoservice.Code.SendCode(c.Request.Context(), phone, c.ClientIP(), codeType)
	if err != nil {
		c.JSONE(1, "生成code失败", err)
		return
	}
	if isSent {
		c.JSONOK(phoneCodeRes{TTL: codeTTL})
		return
	}

	// 只有线上环境才需要发送验证码短信
	if econf.GetString("mode") == "pro" {
		if _, err := invoker.GrpcNotify.SendMsg(c.Request.Context(), &notifyv1.SendMsgReq{
			TplId:  2, // 默认sms模板
			Msgs:   []*notifyv1.Msg{{Receiver: phone, Vars: map[string]string{"code": msgCode}}},
			VarSms: &notifyv1.Sms{},
		}); err != nil {
			c.JSONE(1, "发送code失败", err)
			return
		}
	}
	needRes = true
	return
}
