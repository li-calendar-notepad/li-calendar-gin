package models

// 事件风格表
type Style struct {
	BaseModel
	Title           string `gorm:"type:varchar(50)" json:"title"`            // 标题
	Sort            int    `gorm:"type:int(11)" json:"sort"`                 // 排序
	ClassName       string `gorm:"index:;type:varchar(50)" json:"className"` // 类名称
	TextColor       string `gorm:"type:varchar(50)" json:"textColor"`        // 字体颜色
	BackgroundColor string `gorm:"type:varchar(50)" json:"backgroundColor"`  // 背景颜色
	BorderColor     string `gorm:"type:varchar(50)" json:"borderColor"`      // 边框颜色
}
