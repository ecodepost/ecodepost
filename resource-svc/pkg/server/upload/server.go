package upload

import (
	"bytes"
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/service"

	errcodev1 "ecodepost/pb/errcode/v1"
	uploadv1 "ecodepost/pb/upload/v1"
)

type GrpcServer struct{}

var _ uploadv1.UploadServer = (*GrpcServer)(nil)

// GetToken 获取一次性上传文件的token
func (GrpcServer) GetToken(ctx context.Context, req *uploadv1.GetTokenReq) (*uploadv1.GetTokenRes, error) {
	config := invoker.AliSts.GetConfig()
	// 获取credential
	publicCredentials, err := invoker.AliSts.GetToken(900, config.Bucket, req.GetUid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("GetToken fail,err:" + err.Error())
	}

	output := make([]*uploadv1.Token, 0)
	output = append(output, &uploadv1.Token{
		Region:          config.RegionId,
		AccessKeyId:     publicCredentials.AccessKeyId,
		AccessKeySecret: publicCredentials.AccessKeySecret,
		StsToken:        publicCredentials.SecurityToken,
		Bucket:          config.Bucket,
		Expiration:      publicCredentials.Expiration,
	})
	//output = append(output, &uploadv1.Token{
	//	Region:          config.Region,
	//	AccessKeyId:     secretCredentials.AccessKeyId,
	//	AccessKeySecret: secretCredentials.AccessKeySecret,
	//	StsToken:        secretCredentials.SecurityToken,
	//	Bucket:          config.PrivateBucket.Name,
	//	Expiration:      secretCredentials.Expiration,
	//})

	return &uploadv1.GetTokenRes{List: output}, nil
}

// GetPath 获取上传文件的Path
func (GrpcServer) GetPath(ctx context.Context, req *uploadv1.GetPathReq) (*uploadv1.GetPathRes, error) {
	return service.File.GetUploadPath(ctx, req)
}

// UploadLocalFile 上传本地文件
func (GrpcServer) UploadLocalFile(ctx context.Context, req *uploadv1.UploadLocalFileReq) (*uploadv1.UploadLocalFileRes, error) {
	// 调用 sdk
	reader := bytes.NewReader(req.File)
	err := invoker.AliSts.PutObject(req.DstPath, reader)
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("File Save,err:" + err.Error())
	}
	return &uploadv1.UploadLocalFileRes{}, nil
}

// ShowImage 上传本地文件
func (GrpcServer) ShowImage(ctx context.Context, req *uploadv1.ShowImageReq) (*uploadv1.ShowImageRes, error) {
	// 调用 sdk
	object, err := invoker.AliSts.GetObject(req.Path)
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("ShowImage,err:" + err.Error())
	}
	return &uploadv1.ShowImageRes{
		File: object,
	}, nil
}
