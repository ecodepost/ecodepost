package service

import (
	"context"
	"fmt"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	"github.com/samber/lo"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	spacev1 "ecodepost/pb/space/v1"
	// trackv1 "ecodepost/pb/track/v1"
	userv1 "ecodepost/pb/user/v1"

	"gorm.io/gorm"
)

type space struct{}

func InitSpace() *space {
	return &space{}
}

func (*space) List(ctx context.Context, db *gorm.DB, uid int64) (pbSpaceList []*spacev1.TreeSpace, pbSpaceGroupList []*spacev1.TreeSpaceGroup, err error) {
	var isAllowCreateSpaceBool, isAllowUpdateSpace, isAllowUpdateSpaceGroup bool
	if uid != 0 {
		if isAllowCreateSpaceBool, err = pmspolicy.Check(ctx, commonv1.PMS_SPACE_CREATE, uid, ""); err != nil {
			return nil, nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}
		if isAllowUpdateSpace, err = pmspolicy.Check(ctx, commonv1.PMS_SPACE_UPDATE, uid, ""); err != nil {
			return nil, nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}
		if isAllowUpdateSpaceGroup, err = pmspolicy.Check(ctx, commonv1.PMS_SPACE_GROUP_UPDATE, uid, ""); err != nil {
			return nil, nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}
	}

	spaceList, spaceGroups, err := mysql.SpaceAndGroupList(db, uid)
	pbSpaceList = lo.Map(spaceList, func(t *mysql.Space, i int) *spacev1.TreeSpace {
		pbInfo := t.ToTreePb()
		pbInfo.IsAllowSet = isAllowUpdateSpace
		pbInfo.EmojiList = mysql.EmojiList()
		return pbInfo
	})
	pbSpaceGroupList = lo.Map(spaceGroups, func(t *mysql.SpaceGroup, i int) *spacev1.TreeSpaceGroup {
		pbInfo := t.ToTreePb()
		pbInfo.IsAllowCreateSpace = isAllowCreateSpaceBool
		pbInfo.IsAllowSet = isAllowUpdateSpaceGroup
		return pbInfo
	})
	return
}

// Tree
// Deprecated: 待移除
//func (*space) Tree(ctx context.Context, db *gorm.DB, cmtGuid string, uid int64) (output []*spacev1.AntSpaceGroupInfo, err error) {
//	var isAllowCreateSpace, IsAllowSet, isAllowSetSpace int32
//	var isAllowCreateSpaceBool, isAllowUpdateSpace, isAllowUpdateSpaceGroup bool
//	if uid != 0 {
//		isAllowCreateSpaceBool, err = pmspolicy.Check(ctx, commonv1.PMS_SPACE_CREATE, uid, "")
//		if err != nil {
//			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
//		}
//		isAllowUpdateSpace, err = pmspolicy.Check(ctx, commonv1.PMS_SPACE_UPDATE, uid, "")
//		if err != nil {
//			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
//		}
//		isAllowUpdateSpaceGroup, err = pmspolicy.Check(ctx, commonv1.PMS_SPACE_GROUP_UPDATE, uid, "")
//		if err != nil {
//			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
//		}
//	}
//	if isAllowCreateSpaceBool {
//		isAllowCreateSpace = 1
//	}
//	if isAllowUpdateSpace {
//		isAllowSetSpace = 1
//	}
//	if isAllowUpdateSpaceGroup {
//		IsAllowSet = 1
//	}
//
//	spaceGroupList, err := mysql.SpaceTree(db, cmtGuid, uid)
//	output = make([]*spacev1.AntSpaceGroupInfo, 0)
//	for _, spaceGroupInfo := range spaceGroupList {
//		spaceList := make([]*spacev1.AntSpaceInfo, 0)
//		for _, spaceInfo := range spaceGroupInfo.List {
//			spaceList = append(spaceList, &spacev1.AntSpaceInfo{
//				Guid:         spaceInfo.Guid,
//				Name:         spaceInfo.Name,
//				CmtGuid:      spaceInfo.CmtGuid,
//				IconType:     spaceInfo.IconType,
//				Icon:         spaceInfo.Icon,
//				SpaceType:    spaceInfo.Type,
//				SpaceLayout:  spaceInfo.Layout,
//				Visibility:   spaceInfo.Visibility,
//				IsAllowSet:   isAllowSetSpace,
//				Access:       spaceInfo.Access,
//				SpaceOptions: spaceInfo.OptionList.ToPb(spaceInfo.Type),
//				HeadImage:    spaceInfo.HeadImage,
//				Link:         spaceInfo.Link,
//			})
//		}
//		output = append(output, &spacev1.AntSpaceGroupInfo{
//			Guid:               spaceGroupInfo.Guid,
//			Name:               spaceGroupInfo.Name,
//			CmtGuid:            spaceGroupInfo.CmtGuid,
//			IconType:           spaceGroupInfo.IconType,
//			Icon:               spaceGroupInfo.Icon,
//			Children:           spaceList,
//			Visibility:         spaceGroupInfo.Visibility,
//			IsAllowSet:         IsAllowSet,
//			IsAllowCreateSpace: isAllowCreateSpace,
//		})
//	}
//	return
//}

