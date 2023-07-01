package server

import (
	"fmt"
	"net/http"
	"strings"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/api/article"
	"ecodepost/bff/pkg/server/api/column"
	"ecodepost/bff/pkg/server/api/community"
	"ecodepost/bff/pkg/server/api/file"
	"ecodepost/bff/pkg/server/api/home"
	"ecodepost/bff/pkg/server/api/logger"
	"ecodepost/bff/pkg/server/api/my"
	"ecodepost/bff/pkg/server/api/pms"
	profile "ecodepost/bff/pkg/server/api/public-profile"
	"ecodepost/bff/pkg/server/api/question"
	"ecodepost/bff/pkg/server/api/space"
	"ecodepost/bff/pkg/server/api/theme"
	"ecodepost/bff/pkg/server/api/upload"
	"ecodepost/bff/pkg/server/api/user"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/server/mw"
	"ecodepost/bff/pkg/server/ssoapi"
	"ecodepost/bff/pkg/server/ui"
	"ecodepost/bff/pkg/service"
	"ecodepost/bff/pkg/service/ssoservice"
	uploadv1 "ecodepost/pb/upload/v1"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/server/egin"
)

var H = bffcore.Handle

func ServeHttp() *egin.Component {
	r := egin.Load("server.bff").Build(egin.WithEmbedFs(ui.WebUI))
	r.Invoker(
		invoker.Init,
		ssoservice.Init,
		service.Init,
		func() error {
			return registerRouter(r)
		},
	)
	return r
}

