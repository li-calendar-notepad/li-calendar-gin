package initialize

import (
	"calendar-note-gin/lib/cache"
	"calendar-note-gin/lib/global"
	"time"
)

func InitUserToken() {
	global.UserToken = cache.NewGoCache(5*time.Minute, 60*time.Second)
}

func InitVerifyCodeCachePool() {
	global.VerifyCodeCachePool = cache.NewGoCache(10*time.Minute, 60*time.Second)
}
