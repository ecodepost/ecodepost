package service

import (
	"context"
	"fmt"
	"time"

	errcodev1 "ecodepost/pb/errcode/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service/cache"

	columnv1 "ecodepost/pb/column/v1"
	commonv1 "ecodepost/pb/common/v1"
	"gorm.io/gorm"
)

func (f *file) CreateColumn(ctx context.Context, req *columnv1.CreateReq) (createFile *mysql.File, err error) {
	fileType := commonv1.FILE_TYPE_COLUMN
	var parentPointer *mysql.File

	// 默认当前节点为根节点, 父节点为缺省值
	var node = commonv1.FILE_NODE_ROOT
	var parentNode = commonv1.FILE_NODE_INVALID

	// 如果大于0, 说明存在父亲节点, 则此节点设置为叶子节点
	if req.GetParentGuid() != "" {
		node = commonv1.FILE_NODE_LEAF
		parentFile, err := mysql.FileInfoByGuid(invoker.Db.WithContext(ctx), req.GetParentGuid())
		if err != nil {
			return nil, fmt.Errorf("file create fail, err: %w", err)
		}

		// 如果父节点无父节点, 则父节点为根节点, 否则父节点为中间节点
		if parentFile.ParentGuid == "" {
			parentNode = commonv1.FILE_NODE_ROOT
		} else {
			parentNode = commonv1.FILE_NODE_INNER
		}
		// 更新父节点信息
		if e := mysql.FileUpdate(invoker.Db.WithContext(ctx), parentFile.Guid, map[string]any{"file_node": parentNode}); e != nil {
			return nil, fmt.Errorf("file parent file update fail, err: %w", e)
		}
		parentPointer = &parentFile
	}
	nowTime := time.Now().Unix()
	ctime := req.GetCtime()
	utime := req.GetUtime()
	if ctime == 0 {
		ctime = nowTime
	}
	if utime == 0 {
		utime = nowTime
	}
	// 创建File
	createFile, err = f.CreateFile(ctx, CreateOrCopyFileReq{
		Name:       req.GetName(),
		Uid:        req.GetUid(),
		Content:    req.GetContent(),
		SpaceGuid:  req.GetSpaceGuid(),
		HeadImage:  req.GetHeadImage(),
		CreateTime: ctime,
		UpdateTime: utime,
		CntView:    req.GetCntView(),
		Node:       node,
		ParentNode: parentNode,
		ParentGuid: req.ParentGuid,
		Ip:         req.Ip,
		FileFormat: commonv1.FILE_FORMAT_DOCUMENT_SLATE,
	}, fileType)
	if err != nil {
		return nil, err
	}

	edge := &mysql.Edge{
		SpaceGuid:  req.GetSpaceGuid(),
		FileType:   fileType,
		FileGuid:   createFile.Guid,
		ParentType: commonv1.FILE_TYPE_COLUMN,
		Sort:       time.Now().UnixMilli(),
	}

	// 如果存在父亲节点，那么就将父亲节点数据放上去
	if parentPointer != nil {
		edge.ParentType = parentPointer.FileType
		edge.ParentGuid = parentPointer.Guid
	}

	// 更新edge
	if err = mysql.EdgeCreate(invoker.Db.WithContext(ctx), edge); err != nil {
		return nil, fmt.Errorf("file create edge create fail, err: %w", err)
	}
	return
}

func (f *file) UpdateColumn(ctx context.Context, req *columnv1.UpdateReq) (err error) {
	_, err = f.UpdateFile(ctx, UpdateFileRequest{
		Uid:        req.GetUid(),
		FileGuid:   req.GetFileGuid(),
		NewContent: req.Content,
		HeadImage:  req.HeadImage,
	}, map[string]any{
		"name": req.GetName(),
	})
	return
}

// ChangeSort 更改顺序
// 顺序是按照sort asc，越小越往前
// 1 放到某个file guid后面
// 2 放到某个parent guid下面
func (f *file) ChangeSort(ctx context.Context, currentGuid string, targetGuid *string, dropPosition *string, parentGuid string) (err error) {
	if currentGuid == "" {
		return fmt.Errorf("current guid is empty")
	}
	if targetGuid != nil {
		if currentGuid == *targetGuid {
			return fmt.Errorf("current guid cant eq after guid")
		}
	}

	if currentGuid == parentGuid {
		return fmt.Errorf("current guid cant eq parent guid")
	}

	spaceGuid, err := mysql.FileSpaceGuidByGuid(invoker.Db.WithContext(ctx), currentGuid)
	if err != nil {
		return fmt.Errorf("file ChangeSortByTargetGuid fail1, err: %w", err)
	}
	if spaceGuid == "" {
		return fmt.Errorf("file ChangeSortByTargetGuid space guid not exist")
	}
	spaceInfo, err := mysql.SpaceGetInfoByGuid(invoker.Db.WithContext(ctx), "guid,type", spaceGuid)
	if err != nil {
		return fmt.Errorf("file ChangeSortByTargetGuid space info fail, err: %w", err)
	}

	if targetGuid != nil {
		return f.ChangeSortByTargetGuid(ctx, spaceInfo, currentGuid, *targetGuid, dropPosition)
	}

	if parentGuid == "" {
		return fmt.Errorf("parent guid is empty")
	}

	return f.ChangeSortByParentGuid(ctx, spaceInfo, currentGuid, parentGuid)
}

