package mysql

import (
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	pmsv1 "ecodepost/pb/pms/v1"

	"gorm.io/gorm"
)

// PmsPolicy
// 权限有三个维度
// 1. 人员： 用户或者用户组
// 2. 行动： 可以编辑，可以查看，可以导出，允许邀请，写在代码里
// 3. 资源： 作用于空间，空间分组
// subject: role id, uid
// action: CAN_EDIT
// resource: space, space_group, system
// subject，action，resource
type PmsPolicy struct {
	Id           int64            `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	SubjectId    int64            `gorm:"not null; default:0; comment:subjectID"`
	SubjectType  commonv1.PMS_SUB `gorm:"not null; default:0; comment:类型"`
	ActionName   string           `gorm:"not null; default:; comment:行为名称"`
	ResourceGuid string           `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';comment:资源GUID"`
	ResourceType commonv1.PMS_RSC `gorm:"not null; default:0; comment:资源类型"`
	CreatedBy    int64            `gorm:"not null; default:0; comment:操作人"`
	Ctime        int64            `gorm:"not null; default:0; comment:创建时间"`
}

func (PmsPolicy) TableName() string {
	return "pms_policy"
}

type PmsPolicies []PmsPolicy

// CheckRolePolicy 检测权限，目前只有role角色，可以check 权限
func CheckRolePolicy(db *gorm.DB, roleIds []int64, actionName commonv1.PMS_ACT) (flag bool, err error) {
	var cnt int64
	err = db.Model(PmsPolicy{}).Where("subject_id in (?) and subject_type = ?  and action_name = ? and resource_type = ?", roleIds, commonv1.PMS_SUB_ROLE.Number(), actionName.String(), commonv1.PMS_RSC_SYSTEM.Number()).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("check policy fail, err: %w", err)
		return
	}
	// 说明有这个策略，那么可以有这个权限
	if cnt > 0 {
		flag = true
	}
	return
}

// CheckPolicy 检测权限，目前只有role角色，可以check 权限
func CheckPolicy(db *gorm.DB, roleIds []int64, actionName commonv1.PMS_ACT, resourceGuid string, resourceType commonv1.PMS_RSC) (flag bool, err error) {
	var cnt int64
	err = db.Model(PmsPolicy{}).Where("subject_id in (?) and subject_type = ? and action_name = ? and resource_guid = ? and resource_type = ?", roleIds, commonv1.PMS_SUB_ROLE.Number(), actionName.String(), resourceGuid, resourceType.Number()).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("check policy fail, err: %w", err)
		return
	}
	// 说明有这个策略，那么可以有这个权限
	if cnt > 0 {
		flag = true
	}
	return
}

// PutRolePolicy 增加管理员和超级管理员权限， 低频率操作，所以用写扩散方式，将权限分配。
func PutRolePolicy(db *gorm.DB, putItems *pmsv1.PutRolePermissionReq) (oldRolePolicyList []PmsPolicy, newRolePolicyList []PmsPolicy, err error) {
	err = db.Where("subject_id = ? AND subject_type = ? AND resource_type = ?", putItems.GetRoleId(), commonv1.PMS_SUB_ROLE.Number(), commonv1.PMS_RSC_SYSTEM.Number()).Find(&oldRolePolicyList).Error
	if err != nil {
		err = fmt.Errorf("put role policy fail, err: %w", err)
		return
	}

	err = db.Where("subject_id = ? AND subject_type = ?  AND resource_type = ?", putItems.GetRoleId(), commonv1.PMS_SUB_ROLE.Number(), commonv1.PMS_RSC_SYSTEM.Number()).Delete(PmsPolicy{}).Error
	if err != nil {
		err = fmt.Errorf("put role policy fail, err: %w", err)
		return
	}

	nowTime := time.Now().Unix()
	// role的自身权限
	newRolePolicyList = make([]PmsPolicy, 0)
	for _, value := range putItems.GetList() {
		if value.Flag == 0 {
			continue
		}
		newRolePolicyList = append(newRolePolicyList, PmsPolicy{
			SubjectId:    putItems.GetRoleId(),
			SubjectType:  commonv1.PMS_SUB_ROLE,
			ActionName:   value.GetActionName(),
			ResourceGuid: "",
			ResourceType: commonv1.PMS_RSC_SYSTEM,
			CreatedBy:    putItems.GetOperateUid(),
			Ctime:        nowTime,
		})
	}

	err = db.CreateInBatches(newRolePolicyList, len(newRolePolicyList)).Error
	if err != nil {
		err = fmt.Errorf("PutRolePolicy fail, err: %w", err)
		return
	}
	return
}

