package bffcore

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"ecodepost/bff/pkg/consts"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eerrors"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/transport"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
)

// HandlerFunc core封装后的handler
type HandlerFunc func(c *Context)

var pbMarshaler = jsonpb.MarshalOptions{
	// EnumsAsInts:  true,
}

type Pagination struct {
	// 当前页数
	CurrentPage int32 `json:"currentPage" form:"currentPage"`
	// 每页总数
	PageSize int32 `json:"pageSize" form:"pageSize"`
	// 排序字符串
	Sort string `json:"sort" form:"sort"`
}

type ListPage struct {
	// 列表
	List any `json:"list"`
}

func (p Pagination) ToPb() *commonv1.Pagination {
	pg := &commonv1.Pagination{}
	if p.CurrentPage == 0 {
		pg.CurrentPage = 1
	}
	if p.PageSize == 0 {
		pg.PageSize = 10
	}
	return pg
}

// Handle 将core.HandlerFunc转换为gin.HandlerFunc
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			c,
		}
		h(ctx)
	}
}

// Context core封装后的Context
type Context struct {
	*gin.Context
}

const (
	// CodeOK 表示响应成功状态码
	CodeOK = 0
	// CodeErr 表示默认响应失败状态码
	CodeErr     = 1
	CodeInvalid = 401
)

// Ctx 返回Request.Context()
func (c *Context) Ctx() context.Context {
	return c.Request.Context()
}

// InjectGuid 注入guid
func (c *Context) InjectGuid(guid string) *Context {
	c.Request = c.Request.WithContext(transport.WithValue(c.Ctx(), consts.XLinkSpaceGuid, guid))
	return c
}

// InjectSpc 注入spaceGuid
func (c *Context) InjectSpc(spcGuid string) *Context {
	c.Request = c.Request.WithContext(transport.WithValue(c.Ctx(), consts.XLinkSpaceSpc, spcGuid))
	return c
}

// Res 标准JSON输出格式
type Res struct {
	// Code 响应的业务错误码。0表示业务执行成功，非0表示业务执行失败。
	Code int `json:"code"`
	// Msg 响应的参考消息。前端可使用msg来做提示
	Msg string `json:"msg"`
	// Data 响应的具体数据
	Data any `json:"data"`
}

// ResProto 标准JSON输出格式
type ResProto struct {
	// Code 响应的业务错误码。0表示业务执行成功，非0表示业务执行失败。
	Code int `json:"code"`
	// Msg 响应的参考消息。前端可使用msg来做提示
	Msg string `json:"msg"`
	// Data 响应的具体数据
	Data json.RawMessage `json:"data"`
}

// ResPage 带分页的标准JSON输出格式
type ResPage struct {
	Res
	Pagination *commonv1.Pagination `json:"pagination"`
}

// JSON 输出响应JSON
// 形如 {"code":<code>, "msg":<msg>, "data":<data>}
func (c *Context) JSON(httpStatus int, res Res) {
	c.Context.JSON(httpStatus, res)
}

// JSONOK 输出响应成功JSON，如果data不为零值，则输出data
// 形如 {"code":0, "msg":"成功", "data":<data>}
func (c *Context) JSONOK(data ...interface{}) {
	j := new(Res)
	j.Code = CodeOK
	j.Msg = "成功"
	if len(data) > 0 {
		j.Data = data[0]
	} else {
		j.Data = ""
	}
	c.Context.JSON(http.StatusOK, j)
	return
}

// ProtoJSONOK 输出响应成功JSON，如果data不为零值，则输出data
// 形如 {"code":0, "msg":"成功", "data":<data>}
func (c *Context) ProtoJSONOK(data proto.Message) {
	j := new(ResProto)
	j.Code = CodeOK
	j.Msg = "成功"
	bts, _ := pbMarshaler.Marshal(data)
	j.Data = bts
	c.Context.JSON(http.StatusOK, j)
	return
}

// JSONE 输出失败响应
// 形如 {"code":<code>, "msg":<msg>, "data":<data>}
func (c *Context) JSONE(code int, msg string, data interface{}) {
	j := new(Res)
	j.Code = code
	j.Msg = msg

	if econf.GetBool("debug") == true {
		switch d := data.(type) {
		case error:
			j.Data = d.Error()
		default:
			j.Data = data
		}
	}

	elog.Warn("biz warning", elog.FieldValue(msg), elog.FieldValueAny(data), elog.FieldTid(etrace.ExtractTraceID(c.Ctx())))
	c.Context.JSON(http.StatusOK, j)
	return
}

type ResPageData struct {
	// List item列表
	List []any `json:"list"`
	// 分页数据
	Pagination commonv1.Pagination `json:"pagination"`
}

// JSONListPage 输出带分页信息的JSON
func (c *Context) JSONListPage(data any, pagination *commonv1.Pagination) {
	j := new(ResPage)
	j.Code = CodeOK
	j.Data = ListPage{
		List: data,
	}
	j.Pagination = pagination
	c.Context.JSON(http.StatusOK, j)
}

// Bind 将请求消息绑定到指定对象中，并做数据校验。如果校验失败，则返回校验失败错误中文文案
// 并将HTTP状态码设置成400
func (c *Context) Bind(obj interface{}) (err error) {
	return validate(c.Context.Bind(obj))
}

// ShouldBind 将请求消息绑定到指定对象中，并做数据校验。如果校验失败，则返回校验失败错误中文文案
// 类似Bind，但是不会将HTTP状态码设置成400
func (c *Context) ShouldBind(obj interface{}) (err error) {
	return validate(c.Context.ShouldBind(obj))
}

func (c *Context) EgoJsonI18N(err error) {
	egoErr := eerrors.FromError(err)
	msg := errcodev1.ReasonI18n(egoErr, c.GetString(ContextLanguage))
	if msg == "" {
		msg = err.Error()
	}
	j := new(Res)
	j.Code = 1
	if int(errcodev1.Error_value[strings.TrimPrefix(egoErr.GetReason(), "errcode.v1.")]) != 0 {
		j.Code = int(errcodev1.Error_value[strings.TrimPrefix(egoErr.GetReason(), "errcode.v1.")])
	}
	j.Msg = msg

	if econf.GetBool("debug") == true {
		j.Data = err.Error()
	} else {
		j.Data = egoErr.GetReason()
	}

	j.Data = egoErr.GetReason()
	elog.Warn("biz warning", elog.FieldValue(msg), elog.FieldErr(err), elog.FieldTid(etrace.ExtractTraceID(c.Ctx())))
	c.Context.JSON(http.StatusOK, j)
}
