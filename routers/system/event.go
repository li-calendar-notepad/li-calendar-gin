package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitEventRouter(router *gin.RouterGroup) {
	itemApi := v1.ApiGroupApp.ItemApi
	eventApi := v1.ApiGroupApp.EventApi

	// 验证项目的权限(有访问密码的需要验证访问token)
	private := router.Group("")
	private.Use(itemApi.MiddlewarePrivate)
	{
		private.POST("/item/event/createOne", eventApi.CreateOne)
		private.POST("/item/event/updateByEventId", eventApi.UpdateByEventId)
		private.POST("/item/event/getListByTimeScope", eventApi.GetListByTimeScope)
		private.POST("/item/event/getListAndSpecialDayByTimeScope", eventApi.GetListAndSpecialDayByTimeScope)
		private.POST("/item/event/deleteByIdAndItemId", eventApi.DeleteByIdAndItemId)
		private.POST("/item/event/getInfo", eventApi.GetInfo)
		private.POST("/item/event/searchWord", eventApi.SearchWord)
	}

}