// PutRoleSpaceGroupPolicy 增加管理员和超级管理员权限， 低频率操作，所以用写扩散方式，将权限分配。
func PutRoleSpaceGroupPolicy(db *gorm.DB, putItems *pmsv1.PutRoleSpaceGroupPermissionReq) (oldRolePolicyList []PmsPolicy, newRolePolicyList []PmsPolicy, err error) {
	err = db.Where("subject_id = ? and subject_type = ? and resource_type = ?", putItems.GetRoleId(), commonv1.PMS_SUB_ROLE.Number(), commonv1.PMS_RSC_SPACE_GROUP.Number()).Find(&oldRolePolicyList).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpaceGroupPolicy fail1, err: %w", err)
		return
	}

	err = db.Where("subject_id = ? and subject_type = ? and resource_type = ?", putItems.GetRoleId(), commonv1.PMS_SUB_ROLE.Number(), commonv1.PMS_RSC_SPACE_GROUP.Number()).Delete(PmsPolicy{}).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpaceGroupPolicy fail2, err: %w", err)
		return
	}

	// 删除role和space group关系
	err = db.Where("role_id = ?  and guid_type = ?", putItems.GetRoleId(), commonv1.CMN_GUID_SPACE_GROUP.Number()).Delete(PmsRoleSpace{}).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpaceGroupPolicy fail3, err: %w", err)
		return
	}

	nowTime := time.Now().Unix()

	// role 和 space group关系
	spaceGroupRoleList := make([]PmsRoleSpace, 0)

	// 加入space group权限信息
	newRolePolicyList = make([]PmsPolicy, 0)
	// todo 未校验是不是自己管理的guid
	for _, spaceGroupInfo := range putItems.GetList() {
		spaceGroupRoleList = append(spaceGroupRoleList, PmsRoleSpace{
			RoleId:   putItems.GetRoleId(),
			Guid:     spaceGroupInfo.GetGuid(),
			GuidType: commonv1.CMN_GUID_SPACE_GROUP,
			Ctime:    time.Now().Unix(),
		})

		for _, value := range spaceGroupInfo.GetList() {
			if value.Flag == 0 {
				continue
			}
			newRolePolicyList = append(newRolePolicyList, PmsPolicy{
				SubjectId:    putItems.GetRoleId(),
				SubjectType:  commonv1.PMS_SUB_ROLE,
				ActionName:   value.GetActionName(),
				ResourceGuid: spaceGroupInfo.GetGuid(),
				ResourceType: commonv1.PMS_RSC_SPACE_GROUP,
				CreatedBy:    putItems.GetOperateUid(),
				Ctime:        nowTime,
			})
		}
	}
	err = db.CreateInBatches(spaceGroupRoleList, len(spaceGroupRoleList)).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpaceGroupPolicy fail4, err: %w", err)
		return
	}

	err = db.CreateInBatches(newRolePolicyList, len(newRolePolicyList)).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpaceGroupPolicy fail5, err: %w", err)
		return
	}
	return
}

