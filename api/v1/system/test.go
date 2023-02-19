package system

import (
	// . "calendar-note-gin/api/v1/common/apiReturn"

	// "calendar-note-gin/api/v1/common/apiReturn"

	"calendar-note-gin/api/v1/common/apiReturn"

	"github.com/gin-gonic/gin"
)

type TestApi struct {
}

var ApiReturn = apiReturn.ApiReturn

func (l TestApi) Test(c *gin.Context) {

	// 返回token等基本信息
	// ApiReturn.SuccessData(c, gin.H{})
}
