package file

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	uploadv1 "ecodepost/pb/upload/v1"
	"encoding/base64"
)

type UploadLocalFileReq struct {
	DstPath string `json:"dstPath" binding:"required"`
	File    string `json:"file" binding:"required"`
}

func UploadLocalFile(c *bffcore.Context) {
	var req UploadLocalFileReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err.Error())
		return
	}

	imageFile, err := base64.StdEncoding.DecodeString(req.File)
	if err != nil {
		c.JSONE(1, "转换错误", err.Error())
		return
	}

	_, err = invoker.GrpcUpload.UploadLocalFile(c, &uploadv1.UploadLocalFileReq{
		DstPath: req.DstPath,
		File:    imageFile,
	})

	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	c.JSONOK()
}
