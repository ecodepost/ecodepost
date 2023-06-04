package job

import (
	"encoding/json"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	communityv1 "ecodepost/pb/community/v1"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"github.com/ego-component/eoauth2/storage/dao"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/task/ejob"
	"google.golang.org/protobuf/proto"
)

func init() {
	eflag.Register(&eflag.StringFlag{
		Name:    "redirectUri",
		Usage:   "--redirectUri",
		Default: "http://localhost:8000/api/oauth/code",
		Action:  func(string, *eflag.FlagSet) {},
	})
	eflag.Register(&eflag.StringFlag{
		Name:    "clientId",
		Usage:   "--clientId",
		Default: "",
		Action:  func(string, *eflag.FlagSet) {},
	})
}

func InitData(ctx ejob.Context) (err error) {
	InitSsoData(ctx)
	InitSystem(ctx)
	return nil
}

func InitSsoData(ctx ejob.Context) (err error) {
	unixtime := time.Now().Unix()
	appCreate := &dao.App{
		ClientId: "clientid123456",
		//ClientId:    random.String(16, random.Alphanumeric),
		Name:   "post-local",
		Secret: "secret123456",
		//Secret:      random.String(32, random.Alphanumeric),
		RedirectUri: eflag.String("redirectUri"),
		Extra:       "",
		Status:      1,
		Url:         "",
		Ctime:       unixtime,
		Utime:       unixtime,
	}
	err = invoker.TokenComponent.GetAPI().CreateClient(ctx.Ctx, appCreate)
	if err != nil {
		return
	}
	return nil
}

func InitSystem(ctx ejob.Context) (err error) {
	err = mysql.InitSystem(invoker.Db.WithContext(ctx.Ctx))
	if err != nil {
		return
	}

	err = mysql.PutSystemBase(invoker.Db.WithContext(ctx.Ctx), &communityv1.UpdateReq{
		Name:           proto.String("Ecode Post"),
		Description:    proto.String("Ecode Post Description"),
		Logo:           proto.String("https://avatars.githubusercontent.com/u/67905417?s=400&u=27a66894115162b6bea7e635deb7af84ae0c5c4f&v=4"),
		GongxinbuBeian: nil,
		GonganbuBeian:  nil,
	})
	if err != nil {
		return
	}
	/*
		//		ArticleSortByLogin:       commonv1.CMN_SORT_CREATE_TIME,
		//		ArticleSortByNotLogin:    commonv1.CMN_SORT_CREATE_TIME,
		//		ArticleHotShowSum:        5,
		//		ArticleHotShowWithLatest: 10,
		//		ActivityLatestShowSum:    5,
		//		DefaultPageByNewUser:     "sys_home",
		//		DefaultPageByNotLogin:    "sys_home",
		//		DefaultPageByLogin:       "sys_home",
	*/
	err = mysql.PutSystemHomeOption(invoker.Db.WithContext(ctx.Ctx), &communityv1.PutHomeOptionReq{
		IsSetHome:                proto.Bool(true),
		IsSetBanner:              proto.Bool(false),
		ArticleSortByLogin:       commonv1.CMN_SORT_CREATE_TIME,
		ArticleSortByNotLogin:    commonv1.CMN_SORT_CREATE_TIME,
		ArticleHotShowSum:        proto.Int32(5),
		ArticleHotShowWithLatest: proto.Int32(10),
		DefaultPageByNewUser:     proto.String("sys_home"),
		DefaultPageByNotLogin:    proto.String("sys_home"),
		DefaultPageByLogin:       proto.String("sys_home"),
	})
	if err != nil {
		return
	}
	/**
		//	themeInfo.ThemeName = "default"
	//	themeInfo.CustomColor = mysql.DefaultColor
	//	themeInfo.DefaultAppearance = "light"
	*/
	info := mysql.DefaultColor
	infoBytes, _ := json.Marshal(info)
	err = mysql.PutSystemTheme(invoker.Db.WithContext(ctx.Ctx), &communityv1.SetThemeReq{
		IsCustom:          false,
		ThemeName:         "default",
		CustomColor:       string(infoBytes),
		DefaultAppearance: "light",
	})
	if err != nil {
		return
	}
	return nil
}

func UpdateSsoData(ctx ejob.Context) (err error) {
	err = invoker.TokenComponent.GetAPI().UpdateClient(ctx.Ctx, eflag.String("clientId"), map[string]interface{}{"utime": time.Now().Unix()})
	if err != nil {
		return
	}
	return nil
}
