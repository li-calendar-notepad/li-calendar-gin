package models

import (
	"calendar-note-gin/lib/global"

	"gorm.io/gorm"
)

// 提醒
type EventReminder struct {
	BaseModel
	EventId      uint   `json:"eventId"`                            // 事件id
	ReminderTime string `json:"time" gorm:"type:varchar(20),index"` // 提醒时间(索引)，格式：202303162126
	Method       uint   `json:"method" gorm:"type:int(2)"`          // 提醒方式 1.不重复提醒 2.每天 3.每周 4.每月 5.每年
	Status       uint   `json:"status" gorm:"type:int(1)"`          // 状态 0.待执行 1.已执行 2.已过期

	Event Event
}

// 添加一个
func (m *EventReminder) Add() error {
	return global.Db.Create(m).Error
}

// 删除一个
func (m *EventReminder) Del() {

}

// 根据事件获取某时间段的提醒任务列表
func (m *EventReminder) GetListByReminderTime(reminderTime string) (eventReminderList []EventReminder, err error) {
	err = Db.Model(&EventReminder{}).Find(&eventReminderList, "reminder_time=?", reminderTime).Error
	return
}

func (m *EventReminder) BeforeCreate(tx *gorm.DB) error {
	m.Status = 0 // 创建必须设置为0
	return nil
}
