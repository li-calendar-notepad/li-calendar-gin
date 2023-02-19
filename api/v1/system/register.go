package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/api/v1/common/base"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/mail"
	"calendar-note-gin/lib/systemSetting"
	"calendar-note-gin/models"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type registerInfo struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	Passwd    string `json:"passwd"`
	Vcode     string `json:"vcode"`
	EmailCode string `json:"emailCode"`
	//真正注册的标志
	IsCommitEmailCode bool `json:"isCommitEmailCode"`
	IsTest            bool `json:"isTest"`
}

const EmailCodeCapacity = 1000

// 获取注册验证码
func (l LoginApi) Register(c *gin.Context) {
	info := registerInfo{}
	err := c.ShouldBindJSON(&info)
	info.Email = info.UserName
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	errMsg, err := base.ValidateInputStruct(info)
	if err != nil {
		apiReturn.Error(c, errMsg)
		return
	}

	// 验证是否开启注册和后缀格式是否正确
	{
		systemSettingInfo := systemSetting.ApplicationSetting{}
		if err := systemSetting.GetValueByInterface("system_application", &systemSettingInfo); err != nil || !systemSettingInfo.Register.OpenRegister {
			apiReturn.Error(c, global.Lang.Get("register.unopened_register"))
			return
		}

		if systemSettingInfo.Register.EmailSuffix != "" && !cmn.VerifyFormat("^.*"+systemSettingInfo.Register.EmailSuffix+"$", info.Email) {
			apiReturn.Error(c, global.Lang.GetWithFields("register.emailSuffix_error", map[string]string{"EmailSuffix": systemSettingInfo.Register.EmailSuffix}))
			return
		}
	}

	// 验证邮箱是否被注册
	{
		userCheck := &models.User{Mail: info.UserName}
		if _, err := userCheck.GetUserInfoByUsername(info.UserName); err == nil && err != gorm.ErrRecordNotFound {
			apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
			return
		}
	}

	emailCode := generateEmailCode()
	if global.VerifyCodeCachePool.ItemCount() >= EmailCodeCapacity {
		global.VerifyCodeCachePool.Flush()
	}
	global.VerifyCodeCachePool.Set(info.Email, emailCode, 0)
	err = mail.SendRegisterEmail(info.Email, emailCode)
	if err != nil {
		apiReturn.Error(c, global.Lang.Get("mail.send_mail_fail"))
		global.Logger.Errorf("[register] fail to send email to%s", info.UserName)
		return
	}
	apiReturn.Success(c)
}

// 注册提交（开始注册）
func (l *LoginApi) Commit(c *gin.Context) {
	info := registerInfo{}
	err := c.ShouldBindJSON(&info)
	info.Email = info.UserName
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	errMsg, err := base.ValidateInputStruct(info)
	if err != nil {
		apiReturn.Error(c, errMsg)
		return
	}

	// 验证是否开启注册和后缀格式是否正确
	{
		systemSettingInfo := systemSetting.ApplicationSetting{}
		if err := systemSetting.GetValueByInterface("system_application", &systemSettingInfo); err != nil || !systemSettingInfo.Register.OpenRegister {
			apiReturn.Error(c, global.Lang.Get("register.unopened_register"))
			return
		}

		if systemSettingInfo.Register.EmailSuffix != "" && !cmn.VerifyFormat("^.*"+systemSettingInfo.Register.EmailSuffix+"$", info.Email) {
			apiReturn.Error(c, global.Lang.GetWithFields("register.emailSuffix_error", map[string]string{"EmailSuffix": systemSettingInfo.Register.EmailSuffix}))
			return
		}
	}

	// 验证邮箱是否被注册
	{
		userCheck := &models.User{Mail: info.UserName}
		if _, err := userCheck.GetUserInfoByUsername(info.UserName); err == nil && err != gorm.ErrRecordNotFound {
			apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
			return
		}
	}

	// 验证码验证
	{
		v, ok := global.VerifyCodeCachePool.Get(info.Email)
		if !ok {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			//验证码不存在
			return
		}
		if v.(string) != info.EmailCode {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			return
			//验证码有误
		}
	}

	//验证通过，注册
	user := &models.User{Mail: info.UserName, Username: info.UserName, Password: cmn.PasswordEncryption(info.Passwd), Status: 1}
	_, err = user.CreateOne()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	//删除旧的验证码
	global.VerifyCodeCachePool.Delete(info.Email)
	apiReturn.Success(c)
}

func generateEmailCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}
