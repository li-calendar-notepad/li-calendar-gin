package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/api/v1/common/base"
	"calendar-note-gin/lib/captcha"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/mail"
	"calendar-note-gin/lib/systemSetting"
	"calendar-note-gin/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginApi struct {
}

// 登录输入验证
type LoginLoginVerify struct {
	Username string `validate:"required,min=5"`
	Password string `validate:"required,min=5,max=20"`
	VCode    string `validate:"max=6"`
	Email    string `json:"email"`
}

// @Summary 登录账号
// @Accept application/json
// @Produce application/json
// @Param LoginLoginVerify body LoginLoginVerify true "登陆验证信息"
// @Tags user
// @Router /login [post]
func (l LoginApi) Login(c *gin.Context) {
	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.Error(c, errMsg)
		return
	}

	settings := systemSetting.ApplicationSetting{}
	systemSetting.GetValueByInterface("system_application", &settings)

	// 验证验证码
	if settings.Login.LoginCaptcha {
		var captchaId string
		var err error

		// 获取captchaId
		if captchaId, err = captcha.CaptchaGetIdByCookieHeader(c, "CaptchaId"); err != nil {
			apiReturn.Error(c, global.Lang.Get("login.err_captcha_check_fail"))
			return
		}

		// 验证码错误
		if !captcha.CaptchaVerifyHandle(captchaId, param.VCode) {
			apiReturn.Error(c, global.Lang.Get("captcha.api_void_captcha_id"))
			return
		}
	}

	mUser := models.User{}
	var (
		err  error
		info models.User
	)
	bToken := ""
	if info, err = mUser.GetUserInfoByUsernameAndPassword(param.Username, cmn.PasswordEncryption(param.Password)); err != nil {
		// 未找到记录 账号或密码错误
		if err == gorm.ErrRecordNotFound {
			apiReturn.Error(c, global.Lang.Get("login.err_username_password"))
			return
		} else {
			// 未知错误
			apiReturn.Error(c, err.Error())
			return
		}

	}

	// 停用或未激活
	if info.Status != 1 {
		apiReturn.Error(c, global.Lang.Get("login.err_username_deactivation"))
		return
	}

	bToken = info.Token
	if info.Token == "" {
		// 生成token
		buildTokenOver := false
		for !buildTokenOver {
			bToken = cmn.BuildRandCode(32, cmn.RAND_CODE_MODE2)
			if _, err := mUser.GetUserInfoByToken(bToken); err != nil {
				// 保存token
				mUser.UpdateUserInfoByUserId(info.ID, map[string]interface{}{
					"token": bToken,
				})
				buildTokenOver = true
			}
		}
	}

	global.UserToken.Set(bToken, info, 604800*time.Second)

	// 设置当前用户信息
	c.Set("userInfo", info)
	// 返回token等基本信息
	apiReturn.SuccessData(c, gin.H{
		"token":     bToken,
		"headImage": info.HeadImage,
		"name":      info.Name,
		"username":  info.Username,
	})
}

func (l *LoginApi) Logout(c *gin.Context) {
	// userInfo, _ := base.GetCurrentUserInfo(c)
	apiReturn.Success(c)
}

// 获取重置密码的验证码
func (l *LoginApi) SendResetPasswordVCode(c *gin.Context) {
	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	emailVCode := cmn.BuildRandCode(6, cmn.RAND_CODE_MODE2)
	global.VerifyCodeCachePool.Set(param.Email, emailVCode, 0)

	userCheck := &models.User{Mail: param.Email}
	userInfo := userCheck.GetUserInfoByMail()
	if userInfo == nil {
		apiReturn.Error(c, "账号不存在")
		return
	}
	if err := mail.SendResetPasswordVCode(param.Email, emailVCode); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	apiReturn.Success(c)

}

// 使用邮箱验证码重置密码
func (l *LoginApi) ResetPasswordByVCode(c *gin.Context) {
	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	userCheck := &models.User{Mail: param.Email}
	userInfo := userCheck.GetUserInfoByMail()
	if userInfo == nil {
		apiReturn.Error(c, "账号不存在")
		return
	}

	// 校验验证码
	{
		if vCode, ok := global.VerifyCodeCachePool.Get(param.Email); !ok || param.VCode != vCode {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			return
		}
		global.VerifyCodeCachePool.Delete(param.Email)
	}

	updateData := map[string]interface{}{
		"password": cmn.PasswordEncryption(param.Password),
		"token":    "",
	}

	if err := userInfo.UpdateUserInfoByUserId(userInfo.ID, updateData); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)

}
