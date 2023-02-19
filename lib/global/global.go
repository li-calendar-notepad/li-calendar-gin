package global

import (
	"calendar-note-gin/lib/cache"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/language"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Log *cmn.LogStruct
var Config *cmn.IniConfig
var Db *gorm.DB
var Lang *language.LangStructObj
var UserToken *cache.GoCacheStruct
var Logger *zap.SugaredLogger
var LoggerLevel = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
var VerifyCodeCachePool *cache.GoCacheStruct
