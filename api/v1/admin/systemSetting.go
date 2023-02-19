package admin

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/systemSetting"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SystemSettingApi struct {
}

// ==========
// 系统邮箱
// ==========

func (a *SystemSettingApi) SetEmail(c *gin.Context) {
	param := systemSetting.Email{}
	if err := c.ShouldBindWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		return
	}

	if err := systemSetting.Set("system_email", param); err != nil {
		apiReturn.Error(c, global.Lang.Get("admin.api_mail_setting_save_err"))
	}
	apiReturn.Success(c)
}

func (a *SystemSettingApi) GetEmail(c *gin.Context) {
	info := systemSetting.Email{}
	if err := systemSetting.GetValueByInterface("system_email", &info); err != nil && err != systemSetting.ErrorNoExists {
		apiReturn.Error(c, global.Lang.Get("admin.api_mail_setting_get_err"))
		return
	}
	apiReturn.SuccessData(c, info)
}

// =========
// 应用设置
// =========

func (a *SystemSettingApi) SetApplicationSetting(c *gin.Context) {
	param := systemSetting.ApplicationSetting{}
	if err := c.ShouldBindWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		return
	}

	if err := systemSetting.Set("system_application", param); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_save_fail"))
	}
	apiReturn.Success(c)
}

func (a *SystemSettingApi) GetApplicationSetting(c *gin.Context) {
	info := systemSetting.ApplicationSetting{}
	if err := systemSetting.GetValueByInterface("system_application", &info); err != nil && err != systemSetting.ErrorNoExists {
		apiReturn.Error(c, global.Lang.Get("common.api_get_fail"))
		return
	}
	apiReturn.SuccessData(c, info)
}
