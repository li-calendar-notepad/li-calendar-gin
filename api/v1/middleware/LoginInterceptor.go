package middleware

import (
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"

	"calendar-note-gin/api/v1/common/apiReturn"

	"github.com/gin-gonic/gin"
)

func LoginInterceptor(c *gin.Context) {

	// 继续执行后续的操作，再回来
	// c.Next()

	// 获得token
	token := c.GetHeader("token")

	// 没有token信息视为未登录
	if token == "" {
		apiReturn.ErrorCode(c, 1000, global.Lang.Get("login.err_not_login"), nil)
		c.Abort() // 终止执行后续的操作，一般配合return使用
		return
	}

	userInfoInterface, success := global.UserToken.Get(token)

	if !success {
		// 没有找到token信息,登录状态已经过期
		apiReturn.ErrorCode(c, 1001, global.Lang.Get("login.err_token_expire"), nil)
		c.Abort()
		return
	} else {
		if _, ok := userInfoInterface.(models.User); !ok {
			apiReturn.ErrorCode(c, 1001, global.Lang.Get("login.err_token_expire"), nil)
			c.Abort()
			return
		}
	}
	mUser := models.User{}
	// 去库中查询是否存在该用户；否则返回错误
	if info, err := mUser.GetUserInfoByToken(token); err != nil || info.ID == 0 {
		apiReturn.ErrorCode(c, 1001, global.Lang.Get("login.err_token_expire"), nil)
		c.Abort()
		return
	} else {
		// 通过
		// 设置当前用户信息
		c.Set("userInfo", info)
	}
}

// 不验证缓存直接验证库省去没有缓存每次都要手动登录的问题
func LoginInterceptorDev(c *gin.Context) {

	// 获得token
	token := c.GetHeader("token")
	mUser := models.User{}

	// 去库中查询是否存在该用户；否则返回错误
	if info, err := mUser.GetUserInfoByToken(token); err != nil || info.ID == 0 {
		apiReturn.ErrorCode(c, 1001, global.Lang.Get("login.err_token_expire"), nil)
		c.Abort()
		return
	} else {
		// 通过
		// 设置当前用户信息
		c.Set("userInfo", info)
	}
}
