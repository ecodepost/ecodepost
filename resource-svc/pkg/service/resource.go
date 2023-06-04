package service

import (
	"context"
	"fmt"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	"gorm.io/gorm"

	"github.com/ego-component/eguid"

	commonv1 "ecodepost/pb/common/v1"
)

type resource struct {
	suid      *eguid.Component
	spaceGuid *eguid.Component
}

func initResource() *resource {
	return &resource{
		spaceGuid: eguid.Load("resource-svc.spaceGuid").Build(),
		suid:      eguid.Load("resource-svc.guid").Build(),
	}
}

func (r *resource) GenerateGuid(ctx context.Context, guidType commonv1.CMN_GUID, uid int64) (string, error) {
	var guidObj *eguid.Component
	switch guidType {
	case commonv1.CMN_GUID_SPACE:
		guidObj = r.spaceGuid
	case commonv1.CMN_GUID_SPACE_GROUP:
		guidObj = r.spaceGuid
	case commonv1.CMN_GUID_FILE:
		guidObj = r.suid
	case commonv1.CMN_GUID_COMMENT:
		guidObj = r.suid
	default:
		return "", fmt.Errorf("not exist type: " + guidType.String())
	}
	var guid string
	err := invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := &mysql.Resource{
			GuidType:  guidType,
			Ctime:     time.Now().Unix(),
			CreatedBy: uid,
		}
		err := mysql.ResourceCreate(tx, data)
		if err != nil {
			return fmt.Errorf("GenerateGuid failed, err: %w", err)
		}
		if data.Id == 0 {
			return fmt.Errorf("GenerateGuid failed2, err: data id is 0")
		}

		guid, err = guidObj.EncodeRandomInt64(data.Id)
		if err != nil {
			return fmt.Errorf("GenerateGuid failed3, err: %w", err)
		}
		err = tx.Model(mysql.Resource{}).Where("id = ?", data.Id).Updates(map[string]any{
			"guid": guid,
		}).Error
		if err != nil {
			return fmt.Errorf("GenerateGuid failed4, err: %w", err)
		}
		return nil
	})

	return guid, err
}
