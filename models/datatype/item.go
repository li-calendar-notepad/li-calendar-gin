package datatype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// 项目配置
type ItemStyleConfig struct {
	WeekStartDay                     int    `json:"weekStartDay"`
	MonthViewMaxEvent                string `json:"monthViewMaxEvent"`
	WeekNumbers                      int    `json:"weekNumbers"`
	DayTimeDisplayMode               int    `json:"dayTimeDisplayMode"`
	SpecialDayHolidayBackgroundColor string `json:"specialDayHolidayBackgroundColor"` // 特殊的日期假期背景色
	SpecialDaySpecialID              int    `json:"specialDaySpecialId"`              // 特殊的日期ID
	SpecialDayTextDisplay            bool   `json:"specialDayTextDisplay"`            // 特殊的日期文字显示
}

// 查询的时候解析
func (j *ItemStyleConfig) Scan(value interface{}) error {
	var bytes []byte
	// bytes, ok := value.([]byte)
	if reflect.TypeOf(value).Kind() == reflect.String {
		if s, ok := value.(string); ok {
			bytes = []byte(s)
		} else {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}
	} else {
		if b, ok := value.([]byte); !ok {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		} else {
			bytes = b
		}
	}

	err := json.Unmarshal(bytes, j)
	// 解决空值报错的问题
	if err != nil {
		j = &ItemStyleConfig{}
	}
	return nil
}

// 保存时的编译
func (j ItemStyleConfig) Value() (driver.Value, error) {
	str, err := json.Marshal(j)
	if err != nil {
		return string(str), err
	}
	return string(str), nil
}
