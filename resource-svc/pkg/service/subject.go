package service

import (
	"fmt"

	"ecodepost/resource-svc/pkg/constx"
	"ecodepost/resource-svc/pkg/model/dao"

	commonv1 "ecodepost/pb/common/v1"

	"gorm.io/gorm"
)

var SubjectNotFound = fmt.Errorf("SubjectNotFound")

type subjectService struct{}

func (subjectService) getKeyGuid(bizGuid string, bizType commonv1.CMN_BIZ) string {
	return constx.CommentSubjectCache + fmt.Sprintf("%s_%d", bizGuid, int32(bizType))
}

func (s subjectService) GetInfoByMySQLByGuid(db *gorm.DB, bizGuid string, bizType commonv1.CMN_BIZ) (int64, int32, error) {
	reply, err := dao.CommentSubject.InfoByBizInfo(db, bizGuid, bizType)
	if err != nil {
		return 0, 0, err
	}

	if reply.Id == 0 {
		return 0, 0, SubjectNotFound
	}
	return reply.Id, reply.CntComment, nil
}
