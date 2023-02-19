package models

// 状态定义
type StatusDefine struct {
	BaseModel
	Title  string `gorm:"type:varchar(255)"` // 状态标题
	Flag   string `gorm:"type:varchar(32)"`  // 英文标识
	ItemId uint
	Item   Item
}
