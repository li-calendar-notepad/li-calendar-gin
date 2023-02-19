package models

// 事件主题
type Subject struct {
	BaseModel
	Title     string `gorm:"type:varchar(255)" json:"title"`    // 主题名称
	Content   string `gorm:"type:text" json:"content"`          // 源数据 json
	ClassName string `json:"className" gorm:"type:varchar(50)"` // 样式类名
	ItemId    uint   `json:"itemId"`

	Item Item

	SubjectId uint `gorm:"-"`
}
