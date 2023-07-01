package ssoapi

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	userv1 "ecodepost/pb/user/v1"
	"github.com/ego-component/ewechat/oauth"
	"github.com/google/uuid"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	ua "github.com/mssola/user_agent"
	"go.uber.org/zap"
)

const (
	UserOpenWechatWeb  int32 = 1
	UserOpenWechatH5         = 2
	UserOpenWechatApp        = 3
	UserOpenWechatMini       = 4
)

const (
	wcStateKeyPrefix = "wc:st:"
)

func wcStateKey(stateKey string) string {
	return wcStateKeyPrefix + stateKey
}

func LoginWechat(c *bffcore.Context) {
	var reqView dto.OauthRequest
	if err := c.Bind(&reqView); err != nil {
		elog.Error("json marshal request body error", zap.Error(err), elog.FieldCtxTid(c.Request.Context()))
		c.JSONE(1, "参数错误", err)
		return
	}

	state, err := json.Marshal(reqView)
	if err != nil {
		c.JSONE(2, "编码错误", err)
		return
	}

	stateKey := strings.ReplaceAll(uuid.New().String(), "-", "")
	if e := invoker.Redis.Set(c.Request.Context(), wcStateKey(stateKey), state, 300*time.Second); e != nil {
		c.JSONE(2, "redis set state fail", err)
		return
	}
	// sEnc := base64.RawURLEncoding.EncodeToString(state)
	if ua.New(c.GetHeader("User-Agent")).Mobile() {
		url, _ := invoker.WechatOauthH5.GetRedirectURL(econf.GetString("wechat2.redirectUri"), econf.GetString("wechat2.scope"), stateKey)
		c.Redirect(http.StatusFound, url)
	} else {
		c.Redirect(http.StatusFound, invoker.WechatOauthWeb.AuthCodeURL(econf.GetString("wechat1.redirectUri"), econf.GetString("wechat1.scope"), stateKey))
	}
	return
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

func WechatValidate(c *bffcore.Context) {
	signature := c.Query("signature")
	arr := []string{c.Query("timestamp"), invoker.WechatOauthH5.Token, c.Query("nonce")}
	sort.Strings(arr)
	if signature == SHA1(strings.Join(arr, "")) {
		_, _ = c.Writer.WriteString(c.Query("echostr"))
	} else {
		c.Context.JSON(401, nil)
	}
}

func CodeWechatWeb(c *bffcore.Context) {
	CodeWechat(c, "web")
}

func CodeWechatH5(c *bffcore.Context) {
	CodeWechat(c, "h5")
}

func GetWechatCmp(wechat string) *oauth.Oauth {
	if wechat == "web" {
		return invoker.WechatOauthWeb
	}
	return invoker.WechatOauthH5
}

// CodeWechat ...
//  1. 在 /code/wechat 里面判断下当前登录用户 wechatOpenID 是否已经存在。
//  2. 如果 wechatOpenID 存在，则判断下wechat open表里，uid是否存在。
//  3. 如果uid 不存在，那么跳转到前端绑定手机号页面；用户绑定手机号，完成登录（用哪个接口？）。
//  4. 如果uid存在，那么完成登录过程。
//  5. 如果 wechatOpenID 不存在，那么创建第三方用户。并跳转到前端绑定手机号页面。用户绑定手机号，完成登录（用哪个接口？。
func CodeWechat(c *bffcore.Context, wechat string) {
	code := c.Query("code")
	if code == "" {
		c.JSONE(1, "code can't empty", nil)
		return
	}

	stateKey := c.Query("state")
	resBytes, err := invoker.Redis.GetBytes(c.Request.Context(), wcStateKey(stateKey))
	if err != nil {
		c.JSONE(1, "redis get state fail, stateKey:"+stateKey+", err:"+err.Error(), err)
		return
	}

	// stateBase64 := c.Query("state")
	// resBytes, err := base64.RawURLEncoding.DecodeString(stateBase64)
	// if err != nil {
	// 	c.JSONE(1, "base64 decode err: "+err.Error(), err)
	// 	return
	// }
	reqView := dto.OauthRequest{}
	if err = json.Unmarshal(resBytes, &reqView); err != nil {
		c.JSONE(1, "json decode err: "+err.Error(), err)
		return
	}
	elog.Info("form info", zap.Any("form", c.Request.Form))

	result, err := GetWechatCmp(wechat).ExchangeTokenURL(code)
	if err != nil {
		c.JSONE(1, "exchange token err: "+err.Error(), err)
		return
	}

	userInfo, err := GetWechatCmp(wechat).GetUserInfo(result.AccessToken, result.OpenID, "")
	if err != nil {
		c.JSONE(1, "get user info err: "+err.Error(), err)
		return
	}

	if reqView.RedirectUri == "" {
		c.JSONE(1, "参数错误：redirect uri 不能为空", nil)
		return
	}

	grpcUser, err := invoker.GrpcUser.LoginUserOpen(c.Context, &userv1.LoginUserOpenReq{
		Genre:    WechatGenre(wechat),
		OpenId:   userInfo.OpenID,
		UnionId:  userInfo.Unionid,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.HeadImgURL,
		Sex:      userInfo.Sex,
		Country:  userInfo.Country,
		Province: userInfo.Province,
		City:     userInfo.City,
	})
	if err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}

	// uid不为0说明已创建user记录, 且绑定了手机号
	if grpcUser.Uid != 0 {
		//ssoUser := ssov1.User{
		//	Uid:      grpcUser.Uid,
		//	Nickname: grpcUser.Nickname,
		//	Username: grpcUser.Nickname,
		//	Avatar:   grpcUser.Avatar,
		//	Email:    "",
		//}
		responseSsoLogin(c, grpcUser.Uid, true)
		return
	}

	// uid为0说明未创建user记录, 且未绑定手机号, 则跳转到手机绑定页面，绑定手机并创建用户
	wechatTokenBytes, _ := json.Marshal(wechatTokenPayload{UnionID: userInfo.Unionid})
	token := strings.ReplaceAll(uuid.New().String(), "-", "")
	if e := invoker.Redis.Set(c.Request.Context(), wechatTokenKey(token), string(wechatTokenBytes), 600*time.Second); e != nil {
		c.JSONE(1, "wechat token set fail, err:"+e.Error(), e)
		return
	}
	// 流程上是 bff -> oauth2 -> sso -> oauth2 -> wechat
	// 那么微信第一次oauth2 返回给 sso后
	// sso还需要在登录的时候，反解析，在oauth2授权后给bff
	// 需要携带state信息，用于给后面手机绑定登录后，做oauth2认证的跳转
	c.Redirect(http.StatusFound, "/sso/loginBind?client_id="+reqView.ClientId+"&redirect_uri="+reqView.RedirectUri+"&token="+token+"&referer="+reqView.Referer+"&response_type="+reqView.ResponseType+"&scope="+reqView.Scope+"&state="+reqView.State)
}

func WechatGenre(wechat string) int32 {
	if wechat == "web" {
		return UserOpenWechatWeb
	}
	// TODO 后续可能再细分
	return UserOpenWechatH5
}

func wechatTokenKey(token string) string {
	return "wc:t:" + token
}

type wechatTokenPayload struct {
	UnionID string
}
