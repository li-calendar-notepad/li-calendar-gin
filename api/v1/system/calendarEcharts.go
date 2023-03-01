package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CalendarEchartsApi struct {
}

// 日历图表日期参数
type ParamCalendarEchartsDate struct {
	ItemId    uint   `json:"itemId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type ParamLineDates struct {
	ItemId            uint     `json:"itemId"`
	Dates             []string `json:"dates"`
	StatisticalMethod string   `json:"statisticalMethod"` // 统计方式 day,month
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

// 折线图和日期的统计
func (a *CalendarEchartsApi) LineTopicCount(c *gin.Context) {
	param := ParamLineDates{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		c.Abort()
		return
	}

	mSubject := models.Subject{ItemId: param.ItemId}
	subjectList := mSubject.GetList()

	// start, _ := cmn.StrToTime(cmn.TimeFormatMode1, param.StartTime)
	// end, _ := cmn.StrToTime(cmn.TimeFormatMode1, param.EndTime)

	// start.Add(1 * 24 * time.Hour).Format(cmn.TimeFormatMode1)
	// start.AddDate(0, 0, -1)
	// currentTime := start
	// forlastTime := end.Add(1 * 24 * time.Hour)
	// for currentTime != end {
	// 	currentTime = start.Add(1 * 24 * time.Hour)
	// 	currentTimeStr = currentTime.Format(cmn.TimeFormatMode1)
	datas := map[string][]int64{}
	for _, v := range subjectList {
		var count int64
		// 循环前端给的日期
		for _, currentDate := range param.Dates {
			var endTime time.Time
			var startTimeStr = currentDate
			if param.StatisticalMethod == "month" {
				start, _ := cmn.StrToTime(cmn.TimeYYYY_mm_dd, currentDate)
				end, _ := cmn.StrToTime(cmn.TimeYYYY_mm_dd, currentDate)
				startTimeStr = start.AddDate(0, 0, -start.Day()+1).Format(cmn.TimeYYYY_mm_dd)
				endTime = end.AddDate(0, 1, -end.Day())
			} else {
				endTime, _ = cmn.StrToTime(cmn.TimeYYYY_mm_dd, currentDate)
			}
			fmt.Println("ddd", startTimeStr+" 00:00:00", endTime.Format(cmn.TimeYYYY_mm_dd)+" 23:59:59")
			global.Db.Model(&models.Event{}).
				Where("title like ?", "#"+v.Title+"#").
				Where("start_time >= ?", startTimeStr+" 00:00:00").
				Where("end_time <= ?", endTime.Format(cmn.TimeYYYY_mm_dd)+" 23:59:59").
				Count(&count)
			arrInt := datas[v.Title]
			arrInt = append(arrInt, count)
			datas[v.Title] = arrInt
		}

	}

	apiReturn.SuccessData(c, gin.H{
		"datas": datas,
	})

	// 最终数据{dates:[],datas{},topics{}}
}
