package main

import (
	"calendar-note-gin/initialize"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/language"
)

// @title           日历记事本
// @version         1.0
// @description     This is a sample server celler server.
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:9090
// @BasePath  /api/v1
// @contact.name GgoCoder
// @contact.email GgoCoder@163.com
// @schemes http
// @securityDefinitions.basic  BasicAuth
func main() {

	global.Logger.Infoln("li-calendar start!")

	// gin.SetMode(gin.TestMode) // GIN 框架调试模式
	// 配置初始化
	if config, err, errCode := initialize.Conf(getDefaultConfig()); err != nil && errCode == 0 {
		// 抛出错误
		cmn.Pln(cmn.LOG_ERROR, "配置文件创建错误:"+err.Error())
	} else if errCode == 1 {
		// 配置文件不存在，需要浏览器初始流程
		// 直接启动路由
		cmn.Pln(cmn.LOG_DEBUG, "配置文件不存在")
		global.Logger.Infoln("conf/conf.ini is not exist, start init from web!")
		initialize.RouterNeedWebInit()
	} else {
		global.Config = config
	}

	// 配置文件存在，判断是否为首次运行，是进行初始化
	if initialize.IsNeedInstall() {
		if err := initialize.InstallByConfIni(); err != nil {
			cmn.Pln("Error", err.Error())
			return
		}
	}

	initialize.RunOther()

	// 连接数据库
	if err := initialize.ConnectDb(); err != nil {
		global.Logger.Errorf("failed to init db, err:%+v", err)
		cmn.Pln("Error", "数据库错误："+err.Error())
	}

	// 用户不存在创建用户
	if !initialize.IsExistUser() {
		initialize.CreateAdminUser()
	}

	// 语言
	global.Lang = language.NewLang("zh-cn")
	// global.Lang = language.NewLang("en-us")

	// 测试
	// test()

	// 任务

	// 初始化路由
	initialize.Router()
}

func test() {
	// emailInfo := systemSetting.Email{}
	// systemSetting.GetValueByInterface("system_email", &emailInfo)
	// mailer := mail.NewMail(emailInfo.Mail, emailInfo.Password, emailInfo.Host, emailInfo.Port)
	// appName := global.Lang.Get("common.app_name")
	// title := global.Lang.GetWithFields("mail.register_vcode_title", map[string]string{
	// 	"AppName": appName,
	// })
	// content := global.Lang.GetWithFields("mail.register_vcode_content", map[string]string{
	// 	"AppName": appName,
	// 	"Minute":  "60",
	// })
	// err := mailer.SendMailOfVCode("95302870@qq.com", title, content, "123456")
	// fmt.Println("邮件发送错误", err)
}

func getDefaultConfig() map[string]map[string]string {
	return map[string]map[string]string{
		"base": {
			"http_port":        "9090",
			"source_path":      "./files",      // 存放文件的路径
			"source_temp_path": "./files/temp", // 存放文件的缓存路径
		},
		"sqlite": {
			"file_path": "./database.db",
		},
		"webdav": {
			"upload_max_size": "100", // 上传最大限制（单位：m）
		},
	}

}

func init() {
	global.Logger = cmn.InitLogger("./running.log", global.LoggerLevel)
}
