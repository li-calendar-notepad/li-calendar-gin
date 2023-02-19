package admin

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitSystemSettingRouter(router *gin.RouterGroup) {
	api := v1.ApiGroupApp.AdminApiGroup.SystemSettingApi

	{
		router.POST("systemSetting/getEmail", api.GetEmail)
		router.POST("systemSetting/setEmail", api.SetEmail)
		router.POST("systemSetting/getApplicationSetting", api.GetApplicationSetting)
		router.POST("systemSetting/setApplicationSetting", api.SetApplicationSetting)
	}

}
