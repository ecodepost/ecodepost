package job

import (
	"fmt"

	resmysql "ecodepost/resource-svc/pkg/model/mysql"
	"github.com/ego-component/egorm"

	"ecodepost/user-svc/pkg/model/mysql"

	"github.com/ego-component/eoauth2/storage/dao"
	"github.com/gotomicro/ego/task/ejob"
)

func RunInstall(ctx ejob.Context) error {
	userDb := egorm.Load("mysql").Build()
	models := []interface{}{
		&dao.App{},
		&mysql.CountBehaviorLog{},
		&mysql.CountLog{},
		&mysql.CountStat{},
		&mysql.User{},
		&mysql.UserCollectionGroup{},
		&mysql.UserCollection{},
		&mysql.Logger{},
		&mysql.System{},
		&mysql.UserOpen{},
		&mysql.NotifyLetter{},
		&mysql.NotifyTplChannel{},
		&mysql.NotifyTpl{},
		&mysql.Notify{},
		&mysql.NotifySign{},
		&mysql.NotifyLetterUserTs{},
	}
	err := userDb.Debug().WithContext(ctx.Ctx).Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models...)
	if err != nil {
		return err
	}
	fmt.Println("create user table ok")

	models = []interface{}{
		&resmysql.Image{},
		&resmysql.File{},
		&resmysql.FileRecommend{},
		&resmysql.FileSpaceTop{},
		&resmysql.FileEmoji{},
		&resmysql.FileEmojiStatics{},
		&resmysql.Edge{},
		&resmysql.Space{},
		&resmysql.SpaceGroup{},
		&resmysql.SpaceColumn{},
		&resmysql.SpaceOption{},
		&resmysql.CommentSubject{},
		&resmysql.CommentIndex{},
		&resmysql.CommentContent{},
		&resmysql.Resource{},
		&resmysql.SpaceGroupMember{},
		&resmysql.SpaceMember{},
		&resmysql.PmsRoleSpace{},
		&resmysql.PmsPolicy{},
		&resmysql.PmsRole{},
		&resmysql.PmsRoleMember{},
		&resmysql.PmsSuperAdminMember{},
		&resmysql.AuditIndex{},
		&resmysql.AuditLog{},
	}
	mainDb := egorm.Load("mysql").Build()

	err = mainDb.Debug().WithContext(ctx.Ctx).Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models...)
	if err != nil {
		return err
	}
	fmt.Println("create main table ok")
	return nil
}
