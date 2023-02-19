package models

import (
	"calendar-note-gin/models/datatype"
)

// 项目表
type Item struct {
	BaseModel
	Title       string                   `json:"title" gorm:"type:varchar(50)"`        // 标题
	Description string                   `json:"description" gorm:"type:varchar(250)"` // 描述
	Sort        int                      `json:"sort" gorm:"type:int(11)"`             // 排序
	Password    string                   `json:"password" gorm:"type:varchar(16)"`     // 访问密码
	StyleConfig datatype.ItemStyleConfig `json:"styleConfig" gorm:"type:text"`         // json 风格设置
	UserId      uint                     `json:"userId" `
	User        User
}

func (m *Item) GetList(condition map[string]interface{}) (itemList []Item, count int64) {
	Db.Where(condition).Find(&itemList).Count(&count)
	return
}

func (m *Item) GetInfo(condition map[string]interface{}) (itemInfo Item, err error) {
	err = Db.Where(condition).First(&itemInfo).Error
	return
}

func (m *Item) UpdateByCondition(condition map[string]interface{}, update interface{}) (err error) {
	err = Db.Model(&Item{}).Where(condition).Updates(update).Error
	return
}

// 根据ID删除
func (m *Item) DeleteByIdAndUserId(id, userId uint) error {
	return Db.Delete(&Item{}, "id=? AND user_id=?", id, userId).Error
}

// 创建
func (m *Item) Create(itemInfo Item) (Item, error) {
	err := Db.Create(&itemInfo).Error
	return itemInfo, err
}