// PutRoleSpacePolicy 增加管理员和超级管理员权限， 低频率操作，所以用写扩散方式，将权限分配。
func PutRoleSpacePolicy(db *gorm.DB, putItems *pmsv1.PutRoleSpacePermissionReq) (oldRolePolicyList []PmsPolicy, newRolePolicyList []PmsPolicy, err error) {
	err = db.Where("subject_id = ? and subject_type = ? and resource_type = ?", putItems.GetRoleId(), commonv1.PMS_SUB_ROLE.Number(), commonv1.PMS_RSC_SPACE.Number()).Find(&oldRolePolicyList).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpacePolicy fail1, err: %w", err)
		return
	}

	err = db.Where("subject_id = ? and subject_type = ?  and resource_type = ?", putItems.GetRoleId(), commonv1.PMS_SUB_ROLE.Number(), commonv1.PMS_RSC_SPACE.Number()).Delete(PmsPolicy{}).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpacePolicy fail2, err: %w", err)
		return
	}

	// 删除role和space group关系
	err = db.Where("role_id = ? and guid_type = ?", putItems.GetRoleId(), commonv1.CMN_GUID_SPACE.Number()).Delete(PmsRoleSpace{}).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpacePolicy fail3, err: %w", err)
		return
	}

	nowTime := time.Now().Unix()

	// role 和 space关系
	spaceRoleList := make([]PmsRoleSpace, 0)

	// 加入space权限信息
	newRolePolicyList = make([]PmsPolicy, 0)
	// todo 未校验是不是自己管理的guid
	for _, spaceGuidInfo := range putItems.GetList() {
		spaceRoleList = append(spaceRoleList, PmsRoleSpace{
			RoleId:   putItems.GetRoleId(),
			Guid:     spaceGuidInfo.GetGuid(),
			GuidType: commonv1.CMN_GUID_SPACE,
			Ctime:    time.Now().Unix(),
		})

		for _, value := range spaceGuidInfo.GetList() {
			if value.Flag == 0 {
				continue
			}
			newRolePolicyList = append(newRolePolicyList, PmsPolicy{
				SubjectId:    putItems.GetRoleId(),
				SubjectType:  commonv1.PMS_SUB_ROLE,
				ActionName:   value.GetActionName(),
				ResourceGuid: spaceGuidInfo.GetGuid(),
				ResourceType: commonv1.PMS_RSC_SPACE,
				CreatedBy:    putItems.GetOperateUid(),
				Ctime:        nowTime,
			})
		}
	}

	err = db.CreateInBatches(spaceRoleList, len(spaceRoleList)).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpacePolicy fail4, err: %w", err)
		return
	}

	err = db.CreateInBatches(newRolePolicyList, len(newRolePolicyList)).Error
	if err != nil {
		err = fmt.Errorf("PutRoleSpacePolicy fail5, err: %w", err)
		return
	}

	return
}

func PolicyMap(db *gorm.DB, cmtGuid string, subjectId int64) (policyMap map[string]PmsPolicy, err error) {
	var list PmsPolicies
	err = db.Where("subject_id = ? and subject_type = ? and cmt_guid = ? ", subjectId, commonv1.PMS_SUB_ROLE.Number(), cmtGuid).Find(&list).Error
	if err != nil {
		err = fmt.Errorf("PolicyMap fail, err: %w", err)
		return
	}

	policyMap = make(map[string]PmsPolicy)
	for _, policy := range list {
		policyMap[policy.ActionName] = policy
	}
	return
}

// PolicyPmsList 如果使用了角色，那么会拿到space，space group guid的信息，所以需要在进一步分类，做筛选信息
func PolicyPmsList(db *gorm.DB, subjectId int64) (rolePolicies PmsPolicies, spaceGroupPolicies PmsPolicies, spacePolicies PmsPolicies, err error) {
	var list PmsPolicies
	err = db.Where("subject_id = ? and subject_type = ? ", subjectId, commonv1.PMS_SUB_ROLE.Number()).Find(&list).Error
	if err != nil {
		err = fmt.Errorf("PolicyMap fail, err: %w", err)
		return
	}

	rolePolicies = make(PmsPolicies, 0)
	spacePolicies = make(PmsPolicies, 0)
	spaceGroupPolicies = make(PmsPolicies, 0)

	for _, value := range list {
		if value.ResourceType == commonv1.PMS_RSC_SYSTEM {
			rolePolicies = append(rolePolicies, value)
		}
		if value.ResourceType == commonv1.PMS_RSC_SPACE_GROUP {
			spaceGroupPolicies = append(spaceGroupPolicies, value)
		}
		if value.ResourceType == commonv1.PMS_RSC_SPACE {
			spacePolicies = append(spacePolicies, value)
		}
	}
	return
}
