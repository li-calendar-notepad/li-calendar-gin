package models

// 分享
type Share struct {
	BaseModel
	Title      string `gorm:"type:varchar(50)" json:"title"`          // 标题
	Password   string `gorm:"type:varchar(50)" json:"password"`       // 密码
	OnlyTag    string `gorm:"type:varchar(10);index:" json:"onlyTag"` // 标识
	ItemId     uint   `json:"itemId"`                                 // 项目ID
	UserId     uint   `json:"userId"`                                 // 分享人
	Type       int    `gorm:"type:tiny(1)" json:"type"`               // 分享类型 1.项目 2.时间范围 3.单个事件
	RangeStart string `gorm:"type:varchar(20)" json:"rangeStart"`     // 范围 开始时间
	RangeEnd   string `gorm:"type:varchar(20)" json:"rangeEnd"`       // 范围 结束时间
	Auth       int    `gorm:"type:tiny(2)" json:"auth"`               // 授权 1.只读 2.读写（未启用）+

	User User
	Item Item
}
