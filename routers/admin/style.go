package admin

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitStyleRouter(router *gin.RouterGroup) {
	styleApi := v1.ApiGroupApp.StyleApi
	{
		router.POST("style/getStyleList", styleApi.GetStyleList)
		router.POST("style/edit", styleApi.Edit)
		router.POST("style/deletes", styleApi.Deletes)

		router.POST("style/eventStyleExport", styleApi.EventStyleExport)
		router.POST("style/eventStyleImport", styleApi.EventStyleImport)

	}

}
