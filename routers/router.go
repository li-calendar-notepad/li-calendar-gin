package routers

import (
	"calendar-note-gin/api/v1/middleware"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/routers/admin"
	"calendar-note-gin/routers/common"
	"calendar-note-gin/routers/system"

	_ "calendar-note-gin/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 总路由初始化
func Routers() *gin.Engine {
	router := gin.Default()

	// 静态文件服务
	{

		cfgPath := global.Config.GetValueStringOrDefault("base", "source_path")
		// fmt.Println("静态文件", cfgPath)
		router.Static("/"+cfgPath, cfgPath) // 文件上传文件夹

	}

	// 前端文件
	// 前端文件夹存在将挂载路由
	{

		common.FSListenFile(router, "/", "frontend/index.html")
		common.FSListenFile(router, "/favicon.ico", "frontend/favicon.ico")

		common.FSListenDir(router, "/static", "frontend/static")
		common.FSListenDir(router, "/assets", "frontend/assets")

		// if ok, _ := cmn.PathExists("assets/frontend"); ok {

		// 	// router.StaticFile("/", "assets/frontend/index.html")
		// 	// router.StaticFS("/static", http.Dir("assets/frontend/static"))
		// 	// router.StaticFS("/assets", http.Dir("assets/frontend/assets"))
		// 	// router.Static("/favicon.ico", "assets/frontend/assets/favicon.ico")
		// }
	}

	publicGroup := router.Group("api/v1")
	privateGroup := router.Group("api/v1")
	adminGroup := privateGroup.Group("admin")

	// 路由
	system.InitLoginRouter(publicGroup)
	system.InitCaptchaRouter(publicGroup)
	system.InitStyleCssRouter(publicGroup)
	system.InitOpenRouter(publicGroup)

	// 需要登录
	privateGroup.Use(middleware.LoginInterceptor)
	// privateGroup.Use(middleware.LoginInterceptorDev) // 接口开发使用
	system.InitTestRouter(privateGroup)
	system.InitItemRouter(privateGroup)
	system.InitEventRouter(privateGroup)
	system.InitFileRouter(privateGroup)
	system.InitSubjectRouter(privateGroup)
	system.InitUserRouter(privateGroup)
	system.InitSpecialDayRouter(privateGroup)
	system.InitCalendarEcharts(privateGroup)

	// 管理员
	adminGroup.Use(middleware.LoginInterceptor)
	adminGroup.Use(middleware.AdminInterceptor)
	admin.InitUserRouter(adminGroup)
	admin.InitSystemSettingRouter(adminGroup)
	admin.InitSpecialRouter(adminGroup)
	admin.InitStyleRouter(adminGroup)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
