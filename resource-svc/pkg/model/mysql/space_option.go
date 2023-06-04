package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/resource-svc/pkg/service/spaceoption"

	"github.com/gotomicro/ego/core/elog"
	"gorm.io/gorm"
)

// SpaceOption 空间选项
type SpaceOption struct {
	Id          int64  `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	Guid        string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';"`
	OptionName  string `gorm:"not null; default:; comment:名称"`
	OptionValue int64  `gorm:"not null; default:0; comment:数据"`
	CreatedBy   int64  `gorm:"not null; default:0; comment:操作人"`
	Ctime       int64  `gorm:"not null; default:0; comment:创建时间"`
}

func (SpaceOption) TableName() string {
	return "space_option"
}

type SpaceOptions []SpaceOption

func (options SpaceOptions) ToPb(spaceType commonv1.CMN_APP) []*commonv1.SpaceOption {
	// 根据space type，获取space的配置信息
	optionList, err := spaceoption.GetListBySpaceType(spaceType)
	if err != nil {
		elog.Error("space options fail", elog.FieldErr(err))
		return nil
	}
	calcOptions := make([]*commonv1.SpaceOption, 0)
	for _, value := range optionList {
		outputOption := &commonv1.SpaceOption{
			Name:            value.Name,
			SpaceOptionId:   value.Option,
			SpaceOptionType: value.Type,
		}
		for _, optionValue := range options {
			if optionValue.OptionName == value.Option.String() {
				outputOption.Value = optionValue.OptionValue
			}
		}
		calcOptions = append(calcOptions, outputOption)
	}
	return calcOptions
}

// PutSpaceOption 增加管理员和超级管理员权限， 低频率操作，所以用写扩散方式，将权限分配。
func PutSpaceOption(db *gorm.DB, guid string, list []SpaceOption) (err error) {
	err = db.Where("guid = ?", guid).Delete(SpaceOption{}).Error
	if err != nil {
		err = fmt.Errorf("PutSpaceOption fail, err: %w", err)
		return
	}

	err = db.CreateInBatches(list, len(list)).Error
	if err != nil {
		err = fmt.Errorf("PutRolePolicy fail, err: %w", err)
		return
	}
	return
}

// GetSpaceOptionList 增加管理员和超级管理员权限，低频率操作，所以用写扩散方式，将权限分配。
func GetSpaceOptionList(db *gorm.DB, guid string) (list SpaceOptions, err error) {
	err = db.Select("option_name,option_value").Where("guid = ?", guid).Find(&list).Error
	if err != nil {
		err = fmt.Errorf("GetSpaceOptionList fail, err: %w", err)
		return
	}
	return
}

// BatchGetSpaceOptionList 增加管理员和超级管理员权限，低频率操作，所以用写扩散方式，将权限分配。
func BatchGetSpaceOptionList(db *gorm.DB, guids []string) (list SpaceOptions, err error) {
	err = db.Select("guid,option_name,option_value").Where("guid in (?)", guids).Find(&list).Error
	if err != nil {
		err = fmt.Errorf("GetSpaceOptionList fail, err: %w", err)
		return
	}
	return
}

// GetSpaceOptionInfo 增加管理员和超级管理员权限，低频率操作，所以用写扩散方式，将权限分配。
func GetSpaceOptionInfo(db *gorm.DB, guid string, optionName commonv1.SPC_OPTION) (list SpaceOption, err error) {
	err = db.Select("option_name,option_value").Where("guid = ? and  option_name = ?", guid, optionName.String()).Find(&list).Error
	if err != nil {
		err = fmt.Errorf("GetSpaceOptionList fail, err: %w", err)
		return
	}
	return
}
