package initialize

import (
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"
	"log"
	"os"
	"path"
	"time"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB_DRIVER = "sqlite"
var RUNCODE = "debug"
var ISDOCER = "" // 是否为docker模式

func DbInit() (db *gorm.DB, db_err error) {
	database_drive := global.Config.GetValueStringOrDefault("base", "database_drive")
	// [mysql]
	database_mysql_host := global.Config.GetValueStringOrDefault("mysql", "host")
	database_mysql_port := global.Config.GetValueStringOrDefault("mysql", "port")
	database_mysql_username := global.Config.GetValueStringOrDefault("mysql", "username")
	database_mysql_password := global.Config.GetValueStringOrDefault("mysql", "password")
	database_mysql_db_name := global.Config.GetValueStringOrDefault("mysql", "db_name")
	// database_mysql_wait_timeout := global.Config.GetValueStringOrDefault("mysql", "wait_timeout")
	// [sqlite]
	database_sqlite_path := global.Config.GetValueStringOrDefault("sqlite", "file_path")

	if database_drive == "mysql" {
		DB_DRIVER = "mysql"
		db, db_err = MysqlConnect(database_mysql_host, database_mysql_port, database_mysql_username, database_mysql_password, database_mysql_db_name)
		if db_err != nil {
			return
		}
		db_err = CreateDatabase(db.Set("gorm:table_options", "ENGINE=InnoDB"))
	} else {
		db, db_err = SqlLite3Connect(database_sqlite_path)
		if db_err != nil {
			return
		}
		db_err = CreateDatabase(db)
		if db_err != nil {
			return
		}
	}
	models.Db = db
	global.Db = db

	return
}

// sqllite3连接
func SqlLite3Connect(filePath string) (db *gorm.DB, err error) {
	exists := false
	if exists, err = cmn.PathExists(path.Dir(filePath)); err != nil {
		return
	} else {

		// 创建文件夹
		if !exists {
			if err = os.MkdirAll(path.Dir(filePath), 0666); err != nil {
				return
			}
		}

		db, err = gorm.Open(sqlite.Open(filePath), &gorm.Config{
			Logger: GetLogger(),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}

	return
}

// mysql连接
func MysqlConnect(host, port, username, password, db_name string) (db *gorm.DB, err error) {
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: GetLogger(),
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "blog_",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)  // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxOpenConns(100) // SetMaxOpenConns 设置打开数据库连接的最大数量。
	wait_timeout := global.Config.GetValueInt("mysql", "wait_timeout")
	sqlDb.SetConnMaxLifetime(time.Duration(wait_timeout * int(time.Second))) // SetConnMaxLifetime 设置了连接可复用的最大时间。
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	return
}

// 日志
func GetLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

}

// 创建数据库
func CreateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Subject{},
		&models.Event{},
		&models.SystemSetting{},
		&models.Style{},
		&models.Item{},
		&models.File{},
		&models.SpecialDay{},
		&models.Special{},
		&models.EventReminder{},
	)
	return err
}

// func getDatebase() ...interface{} {
// 	return
// }
