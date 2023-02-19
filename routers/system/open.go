package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitOpenRouter(router *gin.RouterGroup) {
	styleApi := v1.ApiGroupApp.Open
	{
		router.POST("open/getOpenConfig", styleApi.GetOpenConfig)
	}

}
