package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitSpecialDayRouter(router *gin.RouterGroup) {

	specialDayApi := v1.ApiGroupApp.AdminApiGroup.SpecialDayApi
	{
		router.POST("/specialDay/getList", specialDayApi.GetList)
	}
}
