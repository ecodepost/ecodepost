package service

import (
	"context"
	"fmt"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
	questionv1 "ecodepost/pb/question/v1"
)

func (f *file) CreateQuestion(ctx context.Context, req *questionv1.CreateReq) (createFile *mysql.File, err error) {
	nowTime := time.Now().Unix()
	createFile, err = f.CreateFile(ctx, CreateOrCopyFileReq{
		Name:       req.GetName(),
		Uid:        req.GetUid(),
		Content:    req.GetContent(),
		SpaceGuid:  req.GetSpaceGuid(),
		ParentGuid: req.GetParentGuid(),
		CreateTime: nowTime,
		UpdateTime: nowTime,
		FileFormat: commonv1.FILE_FORMAT_DOCUMENT_SLATE,
		Ip:         req.Ip,
	}, commonv1.FILE_TYPE_QUESTION)
	if err != nil {
		return nil, err
	}

	if req.GetParentGuid() != "" {
		if err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetParentGuid(), "cnt_comment", 1); err != nil {
			return nil, fmt.Errorf("file create mysql create fail2, err: %w", err)
		}
	}

	return
}

type CopyQuestionRequest struct {
	Uid        int64  `json:"uid"`        // 用户id
	CmtGuid    string `json:"cmtGuid"`    // 社区guid
	Name       string `json:"name"`       // 文件名称
	SpaceGuid  string `json:"spaceGuid"`  // 空间guid
	ParentGuid string `json:"parentGuid"` // 父级文件 ID，如果存在父亲，说明是answer
	ContentKey string `json:"contentKey"` // 文件内容
	Size       int64  `json:"size"`       // 文件大小
}

//
//func (f *file) CopyQuestion(ctx context.Context, req CopyQuestionRequest) (createFile *mysql.File, err error) {
//	fileType := commonv1.FILE_TYPE_QUESTION
//	fileGuid, err := Resource.GenerateGuid(ctx, commonv1.CMN_GUID_FILE, req.Uid)
//	if err != nil {
//		return nil, fmt.Errorf("file create gen guid fail, err: %w", err)
//	}
//	nowTime := time.Now().Unix()
//	createFile = &mysql.File{
//		Guid:       fileGuid,
//		ParentGuid: req.ParentGuid,
//		SpaceGuid:  req.SpaceGuid,
//		Name:       req.Name,
//		CreatedBy:  req.Uid,
//		UpdatedBy:  req.Uid,
//		ContentKey: GenContentKey(commonv1.FILE_TYPE_QUESTION, commonv1.FILE_BIZSTS_COURSE_PUBLISHED),
//		Status:     commonv1.FILE_STATUS_SUCC,
//		FileType:   fileType,
//		FileFormat: commonv1.FILE_FORMAT_DOCUMENT_SLATE,
//		Ctime:      nowTime,
//		Utime:      nowTime,
//		Size:       req.Size,
//	}
//	tx := invoker.Db.WithContext(ctx).Begin()
//	if err = f.create(ctx, tx, createFile); err != nil {
//		tx.Rollback()
//		return nil, fmt.Errorf("file create mysql create fail, err: %w", err)
//	}
//	if req.ParentGuid != "" {
//		if err = mysql.FileUpdateExpr(tx, req.ParentGuid, "cnt_comment", 1); err != nil {
//			tx.Rollback()
//			return nil, fmt.Errorf("file create mysql create fail2, err: %w", err)
//		}
//	}
//	tx.Commit()
//	if err = invoker.AliOss.CopyObject(req.ContentKey, createFile.ContentKey); err != nil {
//		return nil, fmt.Errorf("file create oss put fail, err: %w", err)
//	}
//	return
//}

func (f *file) UpdateQuestion(ctx context.Context, req *questionv1.UpdateReq) (err error) {
	_, err = f.UpdateFile(ctx, UpdateFileRequest{
		Uid:        req.GetUid(),
		FileGuid:   req.GetGuid(),
		NewContent: req.Content,
	}, map[string]any{
		"name": req.GetName(),
	})
	if err != nil {
		return err
	}
	return
}
