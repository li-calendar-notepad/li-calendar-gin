package ability

import (
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/models"
	"fmt"
	"time"
)

// 事件定时提醒（事件精确到分）
// 每分钟去库中查找一次当前时间需要提醒的事件，
// 然后进行线程提醒

var ticker *time.Ticker

func StartEventReminder() {
	ticker = time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for ticker != nil {
		select {
		case <-ticker.C:
			runTask()
		}
	}
}

// 停止定时器任务
func StopEventReminder() {
	ticker.Stop()
	ticker = nil
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
			fmt.Println("仅执行一次的任务")
		}
	}

	return true
}
