package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitTestRouter(router *gin.RouterGroup) {
	testApi := v1.ApiGroupApp.TestApi
	router.POST("/test", testApi.Test)
}
