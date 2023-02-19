package system

import (
	"calendar-note-gin/lib/captcha"
	"calendar-note-gin/lib/cmn"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

type CaptchaApi struct {
	ErrMsg string // 错误信息
}

// 获取图像
func (c *CaptchaApi) GetImage(ctx *gin.Context) {
	key := cmn.BuildRandCode(16, cmn.RAND_CODE_MODE2)
	// 设置网页验证码的cookie
	ctx.SetCookie("CaptchaId", key, 3600, "/", "", false, false)
	base64Str := captcha.GenerateCaptchaHandler(key)
	_ = base64Str
	// base64 字符串一般会包含头部 data:image/xxx;base64, 需要去除
	baseImg, _ := base64.StdEncoding.DecodeString(base64Str[22:])
	_, _ = ctx.Writer.WriteString(string(baseImg))
}

func (c *CaptchaApi) CheckVCode(id, vcode string) {
	// Captcha.Store = base64Captcha.DefaultMemStore
	// if store.Verify(id, vcode, true) {
	// 	body = map[string]interface{}{"code": 1001, "msg": "ok"}
	// }
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// json.NewEncoder(w).Encode(body)
}
