package main

import (
	"calendar-note-gin/initialize"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var RunMode = "debug"
var IsDocker = "" // 是否为docker模式

func main() {
	initialize.RUNCODE = RunMode
	initialize.ISDOCER = IsDocker

	foostr := flag.NewFlagSet("config", flag.ExitOnError)
	_ = foostr
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			// 生成配置文件
			fmt.Println("正在生成配置文件")
			cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.example.ini")
			cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.ini")
			fmt.Println("配置文件已经创建 conf/conf.ini ", "请按照自己的需求修改")
			os.Exit(0)
		}
	}

	// 配置文件初始化
	if config, err, errCode := initialize.Conf(getDefaultConfig()); err != nil && errCode == 0 {
		// 抛出错误
		cmn.Pln(cmn.LOG_ERROR, "配置文件创建错误:"+err.Error())
	} else if errCode == 1 {
		// 配置文件不存在，需要浏览器初始流程
		global.Logger.Errorln("配置文件 conf/conf.ini 不存在, 请执行 \"li-calendar config\" 来创建配置文件")
		os.Exit(1)
		// initialize.RouterNeedWebInit() // web初始化方式，暂时未开发
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
		// global.Logger.Errorln("failed to init db, err:%+v", err)
		global.Logger.Errorln("数据库连接错误", err.Error())
		os.Exit(1)
	}

	// 用户不存在创建用户
	if !initialize.IsExistUser() {
		initialize.CreateAdminUser()
	}

	// 语言
	initialize.InitLang("zh-cn") // en-us

	global.Logger.Infoln("li-calendar success start!")
	// 测试
	test()

	// 任务
	initialize.RunAfterDb()

	// 初始化路由
	initialize.Router()
}

func test() {
	reminderTime := time.Now().Add(time.Minute)
	e := models.EventReminder{}
	e.ReminderTime = reminderTime.Format(cmn.TIME_MODE_REMINDER_TIME)
	e.EventId = 2700
	e.Method = 1
	for i := 0; i < 100; i++ {
		e.ID = 0
		global.Db.Create(&e)
	}

	e.ReminderTime = reminderTime.Add(time.Minute).Format(cmn.TIME_MODE_REMINDER_TIME)
	for i := 0; i < 5; i++ {
		e.ID = 0
		global.Db.Create(&e)
	}

	// e.ReminderTime = reminderTime.Add(time.Minute).Format(cmn.TIME_MODE_REMINDER_TIME)
	// for i := 0; i < 3; i++ {
	// 	e.ID = 0
	// 	global.Db.Create(&e)
	// }

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
	}

}

func init() {
	gin.SetMode(RunMode) // GIN 运行模式
	runtimePath := "./runtime/runlog"
	if err := os.MkdirAll(runtimePath, 0777); err != nil {
		panic(err)
	}

	var level zap.AtomicLevel
	if initialize.RUNCODE == "debug" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = global.LoggerLevel
	}
	global.Logger = cmn.InitLogger(runtimePath+"/running.log", level)
}
