package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// var ApiReturn = apiReturn.ApiReturn

type EventApi struct {
}

type EventCreateOne struct {
	models.Event
}

type ParamEventCreateData struct {
	EventId   uint
	ItemId    uint
	Title     string
	Content   string
	ClassName string
	StartTime string
	EndTime   string
	Color     string
}

type ParamEventGetList struct {
	ItemId    uint
	StartTime string
	EndTime   string
}

func (a *EventApi) UpdateByEventId(c *gin.Context) {
	param := ParamEventCreateData{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		c.Abort()
		return
	}

	startTime := a.strToSqlNullTime(param.StartTime)
	endTime := a.strToSqlNullTime(param.EndTime)

	mEvent := models.Event{}
	err := mEvent.UpdateByCondition(map[string]interface{}{
		"id": param.EventId,
	}, map[string]interface{}{
		"title":      param.Title,
		"item_id":    param.ItemId,
		"class_name": param.ClassName,
		// "color":      param.Color,
		"content":    param.Content,
		"start_time": startTime,
		"end_time":   endTime,
	})

	if err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_update_fail")+err.Error())
	} else {
		apiReturn.Success(c)
	}

}

func (a *EventApi) GetListByTimeScope(c *gin.Context) {
	param := ParamEventGetList{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		return
	}
	events := []models.Event{}
	var count int64
	err := global.Db.Where("item_id=?", param.ItemId).Find(&events, "start_time >= ? AND start_time <= ?", param.StartTime, param.EndTime).Count(&count).Error
	if err != nil {
		apiReturn.Error(c, global.Lang.Get("common.db_error")+err.Error())
		return
	}
	list := []map[string]interface{}{}
	for _, v := range events {
		list = append(list, map[string]interface{}{
			"eventId":    v.ID,
			"title":      v.Title,
			"content":    v.Content,
			"className":  v.ClassName,
			"startTime":  v.StartTime.Time.Format(cmn.TimeFormatMode1),
			"endTime":    v.EndTime.Time.Format(cmn.TimeFormatMode1),
			"createTime": v.CreatedAt.Format(cmn.TimeFormatMode1),
		})
	}

	// 整理格式化
	apiReturn.SuccessListData(c, list, count)
}

// 获取事件列表和特殊的日期列表
func (a *EventApi) GetListAndSpecialDayByTimeScope(c *gin.Context) {
	param := ParamEventGetList{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		return
	}

	// 事件数据
	events := []models.Event{}
	eventsList := []map[string]interface{}{}
	{
		var count int64
		// 前后各多查10天
		startTime, _ := cmn.StrToTime(cmn.TimeFormatMode1, param.StartTime)
		endTime, _ := cmn.StrToTime(cmn.TimeFormatMode1, param.EndTime)
		startTimeStr := startTime.Add(-10 * 24 * time.Hour).Format(cmn.TimeFormatMode1)
		endTimeStr := endTime.Add(10 * 24 * time.Hour).Format(cmn.TimeFormatMode1)
		// fmt.Println("事件查询时间范围", startTimeStr, endTimeStr)
		err := global.Db.Where("item_id=?", param.ItemId).Find(&events, "start_time >= ? AND start_time <= ?", startTimeStr, endTimeStr).Count(&count).Error
		if err != nil {
			apiReturn.Error(c, global.Lang.Get("common.db_error")+err.Error())
			return
		}

		for _, v := range events {
			eventsList = append(eventsList, map[string]interface{}{
				"eventId":    v.ID,
				"title":      v.Title,
				"content":    v.Content,
				"className":  v.ClassName,
				"startTime":  v.StartTime.Time.Format(cmn.TimeFormatMode1),
				"endTime":    v.EndTime.Time.Format(cmn.TimeFormatMode1),
				"createTime": v.CreatedAt.Format(cmn.TimeFormatMode1),
			})
		}
	}

	// 特殊的日期
	specialDays := []models.SpecialDay{}
	specialDayList := []map[string]interface{}{}
	{
		var itemInfo models.Item
		itemInfoAny, exists := c.Get("itemInfo")
		if exists {
			if v, ok := itemInfoAny.(models.Item); ok {
				itemInfo = v
			}
		}

		if err := global.Db.Where("special_id=?", itemInfo.StyleConfig.SpecialDaySpecialID).Find(&specialDays, "start_time >= ? AND start_time <= ?", param.StartTime, param.EndTime).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
		}

		for _, v := range specialDays {
			specialDayList = append(specialDayList, map[string]interface{}{
				"note":      v.Note,
				"type":      v.Type,
				"startTime": v.StartTime.Time.Format("2006-01-02"),
				"endTime":   v.EndTime.Time.Format("2006-01-02"),
			})
		}
	}

	// 整理格式化
	apiReturn.SuccessData(c, gin.H{
		"events":      eventsList,
		"specialDays": specialDayList,
	})
}

