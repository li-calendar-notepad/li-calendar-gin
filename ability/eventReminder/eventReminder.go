package eventReminder

import (
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"
	"time"
)

// 事件定时提醒（事件精确到分）
// 每分钟去库中查找一次当前时间需要提醒的事件，
// 然后进行线程提醒

// var ticker *time.Ticker

type EventReminder struct {
	Ticker *time.Ticker
}

func (e *EventReminder) Start() {
	go func() {
		e.Ticker = time.NewTicker(60 * time.Second)
		defer e.Ticker.Stop()

		for e.Ticker != nil {
			select {
			case <-e.Ticker.C:
				global.Logger.Debug("定时提醒任务执行中")
				runTask()
			}
		}
	}()

}

// 停止定时器任务
func (e *EventReminder) Stop() {
	e.Ticker.Stop()
	e.Ticker = nil
}

// 运行任务
func runTask() bool {
	// 查询数据库，
	mEventReminder := models.EventReminder{}
	eventReminderList, err := mEventReminder.GetListByReminderTime(time.Now().Format(cmn.TIME_MODE_REMINDER_TIME))
	if err != nil {
		return false
	}

	for _, v := range eventReminderList {
		// 判断任务的方式，Method
		if v.Method == 1 {
			// fmt.Println("仅执行一次的任务")
			global.Logger.Debug("定时提醒任务执行", "任务id:", v.ID, "事件id:", v.EventId)
		}
	}

	return true
}
