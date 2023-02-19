package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitItemRouter(router *gin.RouterGroup) {
	itemApi := v1.ApiGroupApp.ItemApi

	router.POST("/item/getMyItemList", itemApi.GetMyItemList)
	router.POST("/item/getPublicInfoByItemId", itemApi.GetPublicInfoByItemId)
	router.POST("/item/checkVisitPassword", itemApi.CheckVisitPassword)
	router.POST("/item/create", itemApi.Create)
	router.POST("/item/forgotVisitPassword", itemApi.ForgotVisitPassword)
	// 验证项目的权限(有访问密码的需要验证访问token)
	private := router.Group("")
	// 验证项目权限-中间件
	private.Use(itemApi.MiddlewarePrivate)
	{
		private.POST("/item/getPrivateInfo", itemApi.GetPrivateInfo)
		private.POST("/item/saveConfig", itemApi.SaveConfig)
		private.POST("/item/getConfig", itemApi.GetConfig)
		private.POST("/item/deleteOne", itemApi.DeleteOne)
		private.POST("/item/getSummaryData", itemApi.GetSummaryData)
	}

}
