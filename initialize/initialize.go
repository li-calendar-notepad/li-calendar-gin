package initialize

import (
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/language"
	"calendar-note-gin/models"
	"calendar-note-gin/routers"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 连接数据库
func ConnectDb() error {
	_, err := DbInit()
	if err != nil {
		return err
	}
	return nil
}

// 配置初始化
// errCode=1 说明初始化流程
func Conf(defaultConfig map[string]map[string]string) (config *cmn.IniConfig, err error, errCode int) {
	CreateConfExample()
	exists, err := cmn.PathExists("conf/conf.ini")
	if exists {
		config = cmn.NewIniConfig("conf/conf.ini") // 读取配置
		config.Default = defaultConfig
	} else if err != nil {

	} else {
		// docker 运行模式，生成配置文件
		if ISDOCER != "" {
			cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.ini")
			config = cmn.NewIniConfig("conf/conf.ini") // 读取配置
			config.Default = defaultConfig
		} else {
			errCode = 1
		}
	}
	return
}

// 生成示例配置文件
func CreateConfExample() (err error) {
	// 查看配置示例文件是否存在，不存在创建（分别为示例配置和配置文件）
	exists, err := cmn.PathExists("conf/conf.example.ini")
	if err != nil {
		return
	}
	if !exists {
		if err = cmn.AssetsTakeFileToPath("conf.example.ini", "conf/conf.example.ini"); err != nil {
			return
		}
	}

	return nil
}

// 路由初始化
func Router() {
	r := routers.Routers()
	port := global.Config.GetValueString("base", "http_port")
	r.Run(":" + port)
}

// 需要web进行初始化
func RouterNeedWebInit() {
	port := global.Config.GetValueString("base", "http_port")
	fmt.Println("准备启动web")
	cmn.Pln("Info", "首次运行请浏览器前往：http://localhost:"+port+"/#/install")
	r := gin.Default()
	r.Run(":" + port)
}

// 是否需要安装初始化
func IsNeedInstall() bool {
	exists, _ := cmn.PathExists("conf/conf.ini")
	// 如果文件不存在 ||　存在安装时间为空　|| 存在版本不等于1 都需要初始化
	if !exists || global.Config.GetValueString("build", "install_time") == "" || global.Config.GetValueString("build", "conf_version") != "1" {
		return true
	} else {
		return false
	}
}

// 根据配置文件初始化
func InstallByConfIni() error {

	// [init]
	admin_user_username := global.Config.GetValueStringOrDefault("init", "admin_username")
	admin_user_password := global.Config.GetValueStringOrDefault("init", "admin_password")
	if !cmn.VerifyFormat(cmn.VERIFY_EXP_USERNAME, admin_user_username) {
		return errors.New("管理员账号由5-50位组成，可以是字母和数字")
	}

	if len(admin_user_password) < 6 || !cmn.VerifyFormat(cmn.VERIFY_EXP_PASSWORD, admin_user_password) {
		return errors.New("管理员密码由6-16位组成，可以是数字、字母和.、&、@")
	}

	// 测试数据库连接
	// db, db_err := DbInit()
	// // 数据库链接错误
	// if db_err != nil {
	// 	return db_err
	// }

	// // 创建数据库
	// db_err = Create_datebase(db)
	// db_data_init()
	// if db_err != nil {
	// 	return db_err
	// }

	// // 创建管理员
	// if createUserErr := CreateAdminUser("Admin", admin_user_username, admin_user_password); createUserErr != nil {
	// 	cmn.Pln("Warning", "用户已经存在，不创建，仅修改密码")
	// }

	// cmn.Pln("Info", "===================================")
	// cmn.Pln("Info", "请牢记以下账号密码，登陆后可修改")
	// cmn.Pln("Info", "===================================")
	// cmn.Pln("Info", "登录账号："+admin_user_username)
	// cmn.Pln("Info", "登录密码："+admin_user_password)
	// cmn.Pln("Info", "===================================")

	// 初始化完成更新时间 修改配置文件
	global.Config.SetValue("build", "install_time", strconv.Itoa(int(time.Now().Unix())))
	global.Config.DeleteSection("init") // 删除组
	return nil
}

// 是否存在用户
func IsExistUser() bool {
	userInfo := models.User{}
	return global.Db.Model(&models.User{}).First(&userInfo).Error != gorm.ErrRecordNotFound
}

// 创建管理员用户
func CreateAdminUser() error {
	username := "admin" + time.Now().Format("2006")
	password := "123456"
	cmn.Pln("Info", "===================================")
	cmn.Pln("Info", "请牢记以下账号密码，登陆后可修改")
	cmn.Pln("Info", "===================================")
	cmn.Pln("Info", "登录账号："+username)
	cmn.Pln("Info", "登录密码："+password)
	cmn.Pln("Info", "===================================")

	password = cmn.PasswordEncryption(password)

	newUser := models.User{
		Username: username,
		Password: password,
		Name:     "超级管理",
		Status:   1,
		Role:     1,
	}
	return global.Db.Create(&newUser).Error
}

// 运行其他初始化
func RunOther() {
	InitUserToken()
	InitVerifyCodeCachePool()
}

func InitLang(lang string) {
	langPath := "lang/" + lang + ".ini"
	exists, err := cmn.PathExists(langPath)
	if err != nil {
		global.Logger.Errorln("语言文件不存在", err.Error())
		os.Exit(1)
	}

	// 生成语言文件
	if !exists {
		global.Logger.Infoln("输出语言文件:", langPath)
		err := cmn.AssetsTakeFileToPath("lang/zh-cn.ini", "lang/zh-cn.ini")
		if err != nil {
			global.Logger.Errorln("输出语言文件出错:", err.Error())
			os.Exit(1)
		}
		err = cmn.AssetsTakeFileToPath("lang/en-us.ini", "lang/en-us.ini")
		if err != nil {
			global.Logger.Errorln("输出语言文件出错:", err.Error())
			os.Exit(1)
		}
	}

	exists, err = cmn.PathExists(langPath)
	if err != nil || !exists {
		global.Logger.Errorln("语言文件不存在:", langPath)
		os.Exit(1)
	}

	global.Lang = language.NewLang(langPath)

}