func (s *space) CreateGroup(ctx context.Context, req *spacev1.CreateSpaceGroupReq) (resp *spacev1.SpaceGroupInfo, err error) {
	guid, err := Resource.GenerateGuid(ctx, commonv1.CMN_GUID_SPACE_GROUP, req.GetOperateUid())
	if err != nil {
		return nil, fmt.Errorf("generate guid fail, err: %w", err)
	}

	userInfo, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{Uid: req.GetOperateUid()})
	if err != nil {
		return nil, fmt.Errorf("userinfo failed, err: %w", err)
	}
	nowTime := time.Now().Unix()
	data := &mysql.SpaceGroup{
		Guid:                  guid,
		Name:                  req.GetName(),
		Ctime:                 nowTime,
		Utime:                 nowTime,
		CreatedBy:             req.GetOperateUid(),
		UpdatedBy:             req.GetOperateUid(),
		IconType:              req.GetIconType(),
		Icon:                  req.GetIcon(),
		Sort:                  time.Now().UnixMilli(),
		Visibility:            req.GetVisibility(),
		IsAllowReadMemberList: req.GetIsAllowReadMemberList(),
	}
	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = mysql.SpaceGroupCreate(tx, data); err != nil {
			return fmt.Errorf("space group create fail, err: %w", err)
		}
		err = mysql.SpaceGroupMemberCreate(tx, &mysql.SpaceGroupMember{
			Ctime:     nowTime,
			Utime:     nowTime,
			Uid:       data.CreatedBy,
			Nickname:  userInfo.User.GetNickname(),
			Guid:      data.Guid,
			CreatedBy: data.CreatedBy,
		})
		if err != nil {
			return fmt.Errorf("space group member create fail, err: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	resp = data.ToPb()
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:          commonv1.LOG_EVENT_SPACE_GROUP_CREATE,
		Group:          commonv1.LOG_GROUP_SPACE_GROUP,
		OperateUid:     req.GetOperateUid(),
		SpaceGroupGuid: data.Guid,
	})
	return
}

func (*space) Count(db *gorm.DB) (cnt int64, err error) {
	err = db.Model(mysql.Space{}).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("space count fail, err: %w", err)
		return
	}
	return
}

func SpaceOptionAddFileAndCommentAndSortPms(data *mysql.Space, uid int64, nowTime int64) []mysql.SpaceOption {
	return []mysql.SpaceOption{
		{
			Guid:        data.Guid,
			OptionName:  commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_FILE.String(),
			OptionValue: 1,
			CreatedBy:   uid,
			Ctime:       nowTime,
		},
		{
			Guid:        data.Guid,
			OptionName:  commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_COMMENT.String(),
			OptionValue: 1,
			CreatedBy:   uid,
			Ctime:       nowTime,
		},
		{
			Guid:        data.Guid,
			OptionName:  commonv1.SPC_OPTION_FILE_DEFAULT_SORT.String(),
			OptionValue: int64(commonv1.CMN_SORT_CREATE_TIME),
			CreatedBy:   uid,
			Ctime:       nowTime,
		},
	}
}

func SpaceOptionAddFileAndCommentPms(data *mysql.Space, uid int64, nowTime int64) []mysql.SpaceOption {
	return []mysql.SpaceOption{
		{
			Guid:        data.Guid,
			OptionName:  commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_FILE.String(),
			OptionValue: 1,
			CreatedBy:   uid,
			Ctime:       nowTime,
		},
		{
			Guid:        data.Guid,
			OptionName:  commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_COMMENT.String(),
			OptionValue: 1,
			CreatedBy:   uid,
			Ctime:       nowTime,
		},
	}
}

