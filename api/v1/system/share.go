package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/api/v1/common/base"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ItemShareApi struct {
}

// 创建一个分享
func (a *ItemShareApi) Create(c *gin.Context) {
	param := models.Share{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 生成唯一标识
	{
		findRes := models.Share{}
		onlyTag := ""
		for onlyTag == "" {
			randCode := cmn.BuildRandCode(6, cmn.RAND_CODE_MODE1)
			row := global.Db.First(&findRes, "only_tag=?", onlyTag).RowsAffected
			if row == 0 {
				onlyTag = randCode
			}
		}
	}

	user, _ := base.GetCurrentUserInfo(c)

	param.UserId = user.ID
	param.Auth = 1

	if err := global.Db.Create(&param).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
	}
	apiReturn.Success(c)
}
