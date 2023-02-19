package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/systemSetting"

	"github.com/gin-gonic/gin"
)

type Open struct {
}

func (l *Open) GetOpenConfig(c *gin.Context) {
	systemApplication := systemSetting.ApplicationSetting{}
	if err := systemSetting.GetValueByInterface("system_application", &systemApplication); err != nil {
		if err != systemSetting.ErrorNoExists {
			apiReturn.Error(c, "")
		}
	}
	apiReturn.SuccessData(c, systemApplication)
}
