package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	api := v1.ApiGroupApp.UserApi
	router.POST("/user/getInfo", api.GetInfo)
	router.POST("/user/updatePasssword", api.UpdatePasssword)
	router.POST("/user/updateInfo", api.UpdateInfo)
}