func (s *space) Create(ctx context.Context, req *spacev1.CreateSpaceReq) (resp *commonv1.SpaceInfo, err error) {
	tx := invoker.Db.WithContext(ctx).Begin()
	guid, err := Resource.GenerateGuid(ctx, commonv1.CMN_GUID_SPACE, req.GetOperateUid())
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("generate guid fail, err: %w", err)
	}

	icon := "#"
	if req.Icon != "" {
		icon = req.Icon
	}

	// 不同可见类型，初始化access不一样
	access := commonv1.SPC_ACS_INVALID
	if req.GetVisibility() == commonv1.CMN_VISBL_INTERNAL {
		access = commonv1.SPC_ACS_OPEN
	}
	if req.GetVisibility() == commonv1.CMN_VISBL_SECRET {
		access = commonv1.SPC_ACS_DENY_ALL
	}

	nowTime := time.Now().Unix()
	data := &mysql.Space{
		SpaceGroupGuid:        req.GetSpaceGroupGuid(),
		Guid:                  guid,
		Name:                  req.GetName(),
		Ctime:                 nowTime,
		Utime:                 nowTime,
		CreatedBy:             req.GetOperateUid(),
		UpdatedBy:             req.GetOperateUid(),
		IconType:              commonv1.FILE_IT_EMOJI,
		Icon:                  icon,
		Sort:                  time.Now().UnixMilli(),
		Visibility:            req.GetVisibility(),
		Type:                  req.GetSpaceType(),
		Layout:                req.GetSpaceLayout(),
		Access:                access,
		IsAllowReadMemberList: false,
		ChargeType:            commonv1.SPC_CT_FREE,
		Link:                  req.GetLink(),
		Cover:                 req.GetCover(),
	}

	// 查询user
	userInfo, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{Uid: req.GetOperateUid()})
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("userinfo failed, err: %w", err)
	}

	// 查询spaceGroup
	if _, err = mysql.SpaceGroupGetId(tx, data.SpaceGroupGuid); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 根据空间类型，调整其他字段配置
	var opts []mysql.SpaceOption
	switch data.Type {
	case commonv1.CMN_APP_ARTICLE:
		opts = SpaceOptionAddFileAndCommentAndSortPms(data, req.GetOperateUid(), nowTime)
	case commonv1.CMN_APP_QA:
		opts = SpaceOptionAddFileAndCommentAndSortPms(data, req.GetOperateUid(), nowTime)
	case commonv1.CMN_APP_COLUMN:
		opts = SpaceOptionAddFileAndCommentPms(data, req.GetOperateUid(), nowTime)

		// data.Status = commonv1.SPC_STATUS_DRAFT
		data.Visibility = commonv1.CMN_VISBL_DRAFT
	}

	// 插入到space
	if err = mysql.SpaceCreate(tx, data); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("space group create fail, err: %w", err)
	}

	// 如果有options配置，则插入到space_option
	if len(opts) != 0 {
		err = mysql.PutSpaceOption(tx, data.Guid, opts)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 更新spaceMember表
	err = s.CreateMember(ctx, tx, &mysql.SpaceMember{
		Ctime:     nowTime,
		Utime:     nowTime,
		Uid:       data.CreatedBy,
		Nickname:  userInfo.User.GetNickname(),
		Guid:      data.Guid,
		CreatedBy: data.CreatedBy,
	})
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("space member create fail, err: %w", err)
	}
	tx.Commit()

	return data.ToPb(), nil
}

func (*space) DeleteGroup(db *gorm.DB, uid int64, guid string) (err error) {
	err = db.Model(mysql.SpaceGroup{}).Where("guid = ?", guid).Updates(map[string]any{
		"deleted_by": uid,
		"dtime":      time.Now().Unix(),
	}).Error
	if err != nil {
		return fmt.Errorf("space group delete fail, err: %w", err)
	}
	return nil
}

func (*space) Delete(db *gorm.DB, uid int64, guid string) (err error) {
	err = db.Model(mysql.Space{}).Where("guid = ?", guid).Updates(map[string]any{
		"deleted_by": uid,
		"dtime":      time.Now().Unix(),
	}).Error
	if err != nil {
		return fmt.Errorf("space delete fail, err: %w", err)
	}
	return nil
}

