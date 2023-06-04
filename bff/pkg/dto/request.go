package dto

// OauthRequest 用于register and login
type OauthRequest struct {
	Account      string `json:"account,omitempty"`                          // password type need, sms type need,  目前只有手机号可以登录
	Password     string `json:"password,omitempty"`                         // password type need
	OldPassword  string `json:"oldPassword,omitempty"`                      // reset password type need
	Code         string `json:"code,omitempty"`                             // sms type need，短信号码
	Token        string `json:"token,omitempty"`                            // 电话绑定token
	Referer      string `json:"referer,omitempty" form:"referer"`           // 当前用户访问的页面
	RedirectUri  string `json:"redirect_uri,omitempty" form:"redirect_uri"` // redirect by backend
	ClientId     string `json:"client_id,omitempty" form:"client_id"`
	ResponseType string `json:"response_type,omitempty" form:"response_type"`
	State        string `json:"state,omitempty" form:"state"`
	Scope        string `json:"scope,omitempty" form:"scope"`
	Ref          string `json:"ref,omitempty" form:"ref"` // 邀请码
}

// OauthDirectRequest 直接登录和注册
type OauthDirectRequest struct {
	Account  string `json:"account,omitempty"`  // password type need, sms type need,  目前只有手机号可以登录
	Password string `json:"password,omitempty"` // password type need
	Code     string `json:"code,omitempty"`     // sms type need，短信号码
}
