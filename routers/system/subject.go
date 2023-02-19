package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitSubjectRouter(router *gin.RouterGroup) {
	itemApi := v1.ApiGroupApp.ItemApi
	subjectApi := v1.ApiGroupApp.SubjectApi
	_ = subjectApi

	// 验证项目的权限(有访问密码的需要验证访问token)
	private := router.Group("")
	private.Use(itemApi.MiddlewarePrivate)
	{
		private.POST("/item/subject/create", subjectApi.Create)
		private.POST("/item/subject/update", subjectApi.Update)
		private.POST("/item/subject/deletes", subjectApi.Deletes)
		private.POST("/item/subject/getList", subjectApi.GetList)
		// private.POST("/item/event/updateByEventId", eventApi.UpdateByEventId)
		// private.POST("/item/event/getListByTimeScope", eventApi.GetListByTimeScope)
		// private.POST("/item/event/deleteByIdAndItemId", eventApi.DeleteByIdAndItemId)
		// private.POST("/item/event/getInfo", eventApi.GetInfo)
	}

}
