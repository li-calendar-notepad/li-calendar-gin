package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitStyleCssRouter(router *gin.RouterGroup) {
	styleApi := v1.ApiGroupApp.StyleApi
	{
		router.POST("style/getStyleList", styleApi.GetStyleList)
		router.GET("style/getStyleCssCode.css", styleApi.GetStyleCssCode) // 样式文件返回
	}

}
