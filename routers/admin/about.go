package admin

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitAbout(router *gin.RouterGroup) {
	about := v1.ApiGroupApp.AdminApiGroup.AboutApi
	{
		router.GET("about/get", about.Get)
	}
}
