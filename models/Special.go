package models

import "database/sql"

// 特殊的日子，如假期，串休等
type Special struct {
	BaseModel
	Name   string `gorm:"varchar(255)"` // 假期文件名称，中国法定节假日
	OnlyId string `gorm:"varchar(255)"` // 唯一标识ID
	// Year       int          `gorm:"int(4)"`       // 年份
	// CheckCode  string       `gorm:"varchar(50)"`  // 节假日文件校验码
	UpdateTime sql.NullTime
}
