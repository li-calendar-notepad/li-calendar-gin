package eventReminder

import (
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/mail"
	"calendar-note-gin/lib/systemSetting"
	"calendar-note-gin/models"
	"time"
)

// 事件定时提醒（事件精确到分）
// 每分钟去库中查找一次当前时间需要提醒的事件，
// 然后进行线程提醒

// var lock sync.Mutex                        // 互斥锁+
// var ReminderJobs = make(chan ReminderInfo, 1000) // 工作管道的数量

type MailInfo struct {
	Email   string `json:"email"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ReminderInfo struct {
	ItemTitle string // 项目标题
	Title     string // 事件标题
	// Content    string // 事件内容 不支持内容
	RemindTime time.Time // 提醒时间
	StartTime  string    // 开始时间
	EndTime    string    // 结束时间
	// UserInfo      models.User
	// EventInfo     models.Event
	EventReminder models.EventReminder
	Id            uint
}

// type MailWorker struct {
// 	Id      int
// 	Jobs    <-chan MailInfo
// 	Results chan<- int
// 	Wg      *sync.WaitGroup
// }

type EventReminder struct {
	Ticker       *time.Ticker
	ReminderJobs chan ReminderInfo
	ChanNum      int
	ThreadNum    int
}

// 定时提醒任务
//
//	@param chanNum int 管道容量 推荐:1000
//	@param ThreadNum int 线程数量 推荐:10
//	@return *EventReminder
func (e *EventReminder) Start(chanNum int, ThreadNum int) {
	ReminderJobs := make(chan ReminderInfo, chanNum)
	e.ChanNum = chanNum
	e.ThreadNum = ThreadNum
	e.ReminderJobs = ReminderJobs

	go func() {
		e.Ticker = time.NewTicker(60 * time.Second)
		defer e.Ticker.Stop()

		for e.Ticker != nil {
			select {
			case <-e.Ticker.C:
				currentTime := time.Now()
				global.Logger.Debug("#### 定时执行 ", currentTime.Format(cmn.TimeFormatMode1))
				go e.runTask(currentTime)
				// global.Logger.Debug("#### 定时执行结束 ", currentTime.Format(cmn.TimeFormatMode1))
			}
		}
	}()

	// 运行队列
	go e.startThread(e.ThreadNum)
}

// 启动线程
func (e *EventReminder) startThread(thread int) {
	for w := 1; w <= thread; w++ {
		go e.RunWorker(w)
	}
}

// 停止定时器任务
func (e *EventReminder) Stop() {
	e.Ticker.Stop()
	e.Ticker = nil
}

// 开始工作
func (e *EventReminder) RunWorker(workerId int) {
	for v := range e.ReminderJobs {
		global.Logger.Debug("任务开始", v.Title)
		// 目前正在工作
		// 站内消息
		// 发送邮件

		SendMailWorker(workerId, v)
		// 删除该定时
		// v.EventReminder.Del()
		// 生成下次定时
		global.Logger.Debug("任务结束", v.Title)
	}
}

// 邮件工作
func SendMailWorker(workerId int, reminderInfo ReminderInfo) {
	emailInfoConfig := systemSetting.Email{}
	systemSetting.GetValueByInterface("system_email", &emailInfoConfig)
	emailInfo := mail.EmailInfo{
		Username: emailInfoConfig.Mail,
		Password: emailInfoConfig.Password,
		Host:     emailInfoConfig.Host,
		Port:     emailInfoConfig.Port,
	}

	mail.SendEventReminder(mail.NewEmailer(emailInfo), reminderInfo.EventReminder.Event.Item.User.Mail, reminderInfo.EventReminder)
	// time.Sleep(10 * time.Second)
}

// 运行任务
func (e *EventReminder) runTask(currentTime time.Time) bool {

	// 查询数据库
	mEventReminder := models.EventReminder{}
	eventReminderList, err := mEventReminder.GetListByReminderTime(currentTime.Format(cmn.TIME_MODE_REMINDER_TIME))
	if err != nil {
		return false
	}

	taskTitle := currentTime.Format("15:04:05") + " " + currentTime.Format("04") + " "
	global.Logger.Debug("==== runTask开始:", taskTitle)
	defer global.Logger.Debug("==== runTask完成:", taskTitle)
	count := len(eventReminderList)
	if count == 0 {
		return false
	}

	// 将数据内容插入队列
	// global.Logger.Debug(taskTitle, "准备插入任务，任务数:", len(eventReminderList))
	for i, v := range eventReminderList {
		// 判断任务的方式，Method
		if v.Method == 1 {
			// fmt.Println("仅执行一次的任务")
			reminderInfo := ReminderInfo{
				ItemTitle:     v.Event.Item.Title,
				Title:         v.Event.Title,
				RemindTime:    currentTime,
				StartTime:     v.Event.StartTime.Time.Format(cmn.TimeFormatMode1),
				EndTime:       v.Event.EndTime.Time.Format(cmn.TimeFormatMode1),
				EventReminder: v,
			}
			// fmt.Println()
			global.Logger.Debug("插入任务:", count, "/", i+1, " ", reminderInfo.Title)

			// 插入任务到队列
			e.ReminderJobs <- reminderInfo
			v.Del()
		}
	}

	return true
}
