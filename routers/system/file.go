package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitFileRouter(router *gin.RouterGroup) {
	FileApi := v1.ApiGroupApp.FileApi

	// 验证项目的权限(有访问密码的需要验证访问token)
	private := router.Group("")
	{
		private.POST("/file/uploadImg", FileApi.UploadImg)
		private.POST("/file/uploadFiles", FileApi.UploadFiles)
	}

}