func (*space) UpdateGroup(db *gorm.DB, uid int64, guid string, updates map[string]any) (err error) {
	id, err := mysql.SpaceGroupGetId(db, guid)
	if err != nil {
		return err
	}
	updates["utime"] = time.Now().Unix()
	updates["updated_by"] = uid

	err = db.Model(mysql.SpaceGroup{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("update group fail, err: %w", err)
	}
	return nil
}

func (*space) Update(db *gorm.DB, uid int64, guid string, updates map[string]any) (err error) {
	id, err := mysql.SpaceGetId(db, guid)
	if err != nil {
		return fmt.Errorf("update group fail, err: %w", err)
	}
	if id == 0 {
		return fmt.Errorf("update group get space group  fail")
	}
	updates["utime"] = time.Now().Unix()
	updates["updated_by"] = uid

	err = db.Model(mysql.Space{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("update fail, err: %w", err)
	}
	return nil
}

func (*space) ChangeGroupSort(ctx context.Context, currentGuid string, targetGuid string, dropPosition string) (err error) {
	if targetGuid == "" {
		return fmt.Errorf("targetGuid is empty")
	}

	var targetSort int64
	needSearchArr := []string{currentGuid, targetGuid}
	list, err := mysql.SpaceGroupSortGetInfoByIn(invoker.Db.WithContext(ctx), needSearchArr)
	if err != nil {
		return fmt.Errorf("ChangeGroupSort fail, err: %w", err)
	}

	currentInfo := list.FindByGuid(currentGuid)
	if currentInfo == nil {
		return fmt.Errorf("ChangeGroupSort find current info fail")
	}

	// 如果存在该数据
	if targetGuid != "" {
		afterInfo := list.FindByGuid(targetGuid)
		if currentInfo == nil {
			return fmt.Errorf("ChangeGroupSort find after info fail")
		}
		targetSort = afterInfo.Sort
	}

	var sign = ""
	var judgeSign = ""
	var updateSort = targetSort

	if dropPosition == "before" {
		sign = "-"
		judgeSign = "<"
		updateSort--
	} else if dropPosition == "after" {
		sign = "+"
		judgeSign = ">"
		updateSort++
	} else {
		return fmt.Errorf("dropPosition is invalid")
	}

	tx := invoker.Db.WithContext(ctx).Begin()
	// 先把数据+1
	err = tx.Model(&mysql.SpaceGroup{}).Where("sort "+judgeSign+" ?", targetSort).Update("sort", gorm.Expr("`sort` "+sign+" 1")).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ChangeGroupSort update sort fail, err: %w", err)
	}

	err = tx.Model(&mysql.SpaceGroup{}).Where(" guid = ?", currentGuid).Update("sort", updateSort).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ChangeGroupSort update sort fail2, err: %w", err)
	}
	tx.Commit()
	return nil
}

// ChangeSpaceSort 更改顺序
// 顺序是按照sort asc，越小越往前
// 1 放到某个file guid后面
// 2 放到某个parent guid下面
func (s *space) ChangeSpaceSort(ctx context.Context, currentGuid string, targetGuid *string, dropPosition *string, parentGuid *string) (err error) {
	if currentGuid == "" {
		return fmt.Errorf("current guid is empty")
	}
	if targetGuid != nil {
		if currentGuid == *targetGuid {
			return fmt.Errorf("current guid cant eq after guid")
		}
	}
	if parentGuid != nil {
		if currentGuid == *parentGuid {
			return fmt.Errorf("current guid cant eq parent guid")
		}
	}

	if targetGuid != nil {
		return s.ChangeSortByTargetGuid(ctx, currentGuid, targetGuid, dropPosition)
	}

	if *parentGuid == "" {
		return fmt.Errorf("parent guid is empty")
	}

	return s.ChangeSortByParentGuid(ctx, currentGuid, *parentGuid)
}

func (*space) ChangeSortByTargetGuid(ctx context.Context, currentGuid string, targetGuid *string, dropPosition *string) (err error) {
	var targetSort int64
	needSearchArr := []string{currentGuid}
	needSearchArr = append(needSearchArr, *targetGuid)

	commonDb := invoker.Db.WithContext(ctx)
	list, err := mysql.SpaceSortGetInfoByIn(commonDb, needSearchArr)
	if err != nil {
		return fmt.Errorf("ChangeSortByTargetGuid fail, err: %w", err)
	}

	if currentInfo := list.FindByGuid(currentGuid); currentInfo == nil {
		return fmt.Errorf("ChangeSortByTargetGuid find current info fail")
	}

	// 如果存在该数据
	targetInfo := list.FindByGuid(*targetGuid)
	if targetInfo == nil {
		return fmt.Errorf("ChangeSortByTargetGuid find after info fail")
	}
	targetSort = targetInfo.Sort

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
	tx := invoker.Db.WithContext(ctx).Begin()

	// 先把数据+1
	err = tx.Model(&mysql.Space{}).Where(" sort "+judgeSign+" ?", targetSort).Update("sort", gorm.Expr("`sort` "+sign+" 1")).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ChangeSortByTargetGuid update sort fail, err: %w", err)
	}
	err = tx.Model(&mysql.Space{}).Where("guid = ?", currentGuid).Updates(map[string]any{
		"sort":             updateSort,
		"space_group_guid": targetInfo.SpaceGroupGuid,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ChangeSortByTargetGuid update sort fail2, err: %w", err)
	}
	tx.Commit()
	return nil
}

