package upload

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"github.com/gotomicro/ego/core/econf"

	commonv1 "ecodepost/pb/common/v1"
	uploadv1 "ecodepost/pb/upload/v1"
)

type TokenResp struct {
	Mode string            `json:"mode"`
	List []*uploadv1.Token `json:"list"`
}

func Token(c *bffcore.Context) {
	info, err := invoker.GrpcUpload.GetToken(c.Ctx(), &uploadv1.GetTokenReq{
		Uid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(TokenResp{
		Mode: econf.GetString("oss.mode"),
		List: info.List,
	})
}

type PathReq struct {
	FileName  string               `json:"fileName" binding:"required" label:"文件名"`
	Type      commonv1.CMN_UP_TYPE `json:"type"  binding:"required" label:"类型"`
	Size      int64
	SpaceGuid string
}

func Path(c *bffcore.Context) {
	var req PathReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectSpc(req.SpaceGuid)
	// if req.Size == 0 {
	//	c.JSONE(1, "图片大小存在问题", nil)
	//	return
	// }

	// guid := c.Param("guid")
	info, err := invoker.GrpcUpload.GetPath(c.Ctx(), &uploadv1.GetPathReq{
		UploadType: req.Type,
		FileName:   req.FileName,
		Uid:        c.Uid(),
		ClientIp:   c.ClientIP(),
		Refer:      c.GetHeader("referer"),
		SpaceGuid:  req.SpaceGuid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(info)
}
