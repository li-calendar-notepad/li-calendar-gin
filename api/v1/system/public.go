package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/systemSetting"

	"github.com/gin-gonic/gin"
)

type PublicApi struct {
}

// 获取开放配置
func (a *PublicApi) GetOpenConfig(c *gin.Context) {
	value, err := systemSetting.GetValueString("other_open_register")
	// 不存在
	if err == systemSetting.ErrorNoExists {

	}

	apiReturn.SuccessData(c, gin.H{
		"other_open_register": value,
	})
}
