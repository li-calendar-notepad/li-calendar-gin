package eventReminder

import (
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"
	"strconv"
	"time"
)

// 事件定时提醒（事件精确到分）
// 每分钟去库中查找一次当前时间需要提醒的事件，
// 然后进行线程提醒

// var lock sync.Mutex                        // 互斥锁+
var ReminderJobs = make(chan MailInfo, 60) // 工作管道的数量

type MailInfo struct {
	Email   string `json:"email"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
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

// 工作池
func SendMailWorker(id int, jobs <-chan MailInfo, results chan<- int) {
	for j := range jobs {
		// lock.Lock()
		time.Sleep(5 * time.Second)
		global.Logger.Debug("进程:", id, "模拟发送邮件-", j.Email, j.Title)
		// lock.Unlock()
		results <- 1

	}
}

// 运行任务
func runTask() bool {
	// 查询数据库，
	mEventReminder := models.EventReminder{}
	eventReminderList, err := mEventReminder.GetListByReminderTime(time.Now().Format(cmn.TIME_MODE_REMINDER_TIME))
	if err != nil {
		return false
	}

	taskTitle := time.Now().Format(cmn.TimeFormatMode1)
	startTime := time.Now().Unix() // 记录开始时间
	global.Logger.Debug("== 定时提醒任务执行开始", taskTitle)

	count := len(eventReminderList)
	// jobs := make(chan MailInfo, count) // 迁移到全局
	results := make(chan int, count)

	// 开启5个线程
	for w := 1; w <= 5; w++ {
		go SendMailWorker(w, ReminderJobs, results)
	}

	for i, v := range eventReminderList {
		// 判断任务的方式，Method
		if v.Method == 1 {
			// fmt.Println("仅执行一次的任务")
			// global.Logger.Debug("定时提醒任务执行", "任务id:", v.ID, "事件id:", v.EventId)
			m := MailInfo{
				Email: "任务id:" + strconv.Itoa(int(v.ID)),
				Title: "事件id:" + strconv.Itoa(int(v.EventId)),
				// Content: "内容是" + strconv.Itoa(i),
			}

			global.Logger.Debug("发送", taskTitle, "  ", count, " -- ", i)

			// 发送邮件
			ReminderJobs <- m
		}
	}

	// 收集结果
	for a := 1; a <= count; a++ {
		<-results
		if a == count {
			global.Logger.Debug("== 定时提醒任务执行结束", taskTitle, " 用时:", time.Now().Unix()-startTime)
		} else {
			global.Logger.Debug("执行结果", taskTitle, "  ", count, " -- ", a)
		}
	}
	// close(results)

	return true
}
