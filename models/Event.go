package models

import (
	"database/sql"
)

// 事件模型
type Event struct {
	BaseModel
	Title          string       `json:"title" gorm:"type:varchar(2000)"`             // 标题
	Content        string       `json:"content" gorm:"type:text"`                    // 内容
	StartTime      sql.NullTime `json:"startTime"`                                   // 开始时间
	EndTime        sql.NullTime `json:"endTime"`                                     // 结束时间
	UserId         uint         `json:"userId"`                                      // 用户ID
	ItemId         uint         `json:"itemId"`                                      // 项目ID
	SubjectId      uint         `json:"subjectId"`                                   // 主题ID
	StyleId        uint         `json:"styleId"`                                     // 风格ID
	ExtendContent  string       `json:"extendContent" gorm:"type:text"`              // 扩展内容 json 数据（例：心情等）
	EventType      uint         `json:"eventType"`                                   // 事件类型 1.全天 2.普通事件
	ClassName      string       `json:"className" gorm:"type:varchar(50);default:0"` // 样式类名
	ReminderBefore int64        `json:"reminderBefore"`                              // 开始前提醒  0:不提醒 (单位:分钟)
	User           User
	Item           Item
	Subject        Subject
	Style          Style
}

//
func (m *Event) UpdateByCondition(condition map[string]interface{}, update map[string]interface{}) error {
	return Db.Model(&Event{}).Where(condition).Updates(update).Error
}

// 创建
func (m *Event) Create(event Event) (Event, error) {
	err := Db.Create(&event).Error
	return event, err
}

// 獲取
func (m *Event) GetEventById(id uint) (event Event, err error) {
	err = Db.Where("id=?", id).First(&event).Error
	return event, err
}

// 根据ID删除
func (m *Event) DeleteByIdAndItemId(id, itemId uint) error {
	return Db.Delete(&Event{}, "id=? AND item_id=?", id, itemId).Error
}

func (m *Event) GetListByCondition(condition map[string]interface{}) (list []Event, count int64, err error) {
	Db.Find(&list, condition)
	return
}
