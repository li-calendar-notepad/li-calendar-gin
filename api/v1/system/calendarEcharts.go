package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CalendarEchartsApi struct {
}

// 日历图表日期参数
type ParamCalendarEchartsDate struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// 饼图统计每个主题总数
func (a *CalendarEchartsApi) PieChartTopic(c *gin.Context) {
	param := ParamCalendarEchartsDate{}

	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		c.Abort()
		return
	}

	// 查询主题，对每个主题分别查询时间总数
	subjects := []models.Subject{}
	pieDatas := []map[string]interface{}{}
	global.Db.Find(&subjects)

	for _, v := range subjects {
		var count int64
		global.Db.Model(&models.Event{}).Where("title like ?", "#"+v.Title+"#").Count(&count)
		pieDatas = append(pieDatas, map[string]interface{}{
			"name":  v.Title,
			"value": count,
		})
	}
	apiReturn.SuccessData(c, pieDatas)
}