func (*space) ChangeSortByParentGuid(ctx context.Context, currentGuid string, parentGuid string) (err error) {
	needSearchArr := []string{currentGuid}
	needSearchArr = append(needSearchArr)
	commonDb := invoker.Db.WithContext(ctx)
	list, err := mysql.SpaceSortGetInfoByIn(commonDb, needSearchArr)
	if err != nil {
		return fmt.Errorf("ChangeSortByTargetGuid fail, err: %w", err)
	}
	currentInfo := list.FindByGuid(currentGuid)
	if currentInfo == nil {
		return fmt.Errorf("ChangeSortByTargetGuid find current info fail")
	}

	return invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先把数据+1
		err = tx.Model(&mysql.Space{}).Where("sort > ?", 0).
			Update("sort", gorm.Expr("`sort` + 1")).Error
		if err != nil {
			return fmt.Errorf("ChangeSortByTargetGuid update sort fail, err: %w", err)
		}
		// 更新排序
		// 并且将其放在该树形结构下面
		err = tx.Model(&mysql.Space{}).Where("guid = ?", currentGuid).
			Updates(map[string]any{
				"sort":             1,
				"space_group_guid": parentGuid,
			}).Error
		if err != nil {
			return fmt.Errorf("ChangeSortByTargetGuid update sort fail2, err: %w", err)
		}
		return nil
	})
}

func (*space) CreateMember(ctx context.Context, db *gorm.DB, data *mysql.SpaceMember) (err error) {
	// invoker.GrpcCommunity.Join()
	if err = mysql.SpaceMemberCreate(db, data); err != nil {
		return fmt.Errorf("create member fail, err: %w", err)
	}
	// _, err = invoker.GrpcTrack.Total(ctx, &trackv1.TotalReq{
	// 	Event:     commonv1.TRACK_TOTAL_EVENT_SPACE_CREATE_MEMBER,
	// 	Tid:       etrace.ExtractTraceID(ctx),
	// 	Uid:       data.Uid,
	// 	CmtGuid:   data.CmtGuid,
	// 	SpaceGuid: data.Guid,
	// 	Ts:        time.Now().UnixMilli(),
	// })
	// if err != nil {
	// 	return fmt.Errorf("create member fail2, err: %w", err)
	// }
	return nil
}

func (*space) CreateBatchMember(ctx context.Context, db *gorm.DB, list mysql.SpaceMembers) (err error) {
	if err = mysql.SpaceMemberBatchCreate(db, list); err != nil {
		return fmt.Errorf("create member fail, err: %w", err)
	}

	// reqs := make([]*trackv1.TotalReq, 0)
	// for _, data := range list {
	// 	reqs = append(reqs, &trackv1.TotalReq{
	// 		Event:     commonv1.TRACK_TOTAL_EVENT_SPACE_CREATE_MEMBER,
	// 		Tid:       etrace.ExtractTraceID(ctx),
	// 		Uid:       data.Uid,
	// 		CmtGuid:   data.CmtGuid,
	// 		SpaceGuid: data.Guid,
	// 		Ts:        time.Now().UnixMilli(),
	// 	})
	// }
	// if _, err = invoker.GrpcTrack.BatchTotal(ctx, &trackv1.BatchTotalReq{Reqs: reqs}); err != nil {
	// 	return fmt.Errorf("create member fail2, err: %w", err)
	// }
	return nil
}

func (*space) DeleteBatchMember(ctx context.Context, db *gorm.DB, spaceGuid string, uids []int64, event commonv1.TRACK_TOTAL_EVENT) (err error) {
	if err = mysql.SpaceMemberBatchDelete(db, spaceGuid, uids); err != nil {
		return fmt.Errorf("create member fail, err: %w", err)
	}
	// nowMilliTime := time.Now().UnixMilli()
	// reqs := make([]*trackv1.TotalReq, 0)
	// for _, uid := range uids {
	// 	reqs = append(reqs, &trackv1.TotalReq{
	// 		Event:     event,
	// 		Uid:       uid,
	// 		CmtGuid:   cmtGuid,
	// 		SpaceGuid: spaceGuid,
	// 		Ts:        nowMilliTime,
	// 	})
	// }
	//
	// if _, err = invoker.GrpcTrack.BatchTotal(ctx, &trackv1.BatchTotalReq{Reqs: reqs}); err != nil {
	// 	return fmt.Errorf("create member fail2, err: %w", err)
	// }
	return nil
}
