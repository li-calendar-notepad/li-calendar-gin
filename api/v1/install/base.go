package install

import (
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/first_init"
	"calendar-note-gin/server"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// 首次使用初始化
func Install(c *gin.Context) {

	if !first_init.IsNeedInstall() {
		// server.ApiReturnError(c, -1, "already initialized completed")
		server.ApiReturnError(c, -1, "已经初始化过，无需再次初始化")
		return
	}

	source_path := c.PostForm("base_file_path")
	source_temp_path := c.PostForm("base_file_temp_path")

	// base_new_user_mug := c.PostForm("base_new_user_mug")

	database_drive := c.PostForm("database_drive")

	database_sqlite_path := c.PostForm("database_sqlite_path")

	database_mysql_host := c.PostForm("database_mysql_host")
	database_mysql_port := c.PostForm("database_mysql_port")
	database_mysql_db_name := c.PostForm("database_mysql_db_name")
	database_mysql_username := c.PostForm("database_mysql_username")
	database_mysql_password := c.PostForm("database_mysql_password")
	database_mysql_wait_timeout := c.PostForm("database_mysql_wait_timeout")

	admin_user_username := c.PostForm("admin_user_username")
	admin_user_password := c.PostForm("admin_user_password")

	installFileName := "./conf.ini"

	if err := first_init.CreateConf(); err != nil {
		server.ApiReturnError(c, -1, err.Error())
		return
	}
	// 加载配置
	cmn.Config = cmn.NewIniConfig(installFileName)
	// 更新配置项
	cmn.Config.SetValue("mysql", "host", database_mysql_host)
	cmn.Config.SetValue("mysql", "port", database_mysql_port)
	cmn.Config.SetValue("mysql", "username", database_mysql_username)
	cmn.Config.SetValue("mysql", "password", database_mysql_password)
	cmn.Config.SetValue("mysql", "db_name", database_mysql_db_name)
	cmn.Config.SetValue("mysql", "wait_timeout", database_mysql_wait_timeout)
	// [sqlite]
	cmn.Config.SetValue("sqlite", "file_path", database_sqlite_path)

	// [base]
	cmn.Config.SetValue("base", "source_path", source_path)
	cmn.Config.SetValue("base", "source_temp_path", source_temp_path)
	cmn.Config.SetValue("base", "database_drive", database_drive)
	cmn.Config.SetValue("base", "http_port", "9090")

	// [init]
	cmn.Config.SetValue("init", "admin_username", admin_user_username)
	cmn.Config.SetValue("init", "admin_password", admin_user_password)

	if err := first_init.InstallByConfIni(); err != nil {
		server.ApiReturnError(c, -1, err.Error())
		// 初始化失败删除配置文件
		fmt.Println("初始化失败，删除配置文件")
		os.Remove(installFileName)
		return
	}

	server.ApiReturnRight(c, "")
}

// 是否需要初始化
func IsNeedInstall(c *gin.Context) {

	if !first_init.IsNeedInstall() {
		// server.ApiReturnError(c, -1, "already initialized completed")
		server.ApiReturnError(c, -1, "已经初始化过，无需再次初始化")
		return
	}
	server.ApiReturnRight(c, "")
}
