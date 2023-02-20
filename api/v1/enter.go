package v1

import (
	"calendar-note-gin/api/v1/admin"
	"calendar-note-gin/api/v1/system"
)

// 加入路由
type ApiGroup struct {
	LoginApi           system.LoginApi
	TestApi            system.TestApi
	CaptchaApi         system.CaptchaApi
	ItemApi            system.ItemApi
	EventApi           system.EventApi
	FileApi            system.FileApi
	StyleApi           system.StyleApi
	JsonConfig         system.JsonConfig
	SubjectApi         system.SubjectApi
	UserApi            system.UserApi
	Open               system.Open
	CalendarEchartsApi system.CalendarEchartsApi
	AdminApiGroup      admin.ApiGroup
	// UserApi    system.admin
}

var ApiGroupApp = new(ApiGroup)