func (a *EventApi) GetInfo(c *gin.Context) {
	param := ParamEventCreateData{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		return
	}
	eventInfo := models.Event{}
	err := global.Db.Where("item_id=? AND id=?", param.ItemId, param.EventId).First(&eventInfo).Error
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	// 整理格式化
	apiReturn.SuccessData(c, gin.H{
		"eventId":    eventInfo.ID,
		"title":      eventInfo.Title,
		"content":    eventInfo.Content,
		"className":  eventInfo.ClassName,
		"startTime":  eventInfo.StartTime.Time.Format(cmn.TimeFormatMode1),
		"endTime":    eventInfo.EndTime.Time.Format(cmn.TimeFormatMode1),
		"createTime": eventInfo.CreatedAt.Format(cmn.TimeFormatMode1),
	})
}

func (a *EventApi) DeleteByIdAndItemId(c *gin.Context) {
	param := ParamEventCreateData{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		c.Abort()
		return
	}

	mEvent := models.Event{}
	if err := mEvent.DeleteByIdAndItemId(param.EventId, param.ItemId); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.db_error", "[", err.Error(), "]"))
	} else {
		apiReturn.Success(c)
	}

}

// 创建一条
func (a *EventApi) CreateOne(c *gin.Context) {
	param := ParamEventCreateData{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	startTime := a.strToSqlNullTime(param.StartTime)
	endTime := a.strToSqlNullTime(param.EndTime)

	mEvent := models.Event{}
	res, err := mEvent.Create(models.Event{
		Title:     param.Title,
		ItemId:    param.ItemId,
		ClassName: param.ClassName,
		Content:   param.Content,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_create_fail")+err.Error())
	} else {
		apiReturn.SuccessData(c, gin.H{
			"eventId": res.ID,
		})
	}

}

// 关键字查找事件
func (a EventApi) SearchWord(c *gin.Context) {
	type SearchParam struct {
		Word   string
		ItemId uint
	}
	param := SearchParam{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	cond := "%" + param.Word + "%"
	res := []models.Event{}
	err := global.Db.Model(&models.Event{}).
		Where("title like ? OR content like ?", cond, cond).
		Where("item_id=?", param.ItemId).
		Order("start_time DESC").
		Find(&res).Error

	list := []map[string]interface{}{}
	for _, v := range res {
		list = append(list, map[string]interface{}{
			"title":      v.Title,
			"content":    v.Content,
			"startTime":  v.StartTime.Time.Format(cmn.TimeFormatMode1),
			"endTime":    v.EndTime.Time.Format(cmn.TimeFormatMode1),
			"createTime": v.CreatedAt.Format(cmn.TimeFormatMode1),
			"eventId":    v.ID,
		})
	}
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
	} else {
		apiReturn.SuccessListData(c, list, int64(len(res)))
	}

}

// 字符串转sql.NullTime
func (a *EventApi) strToSqlNullTime(timeStr string) (sqlNullTime sql.NullTime) {
	if timeStr != "" {
		if theTime, err := cmn.StrToTime(cmn.TimeFormatMode1, timeStr); err == nil {
			sqlNullTime.Time = theTime
			sqlNullTime.Valid = true
		}
	}
	return
}