// ChangeSortByTargetGuid 放到某个file guid后面
func (*file) ChangeSortByTargetGuid(ctx context.Context, spaceInfo mysql.Space, currentGuid string, targetGuid string, dropPosition *string) (err error) {
	var targetSort int64
	var targetParentGuid string
	var targetParentType commonv1.FILE_TYPE

	needSearchArr := []string{currentGuid, targetGuid}
	// 找到父级节点的guid信息
	targetParentGuid, err = mysql.GetNodeParentGuid(invoker.Db.WithContext(ctx), targetGuid)
	if err != nil {
		return fmt.Errorf("file change sort fail, err: %w", err)
	}
	if targetParentGuid != "" {
		needSearchArr = append(needSearchArr, targetParentGuid)
	}

	list, err := mysql.GetNodeSortInfoByIn(invoker.Db.WithContext(ctx), needSearchArr)
	if err != nil {
		return fmt.Errorf("ChangeSortByTargetGuid fail, err: %w", err)
	}
	if currentInfo := list.FindByGuid(currentGuid); currentInfo == nil {
		return fmt.Errorf("ChangeSortByTargetGuid find current info fail")
	}
	targetInfo := list.FindByGuid(targetGuid)
	if targetInfo == nil {
		return fmt.Errorf("ChangeSortByTargetGuid find after info fail")
	}
	targetSort = targetInfo.Sort
	// 如果after节点，存在parent节点
	if targetParentGuid != "" {
		targetParentInfo := list.FindByGuid(targetParentGuid)
		if targetParentInfo == nil {
			return fmt.Errorf("ChangeSortByTargetGuid find after info fail")
		}
		targetParentType = targetInfo.ParentType
	}

	var sign = ""
	var judgeSign = ""
	var updateSort = targetSort
	if dropPosition != nil {
		if *dropPosition == "before" {
			sign = "-"
			judgeSign = "<"
			updateSort--
		} else if *dropPosition == "after" {
			sign = "+"
			judgeSign = ">"
			updateSort++
		} else {
			return fmt.Errorf("dropPosition is invalid")
		}
	}

	return invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先把数据+1
		err = tx.Model(&mysql.Edge{}).Where("parent_guid = ? and parent_type = ? and sort "+judgeSign+" ?", targetParentGuid, targetParentType.Number(), targetSort).
			Update("sort", gorm.Expr("`sort` "+sign+" 1")).Error
		if err != nil {
			return fmt.Errorf("ChangeSortByTargetGuid update sort fail, err: %w", err)
		}

		// 更新排序
		// 并且将其放在该树形结构下面
		err = tx.Model(&mysql.Edge{}).Where("file_guid = ?", currentGuid).
			Updates(map[string]any{"sort": updateSort, "parent_guid": targetParentGuid, "parent_type": targetParentType.Number()}).Error
		if err != nil {
			return fmt.Errorf("ChangeSortByTargetGuid update sort fail2, err: %w", err)
		}

		if err = mysql.FileUpdate(tx, currentGuid, map[string]any{
			"parent_guid": targetParentGuid,
		}); err != nil {
			return errcodev1.ErrDbError().WithMessage("FileUpdate fail, err: " + err.Error())
		}
		err = cache.FileCache.SetParentGuid(ctx, tx, currentGuid, targetParentGuid)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("SetParentGuid fail2, err: " + err.Error())
		}
		return nil
	})
}

func (*file) ChangeSortByParentGuid(ctx context.Context, spaceInfo mysql.Space, currentGuid string, parentGuid string) (err error) {
	needSearchArr := []string{currentGuid, parentGuid}
	list, err := mysql.GetNodeSortInfoByIn(invoker.Db.WithContext(ctx), needSearchArr)
	if err != nil {
		return fmt.Errorf("ChangeSortByTargetGuid fail, err: %w", err)
	}
	currentInfo := list.FindByGuid(currentGuid)

	if currentInfo == nil {
		return fmt.Errorf("ChangeSortByTargetGuid find current info fail")
	}
	parentInfo := list.FindByGuid(parentGuid)
	if parentInfo == nil {
		return fmt.Errorf("ChangeSortByTargetGuid find current info fail")
	}

	return invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先把数据+1
		err = tx.Model(&mysql.Edge{}).Where("parent_guid = ? and parent_type = ? and sort > ?", parentGuid, parentInfo.FileType.Number(), 0).
			Update("sort", gorm.Expr("`sort` + 1")).Error
		if err != nil {
			return fmt.Errorf("ChangeSortByTargetGuid update sort fail, err: %w", err)
		}
		// 更新排序
		// 并且将其放在该树形结构下面
		err = tx.Model(&mysql.Edge{}).Where("file_guid = ?", currentGuid).
			Updates(map[string]any{"sort": 1, "parent_guid": parentGuid, "parent_type": parentInfo.FileType.Number()}).Error
		if err != nil {
			return fmt.Errorf("ChangeSortByTargetGuid update sort fail2, err: %w", err)
		}

		if err = mysql.FileUpdate(tx, currentGuid, map[string]any{
			"parent_guid": parentGuid,
		}); err != nil {
			return errcodev1.ErrDbError().WithMessage("FileUpdate fail, err: " + err.Error())
		}

		err = cache.FileCache.SetParentGuid(ctx, tx, currentGuid, parentGuid)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("SetParentGuid fail, err: " + err.Error())
		}
		return nil
	})
}
