package models

import "database/sql"

// 特殊的日子 天详情
type SpecialDay struct {
	BaseModel
	Note      string       `gorm:"varchar(255)"` // 假期备注，国庆节，中秋节，串休
	Type      int          `gorm:"tinyint(2)"`   // 节假日 1.普通的日子 2.节假日
	StartTime sql.NullTime `json:"startTime"`    // 开始时间
	EndTime   sql.NullTime `json:"endTime"`      // 结束时间
	SpecialID uint

	Special Special
}
