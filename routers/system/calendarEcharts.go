package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitCalendarEcharts(router *gin.RouterGroup) {
	itemApi := v1.ApiGroupApp.ItemApi
	calendarEchartsApi := v1.ApiGroupApp.CalendarEchartsApi
	r := router.Group("item").Use(itemApi.MiddlewarePrivate)

	r.POST("lineTopicCount", calendarEchartsApi.LineTopicCount)
	r.POST("pieChartTopic", calendarEchartsApi.PieChartTopic)
}
