package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitCaptchaRouter(router *gin.RouterGroup) {
	captchaApi := v1.ApiGroupApp.CaptchaApi
	r := router.Group("captach")
	r.GET("getImage", captchaApi.GetImage)
	// r.POST("/captach/check", captchaApi.CheckVCode)
}
