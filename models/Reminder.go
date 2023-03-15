package models

// 提醒
type EventReminder struct {
	BaseModel
	EventId      uint   `json:"eventId"`                            // 事件id
	ReminderTime string `json:"time" gorm:"type:varchar(20),index"` // 提醒时间(索引)
	Method       uint   `json:"method" gorm:"type:int(2)"`          // 提醒方式 1.不重复提醒 2.每天 3.每周 4.每月 5.每年

	Event Event
}

// 添加一个
func (m *Item) Add() {

}

// 删除一个
func (m *Item) Del() {

}
