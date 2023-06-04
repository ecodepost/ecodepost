package service

import (
	"context"
	"time"

	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/utils/slate"
)

func (f *file) CreateDocument(ctx context.Context, req *articlev1.CreateDocumentReq) (createFile *mysql.File, err error) {
	// 默认当前节点为根节点, 父节点为缺省值
	var node = commonv1.FILE_NODE_ROOT
	var parentNode = commonv1.FILE_NODE_INVALID

	nowTime := time.Now().Unix()
	ctime := req.GetCtime()
	utime := req.GetUtime()
	if ctime == 0 {
		ctime = nowTime
	}
	if utime == 0 {
		utime = nowTime
	}

	createFile, err = f.CreateFile(ctx, CreateOrCopyFileReq{
		Name:       req.GetName(),
		Uid:        req.GetUid(),
		Content:    req.Content,
		SpaceGuid:  req.GetSpaceGuid(),
		HeadImage:  req.GetHeadImage(),
		CreateTime: ctime,
		UpdateTime: utime,
		CntView:    req.GetCntView(),
		Node:       node,
		ParentNode: parentNode,
		FileFormat: req.GetFormat(),
		Ip:         req.Ip,
	}, commonv1.FILE_TYPE_DOCUMENT)
	if err != nil {
		return nil, err
	}
	return
}

type CopyDocumentRequest struct {
	// 用户id
	Uid int64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	// 社区guid
	CmtGuid string `protobuf:"bytes,2,opt,name=cmtGuid,proto3" json:"cmtGuid,omitempty"`
	// 文件名称
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// 空间guid
	SpaceGuid string `protobuf:"bytes,4,opt,name=spaceGuid,proto3" json:"spaceGuid,omitempty"`
	// 父级文件 ID
	ParentGuid string `protobuf:"bytes,5,opt,name=parentGuid,proto3" json:"parentGuid,omitempty"`
	// 内容
	ContentKey string `protobuf:"bytes,7,opt,name=content,proto3" json:"content,omitempty"`
	// headImage
	HeadImage string `protobuf:"bytes,8,opt,name=headImage,proto3" json:"headImage,omitempty"`
}

//
//// CopyDocument 使用模板的时候，直接复制模板，目前复制文档只有富文本类型
//func (f *file) CopyDocument(ctx context.Context, req CopyDocumentRequest) (createFile *mysql.File, err error) {
//	// fileType := commonv1.FILE_TYPE_DOCUMENT
//	// var parentPointer *mysql.File
//	// var parentEdge mysql.Edge
//	var node, parentNode = commonv1.FILE_NODE_ROOT, commonv1.FILE_NODE_INVALID
//
//	// 如果大于0 ，说明存在父亲节点
//	// if req.ParentGuid != "" {
//	// node = commonv1.FILE_NODE_LEAF
//	// var parentFile mysql.File
//	// parentFile, err = mysql.FileInfoByGuidAndCmtGuid(invoker.Db.WithContext(ctx), req.CmtGuid, req.ParentGuid)
//	// if err != nil {
//	//	return nil, fmt.Errorf("file create fail, err: %w", err)
//	// }
//	// if parentFile.Id == 0 {
//	//	return nil, fmt.Errorf("file not exist")
//	// }
//	// // 如果父节点无父节点, 则父节点为根节点, 否则父节点为中间节点
//	// if parentFile.ParentGuid == "" {
//	//	parentNode = commonv1.FILE_NODE_ROOT
//	// } else {
//	//	parentNode = commonv1.FILE_NODE_INNER
//	// }
//	// // 更新父节点信息
//	// if e := mysql.FileUpdate(invoker.Db.WithContext(ctx), parentFile.CmtGuid, parentFile.Guid, map[string]any{"file_node": parentNode}); e != nil {
//	//	return nil, fmt.Errorf("file parent file update fail, err: %w", e)
//	// }
//
//	// parentPointer = &parentFile
//	// parentEdge, err = mysql.EdgeInfoByFileGuidAndFileType(invoker.Db.WithContext(ctx), req.CmtGuid, parentFile.Guid, parentFile.FileType)
//	// if err != nil {
//	//	return nil, fmt.Errorf("file create get edge fail, err: %w", err)
//	// }
//	// }
//	nowTime := time.Now().Unix()
//	// 复制File
//	createFile, err = f.CreateFile(ctx, CreateOrCopyFileReq{
//		Name:       req.Name,
//		Uid:        req.Uid,
//		SpaceGuid:  req.SpaceGuid,
//		HeadImage:  req.HeadImage,
//		Node:       node,
//		ParentNode: parentNode,
//		ContentKey: req.ContentKey,
//		CreateTime: nowTime,
//		UpdateTime: nowTime,
//		FileFormat: commonv1.FILE_FORMAT_DOCUMENT_SLATE,
//	}, commonv1.FILE_TYPE_DOCUMENT)
//	if err != nil {
//		return nil, err
//	}
//
//	// edge := &mysql.Edge{
//	//	SpaceGuid:  req.SpaceGuid,
//	//	FileType:   fileType,
//	//	FileNode:   node,
//	//	FileGuid:   createFile.Guid,
//	//	ParentType: fileType,
//	//	ParentGuid: "",
//	//	RootType:   fileType,
//	//	RootGuid:   createFile.Guid,
//	//	Sort:       time.Now().UnixMilli(),
//	//	CmtGuid:    req.CmtGuid,
//	// }
//	//
//	// // 如果存在父亲节点，那么就将父亲节点数据放上去
//	// if parentPointer != nil {
//	//	edge.ParentType = parentPointer.FileType
//	//	edge.ParentGuid = parentPointer.Guid
//	//	edge.ParentNode = parentNode
//	//	edge.RootType = parentEdge.RootType
//	//	edge.RootGuid = parentEdge.RootGuid
//	// }
//
//	// err = mysql.EdgeCreate(invoker.Db.WithContext(ctx), edge)
//	// if err != nil {
//	//	return nil, fmt.Errorf("file create edge create fail, err: %w", err)
//	// }
//	return
//}

func (f *file) UpdateDocument(ctx context.Context, req *articlev1.UpdateDocumentReq) (err error) {
	// 新创建的file只有slate json模式
	// gocn的有富文本格式
	if req.FileFormat == commonv1.FILE_FORMAT_DOCUMENT_RICH {
		if req.Content != nil {
			slateJsonBytes, err := slate.HtmlToSlateJson(*req.Content)
			if err != nil {
				return errcodev1.ErrInternal().WithMessage("html to slate json fail, err: " + err.Error())
			}
			req.Content = &slateJsonBytes
		}
		req.FileFormat = commonv1.FILE_FORMAT_DOCUMENT_SLATE
	}

	_, err = f.UpdateFile(ctx, UpdateFileRequest{
		Uid:        req.GetUid(),
		FileGuid:   req.GetGuid(),
		NewContent: req.Content,
		HeadImage:  req.HeadImage,
	}, map[string]any{
		"name": req.GetName(),
	})
	return
}
