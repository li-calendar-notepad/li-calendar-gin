package admin

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitSpecialRouter(router *gin.RouterGroup) {
	specialDayApi := v1.ApiGroupApp.AdminApiGroup.SpecialDayApi
	{
		router.POST("specialDay/specialDayImport", specialDayApi.SpecialDayImport)
		router.POST("specialDay/getInfoDays", specialDayApi.GetInfoDays)
		router.POST("specialDay/getList", specialDayApi.GetList)
		router.POST("specialDay/deleteOne", specialDayApi.DeleteOne)
	}
}
