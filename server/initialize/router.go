package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/ebus"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/plugin"
	swaggerFiles "github.com/swaggo/files"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/docs"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 初始化总路由

func Routers() *gin.Engine {

	// 设置为发布模式
	if global.GVA_CONFIG.System.Env == "public" {
		gin.SetMode(gin.ReleaseMode) //DebugMode ReleaseMode TestMode
	}
	Router := gin.New()
	ebus.Publish(global.EVENTROUTERRegisterInit, Router)
	RouterRoot := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	ebus.Publish(global.EVENTROUTERRegisterBegin, Router, RouterRoot)
	InstallPlugin(RouterRoot) // 安装插件
	systemRouter := router.RouterGroupApp.System
	exampleRouter := router.RouterGroupApp.Example
	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面3行注释
	// Router.Static("/favicon.ico", "./dist/favicon.ico")
	// Router.Static("/assets", "./dist/assets")   // dist里面的静态资源
	// Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面

	Router.StaticFS(global.GVA_CONFIG.Local.StorePath, http.Dir(global.GVA_CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	// Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//global.GVA_LOG.Info("use middleware cors")

	if !global.GVA_CONFIG.System.SwaggerDisabled {
		docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix

		docs.SwaggerInfo.Version = global.Version
		docs.SwaggerInfo.Description = fmt.Sprintf(`%s

版 本 号: **%s**
编译时间: **%s**
最近提交: **%s**
提 交 人: **%s**
提交信息: **%s**

%s
`, docs.SwaggerInfo.Description, global.Version, global.BuildTime, global.GitTag, global.GitAuthor, global.GitCommitMsg, strings.ReplaceAll(global.GitCommitLog, "<br/>", "\n"))

		RouterRoot.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		global.GVA_LOG.Info("register swagger handler")
	} else {
		RouterRoot.GET("/swagger/*any", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"ip":  ctx.ClientIP(),
				"now": time.Now(),
			})
		})
	}

	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := RouterRoot.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关
	}
	PrivateGroup := RouterRoot.Group("")
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)    // 注册功能api路由
		systemRouter.InitJwtRouter(PrivateGroup)                 // jwt相关路由
		systemRouter.InitUserRouter(PrivateGroup)                // 注册用户路由
		systemRouter.InitMenuRouter(PrivateGroup)                // 注册menu路由
		systemRouter.InitSystemRouter(PrivateGroup)              // system相关路由
		systemRouter.InitCasbinRouter(PrivateGroup)              // 权限相关路由
		systemRouter.InitAutoCodeRouter(PrivateGroup)            // 创建自动化代码
		systemRouter.InitAuthorityRouter(PrivateGroup)           // 注册角色路由
		systemRouter.InitSysDictionaryRouter(PrivateGroup)       // 字典管理
		systemRouter.InitAutoCodeHistoryRouter(PrivateGroup)     // 自动化代码历史
		systemRouter.InitSysOperationRecordRouter(PrivateGroup)  // 操作记录
		systemRouter.InitSysDictionaryDetailRouter(PrivateGroup) // 字典详情管理
		systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup)  // 字典详情管理
		systemRouter.InitChatGptRouter(PrivateGroup)             // chatGpt接口

		exampleRouter.InitCustomerRouter(PrivateGroup)              // 客户路由
		exampleRouter.InitFileUploadAndDownloadRouter(PrivateGroup) // 文件上传下载功能路由

	}

	if err := global.GVA_PLUGIN.ForeachUnsafe(func(p plugin.Base) error {
		if m1, ok := p.(plugin.WithRouter); ok {
			if e1 := m1.PluginRouter(PrivateGroup, PublicGroup); e1 != nil {
				global.GVA_LOG.Error("plugin ["+p.PluginName()+"] router error", zap.Error(e1))
				return e1
			}
		}
		return nil
	}); err != nil {
		os.Exit(1)
	}

	ebus.Publish(global.EVENTROUTERRegisterSuccess, Router, PrivateGroup, PublicGroup)

	global.GVA_LOG.Info("router register success")
	return Router
}
