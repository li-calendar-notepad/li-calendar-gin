package captcha

import (
	"calendar-note-gin/lib/global"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var Store = base64Captcha.DefaultMemStore

func NewDriver() *base64Captcha.DriverString {
	driver := new(base64Captcha.DriverString)
	driver.Height = 44
	driver.Width = 120
	driver.NoiseCount = 0
	driver.ShowLineOptions = base64Captcha.OptionShowSlimeLine | base64Captcha.OptionShowHollowLine
	driver.Length = 4
	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
	driver.Fonts = []string{"wqy-microhei.ttc"}
	return driver
}

// 生成图形验证码
func GenerateCaptchaHandler(id string) string {
	var driver = NewDriver().ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, Store)
	_, content, answer := c.Driver.GenerateIdQuestionAnswer()

	item, _ := c.Driver.DrawCaptcha(content)
	c.Store.Set(id, answer)
	return item.EncodeB64string()
}

// 验证
func CaptchaVerifyHandle(id, vcode string) bool {
	return Store.Verify(id, vcode, true)
}

// 根据key获取验证码ID
func CaptchaGetIdByCookieHeader(c *gin.Context, key string) (captchaId string, err error) {

	captchaId, err = c.Cookie("CaptchaId")
	if err != nil {
		global.Logger.Errorf("failed to get captchaId from cookie, err:%+v\n", err)
		return captchaId, err
	}
	if captchaId == "" {
		captchaId = c.GetHeader(key)
	}
	return
}