// BACKEND FOR FRONTEND
func registerRouter(r *egin.Component) error {
	r.NoRoute(H(func(ctx *bffcore.Context) {
		// API
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/") || strings.HasPrefix(ctx.Request.URL.Path, "/sso/api/") {
			ctx.String(http.StatusNotFound, "404 not found")
			return
		}
		// 说明是file图片
		if strings.HasPrefix(ctx.Request.URL.Path, "/"+econf.GetString("oss.prefix")) {
			image, err := invoker.GrpcUpload.ShowImage(ctx, &uploadv1.ShowImageReq{
				Path: ctx.Request.URL.Path,
			})
			if err != nil {
				ctx.EgoJsonI18N(err)
				return
			}
			ctx.Writer.Write(image.File)
			return
		}

		maxAge := econf.GetInt("server.http.maxAge")
		if maxAge == 0 {
			maxAge = 86400
		}
		if strings.HasSuffix(ctx.Request.URL.Path, ".js") || strings.HasSuffix(ctx.Request.URL.Path, ".css") || strings.HasSuffix(ctx.Request.URL.Path, ".png") || strings.HasSuffix(ctx.Request.URL.Path, ".img") || strings.HasSuffix(ctx.Request.URL.Path, ".ico") {
			ctx.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			ctx.FileFromFS(ctx.Request.URL.Path, r.HTTPEmbedFs())
			return
		}
		// todo 因为有一些动态路由，html，无法渲染，暂时这么处理
		//ctx.FileFromFS(ctx.Request.URL.Path, r.HTTPEmbedFs())
		ctx.FileFromFS("/", r.HTTPEmbedFs())
		return
	}))

	apiSsoGroup := r.Group("/api/sso")
	apiSsoGroup.GET("/wechat/validate", bffcore.Handle(ssoapi.WechatValidate))                                // 短信登录
	apiSsoGroup.GET("/login/wechat", ssoapi.AuthExist(), bffcore.Handle(ssoapi.LoginWechat))                  // 微信登录
	apiSsoGroup.POST("/login/basic", ssoapi.AuthExist(), bffcore.Handle(ssoapi.LoginDirect))                  // 密码登录
	apiSsoGroup.GET("/code/wechatWeb", ssoapi.AuthExist(), bffcore.Handle(ssoapi.CodeWechatWeb))              // 短信登录
	apiSsoGroup.GET("/code/wechatH5", ssoapi.AuthExist(), bffcore.Handle(ssoapi.CodeWechatH5))                // 短信登录
	apiSsoGroup.POST("/register", ssoapi.AuthExist(), bffcore.Handle(ssoapi.Register))                        // 密码登录
	apiSsoGroup.POST("/phone/sendRegisterCode", ssoapi.AuthExist(), bffcore.Handle(ssoapi.RegisterPhoneCode)) // 发送注册手机验证码

	r.Use(mw.I18nCookie())
	r.GET("/api/oauth/user", invoker.Sso.MayBeLogin(), H(my.OauthUserInfo)) // 用户信息

	// 选择当前分区的cmt guid，这个要单独出来，因为可能从域名里获取
	r.GET("/api/cmt/detail", invoker.Sso.MayBeLogin(), H(community.Detail))

	// 有些请求可能有登录态，也可能没有登录态
	maybeLoginGroup := r.Group("/api")
	// TODO 只有公开的社区才可以，在没有登录态的时候登录 TODO 非公开社区，强制需要登录
	maybeLoginGroup.Use(invoker.Sso.MayBeLogin())
	maybeLoginGroup.GET("/home/page", H(home.Page))                                    // 首页的数据信息
	maybeLoginGroup.GET("/home/files", H(home.Files))                                  // 首页的数据信息
	maybeLoginGroup.GET("/cmt/space-all", H(space.ListSpaceAndGroup))                  // 查看完整group及其下space列表,以及emoji列表等
	maybeLoginGroup.GET("/spaces/:guid/permission", H(space.Permission))               // 用于个人数据权限，不能缓存
	maybeLoginGroup.GET("/spaces/-/groups/:guid/permission", H(space.GroupPermission)) // 用于个人数据权限，不能缓存
	maybeLoginGroup.GET("/spaces/:guid", H(space.Info))                                // 获取space信息，有成员个数，用于公开的数据，可以做缓存
	maybeLoginGroup.GET("/spaces/:guid/column", H(space.ColumnInfo))                   // 获取space信息，有成员个数，用于公开的数据，可以做缓存
	maybeLoginGroup.GET("/cmt/managers", H(community.Managers))                        // 社区管理员信息
	// 文件相关接口
	maybeLoginGroup.GET("/files", H(file.ListPage)) // (0522) 文档列表 (0522 新增用户昵称，头像，摘要)
	maybeLoginGroup.GET("/files/:guid", H(file.GetInfo))
	maybeLoginGroup.GET("/files/:guid/subList", H(file.SubListPage)) // 回答列表
	maybeLoginGroup.GET("/files/-/recommends", H(file.Recommends))   // 推荐列表 (0522 新增用户昵称，头像，摘要)
	maybeLoginGroup.GET("/files/-/spaceTops", H(file.SpaceTops))     // 置顶列表 (0522 新增用户昵称，头像，摘要)
	maybeLoginGroup.GET("/files/-/stats", H(file.Stat))
	maybeLoginGroup.GET("/files/:guid/comments", H(file.CommentList))                            // 评论列表
	maybeLoginGroup.GET("/files/-/comments/:commentGuid/childComment", H(file.ChildCommentList)) // 子评论列表
	maybeLoginGroup.GET("/columns/-/listPermission", H(column.ListPermission))                   // 权限信息
	maybeLoginGroup.GET("/columns/-/files", H(column.ListFiles))                                 // 树型目录
	// 社区里的个人中心信息
	maybeLoginGroup.GET("/notifications/-/total", H(my.NotificationList)) // 我收到的消息列表
	maybeLoginGroup.GET("/files/:guid/permission", H(file.Permission))
	maybeLoginGroup.GET("/files/-/permissions", H(file.PermissionList)) // 权限信息

	// 个人公开数据
	maybeLoginGroup.GET("/users/:name/total", H(profile.UserTotal))
	maybeLoginGroup.GET("/users/:name/articles", H(profile.ArticlesList))   // (0730) 用户文章列表
	maybeLoginGroup.GET("/users/:name/questions", H(profile.QAList))        // (0730) 用户问题、回答列表，如果没有parent guid是问题、有parent guid是回答
	maybeLoginGroup.GET("/users/:name/followers", H(profile.FollowersList)) // 指定用户的followers列表, followers:关注你的人
	maybeLoginGroup.GET("/users/:name/following", H(profile.FollowingList)) // 指定用户的following列表, following:你关注的人

	apiGroup := r.Group("/api")
	apiGroup.Use(invoker.Sso.MustCheckToken())
	apiGroup.GET("/oauth/logout", H(invoker.Sso.OauthLogout)) // 清除team guid的cookie, 退出登录
	apiGroup.POST("/upload/token", H(upload.Token))           // 获取上传TOKEN信息
	apiGroup.POST("/upload/path", H(upload.Path))             // 获取访问路径

	apiGroup.PUT("/my/communities/:guid", H(community.Update))              // 更新社区
	apiGroup.PUT("/my/communities/:guid/banner", H(community.UpdateBanner)) // 更新社区banner信息

	// 个人->关注相关相关接口
	apiGroup.POST("/my/following/:uid", H(my.FollowingCreate))   // 我关注某个用户
	apiGroup.DELETE("/my/following/:uid", H(my.FollowingDelete)) // 我取关某个用户

	// 个人通知相关接口
	apiGroup.PUT("/notifications/-/audits/:auditId/pass", H(my.NotificationAuditPass))     // 我收到的审核消息列表
	apiGroup.PUT("/notifications/-/audits/:auditId/reject", H(my.NotificationAuditReject)) // 我收到的审核消息列表

	// 个人其他接口
	apiGroup.GET("/my", H(user.My))                      // 获取个人信息
	apiGroup.PUT("/my/attr", H(user.UpdateAttr))         // 修改我的通用属性
	apiGroup.PUT("/my/nickname", H(user.UpdateNickname)) // 修改我的昵称
	apiGroup.PUT("/my/avatar", H(user.UpdateAvatar))     // 修改头像，需要先调用upload token获取信息
	apiGroup.PUT("/my/phone", H(user.UpdatePhone))       // 修改手机号
	apiGroup.PUT("/my/email", H(user.UpdateEmail))       // 修改邮箱

	apiGroup.GET("/users/:name", H(user.Info)) // 获取单个用户

	// 社区相关接口
	apiGroup.GET("/cmt/recommendLogos", H(community.ListLogos))   // 社区推荐logo列表
	apiGroup.GET("/cmt/recommendCovers", H(community.ListCovers)) // 社区推荐封面列表

	// 收藏相关接口
	apiGroup.POST("/my/collection-groups", H(my.CollectionGroupCreate))            // (0603) 创建一个收藏分组
	apiGroup.GET("/my/collection-groups", H(my.CollectionGroupList))               // (0603) 收藏分组列表
	apiGroup.DELETE("/my/collection-groups/:cgid", H(my.CollectionGroupDelete))    // (0603) 删除一个收藏分组
	apiGroup.PUT("/my/collection-groups/:cgid", H(my.CollectionGroupUpdate))       // (0603) 修改一个收藏分组
	apiGroup.GET("/my/collection-groups/:cgid/collections", H(my.CollectionList))  // (0730) 查询一个收藏分组下收藏列表
	apiGroup.POST("/my/collection-groups/-/collections", H(my.CollectionCreate))   // (0730) 新增一个收藏，可以加入多个分组id，详细玩法参考B站
	apiGroup.DELETE("/my/collection-groups/-/collections", H(my.CollectionDelete)) // (0730) 删除一个收藏，可以从多个收藏分组删除一个收藏

	// 加入第一个空间的时候，可能不在社区，需要不判断是否在社区
	apiGroup.POST("/spaces/:guid/apply", H(space.ApplyMember))               // (0725) 申请加入
	apiGroup.POST("/spaces/:guid/quit", H(space.QuitMember))                 // (0725) 退出空间
	apiGroup.POST("/spaces/-/groups/:guid/apply", H(space.ApplyGroupMember)) // (0725) 申请加入
	apiGroup.GET("/spaces/-/member-status", H(space.GetMemberStatus))        // (0715) 获取成员状态

	apiGroup.PUT("/files/:guid/increaseEmoji", H(file.IncreaseEmoji))        // (0521) 传入一个emoji id，点赞
	apiGroup.PUT("/files/:guid/decreaseEmoji", H(file.DecreaseEmoji))        // (0521) 传入一个emoji id，去掉点赞
	apiGroup.POST("/files/-/comments", H(file.CreateComment))                // (0524) 创建评论
	apiGroup.DELETE("/files/-/comments/:commentGuid", H(file.DeleteComment)) // (0524) 删除评论
	apiGroup.POST("/files/-/upload/local", H(file.UploadLocalFile))

	// 社区其他接口
	apiGroup.PUT("/cmt/space-trees/change-sort", H(space.TreeChangeSort))   // (0901) space树信息 (0522 增加可见度，权限信息）
	apiGroup.PUT("/cmt/space-trees/change-group", H(space.TreeChangeGroup)) // (0901) space树信息 (0522 增加可见度，权限信息）

	// spaces->groups 相关接口
	apiGroup.POST("/spaces/-/groups", H(space.CreateGroup))                  // (0522) 创建space group (0522 增加可见度，是否打开用户查看列表信息）
	apiGroup.GET("/spaces/-/groups/:guid", H(space.GroupInfo))               // (0715) 获取space group信息，有成员个数
	apiGroup.PUT("/spaces/-/groups/:guid", H(space.UpdateGroup))             // (0522) 更新space group (0522 增加可见度，是否打开用户查看列表信息）
	apiGroup.DELETE("/spaces/-/groups/:guid", H(space.DeleteGroup))          // (0519) 删除一个space group
	apiGroup.GET("/spaces/-/groups/:guid/members", H(space.GroupMemberList)) // (0529) space group的成员列表
	//apiGroup.GET("/spaces/-/groups/:guid/searchMembers", H(space.SearchGroupMember)) // (0529) space group的搜索成员
	apiGroup.POST("/spaces/-/groups/:guid/members", H(space.CreateGroupMember))   // (0529) space group添加成员
	apiGroup.DELETE("/spaces/-/groups/:guid/members", H(space.DeleteGroupMember)) // (0529) space group删除成员

	// spaces 相关接口
	apiGroup.POST("/spaces", H(space.Create))                       // (0522) 创建一个space (0522 增加空间类型，空间第三方类型，空间布局，可见度）
	apiGroup.PUT("/spaces/:guid", H(space.Update))                  // (0522) 更新space基础信息 (0522 增加空间类型，空间第三方类型，空间布局，可见度）
	apiGroup.DELETE("/spaces/:guid", H(space.Delete))               // (0519) 删除一个space
	apiGroup.GET("/spaces/:guid/members", H(space.MemberList))      // (0529) 成员列表
	apiGroup.POST("/spaces/:guid/members", H(space.CreateMember))   // (0529) 添加成员
	apiGroup.DELETE("/spaces/:guid/members", H(space.DeleteMember)) // (0529) 删除成员
	apiGroup.GET("/spaces/:guid/emojis", H(space.Emojis))           // (0604) 空间下emoji列表

	// 文章相关接口
	apiGroup.GET("/articles/-/recommendCovers", H(article.ListCovers))          // (0816) 文章推荐封面列表
	apiGroup.POST("/articles", H(article.CreateArticle))                        // (0519) 创建文档
	apiGroup.PUT("/articles/:guid", H(article.UpdateArticle))                   // (0519) 更新、发布文档
	apiGroup.DELETE("/articles/:guid", H(article.DeleteArticle))                // (0521) 删除文档
	apiGroup.PUT("/articles/:guid/spaceTop", H(article.SpaceTop))               // (0521) 置顶文档
	apiGroup.PUT("/articles/:guid/cancelSpaceTop", H(article.CancelSpaceTop))   // (0521) 取消置顶文档
	apiGroup.PUT("/articles/:guid/recommend", H(article.Recommend))             // (0521) 推荐文档
	apiGroup.PUT("/articles/:guid/cancelRecommend", H(article.CancelRecommend)) // (0521) 取消推荐文档
	apiGroup.PUT("/articles/:guid/openComment", H(article.OpenComment))         // (0521) 文档打开评论
	apiGroup.PUT("/articles/:guid/closeComment", H(article.CloseComment))       // (0521) 文档关闭评论

	// 问答相关接口
	apiGroup.POST("/questions", H(question.Create))                                                             // (0527) 创建问题
	apiGroup.GET("/questions/:guid", H(question.Info))                                                          // (0521) 查询某个社区活动详情
	apiGroup.PUT("/questions/:guid", H(question.Update))                                                        // (0527) 更新问题
	apiGroup.DELETE("/questions/:guid", H(question.Delete))                                                     // (0527) 删除问题
	apiGroup.POST("/questions/:guid/answers", H(question.CreateAnswer))                                         // (0528) 创建回答
	apiGroup.PUT("/questions/-/answers/:answerGuid", H(question.UpdateAnswer))                                  // (0603) 更新回答
	apiGroup.DELETE("/questions/-/answers/:answerGuid", H(question.DeleteAnswer))                               // (0603) 删除回答
	apiGroup.POST("/questions/:guid/like", H(question.LikeQuestion))                                            // (0809) 点赞问题
	apiGroup.DELETE("/questions/:guid/like", H(question.UndoLikeQuestion))                                      // (0809) 取消点赞问题
	apiGroup.POST("/questions/-/answers/:answerGuid/like", H(question.LikeAnswer))                              // (0809) 点赞回答
	apiGroup.DELETE("/questions/-/answers/:answerGuid/like", H(question.UndoLikeAnswer))                        // (0809) 取消点赞回答
	apiGroup.POST("/questions/-/answers/:answerGuid/comments/:commentGuid/like", H(question.LikeComment))       // (0809) 点赞评论
	apiGroup.DELETE("/questions/-/answers/:answerGuid/comments/:commentGuid/like", H(question.UndoLikeComment)) // (0809) 取消点赞评论

	// 社区管理后台
	apiGroup.GET("/cmt/-/searchMembers", H(community.SearchMember))    // (0530) 社区搜索成员 ，必须放这个后面，否则没社区信息
	apiGroup.GET("/cmtAdmin/user/memberList", H(community.MemberList)) // (0530) 社区搜索成员 ，必须放这个后面，否则没社区信息

	apiGroup.GET("/cmtAdmin/pms/managers/members", H(pms.ManagerMemberList))                                           // (0803) 权限成员数据 :managerType superAdmin 为 超级管理员，admin 为管理员
	apiGroup.POST("/cmtAdmin/pms/managers/members", H(pms.CreateManagerMember))                                        // (0803) 添加权限成员
	apiGroup.DELETE("/cmtAdmin/pms/managers/members/:uid", H(pms.DeleteManagerMember))                                 // (0803) 删除权限成员
	apiGroup.GET("/cmtAdmin/pms/roles", H(pms.RoleList))                                                               // (0613) 获取全部role列表
	apiGroup.GET("/cmtAdmin/pms/roles/-/users/:uid/roleIds", H(pms.UserRoleIds))                                       // (0714) 获取某个用户的role ids
	apiGroup.POST("/cmtAdmin/pms/roles", H(pms.CreateRole))                                                            // (0616) 创建Role
	apiGroup.PUT("/cmtAdmin/pms/roles/:roleId", H(pms.UpdateRole))                                                     // (0707) 更新角色
	apiGroup.DELETE("/cmtAdmin/pms/roles/:roleId", H(pms.DeleteRole))                                                  // (0715) 删除角色
	apiGroup.GET("/cmtAdmin/pms/roles/:roleId/members", H(pms.RoleMemberList))                                         // (0613) 获取某个role的成员列表
	apiGroup.GET("/cmtAdmin/pms/roles/:roleId/permissions", H(pms.RolePermission))                                     // (0613) 获取某个role的权限列表
	apiGroup.GET("/cmtAdmin/pms/roles/:roleId/spaceGroupPermissions", H(pms.RoleSpaceGroupPermission))                 // (0613) 获取某个role的权限列表
	apiGroup.GET("/cmtAdmin/pms/roles/:roleId/spacePermissions", H(pms.RoleSpacePermission))                           // (0613) 获取某个role的权限列表
	apiGroup.GET("/cmtAdmin/pms/roles/:roleId/spaces/:guid/initPermissions", H(pms.RoleSpaceInitPermission))           // (0707) role的space的初始权限点
	apiGroup.GET("/cmtAdmin/pms/roles/:roleId/spaceGroups/:guid/initPermissions", H(pms.RoleSpaceGroupInitPermission)) // (0707) role的space group的初始权限点
	apiGroup.PUT("/cmtAdmin/pms/roles/:roleId/permissions", H(pms.PutRolePermission))                                  // (0613) 修改某个role权限点
	apiGroup.PUT("/cmtAdmin/pms/roles/:roleId/spaceGroupPermissions", H(pms.PutRoleSpaceGroupPermission))              // (0613) 修改某个role权限点
	apiGroup.PUT("/cmtAdmin/pms/roles/:roleId/spacePermissions", H(pms.PutRoleSpacePermission))                        // (0613) 修改某个role权限点
	apiGroup.POST("/cmtAdmin/pms/roles/:roleId/members", H(pms.CreateRoleMember))                                      // (0613) 添加某个role成员
	apiGroup.DELETE("/cmtAdmin/pms/roles/:roleId/members", H(pms.DeleteRoleMember))                                    // (0613) 删除某个role成员

	// 操作日志记录
	apiGroup.GET("/cmtAdmin/logger/listPage", H(logger.ListPage))                   // (0805) 操作日志记录数据
	apiGroup.GET("/cmtAdmin/logger/eventAndGroupList", H(logger.EventAndGroupList)) // (0805) 用于筛选的列表数据

	// 后台管理员的应用列表
	apiGroup.GET("/cmtAdmin/home", H(home.Get)) // (0919) 首页信息
	apiGroup.PUT("/cmtAdmin/home", H(home.Put)) // (0919) 首页设置

	// 后台管理员的主题颜色
	apiGroup.GET("/cmtAdmin/theme", H(theme.Get))
	apiGroup.PUT("/cmtAdmin/theme", H(theme.Put))

	// 专栏
	apiGroup.POST("/columns", H(column.Create))                         // (0831) 创建文档
	apiGroup.PUT("/columns/-/change-sort", H(column.SidebarChangeSort)) // (0831) 修改树型目录顺序
	apiGroup.PUT("/columns/:guid", H(column.UpdateArticle))             // (0831) 更新、发布文档
	apiGroup.DELETE("/columns/:guid", H(column.DeleteArticle))          // (0831) 删除文档
	return nil
}
