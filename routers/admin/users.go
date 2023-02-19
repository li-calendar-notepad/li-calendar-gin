package admin

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	userApi := v1.ApiGroupApp.AdminApiGroup.UsersApi
	{
		router.POST("users/create", userApi.Create)
		router.POST("users/update", userApi.Update)
		router.POST("users/getList", userApi.GetList)
		router.POST("users/deletes", userApi.Deletes)
	}
}
