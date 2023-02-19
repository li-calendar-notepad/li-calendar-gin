package admin

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/api/v1/common/base"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UsersApi struct {
}

// type ParamUserInfo struct {
// 	UserId    uint   `json:"userId"`
// 	Username  string `json:"username" validate:"required,email"`
// 	Password  string `json:"password" validate:"required"`
// 	Name      string `json:"name" `
// 	HeadImage string `json:"headImage" `
// 	Status    int    `json:"status" `
// 	Role      int    `json:"role" `
// 	Mail      string `json:"mail" `
// }

func (a UsersApi) Create(c *gin.Context) {
	param := models.User{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	param.Password = "-"
	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	mUser := models.User{
		Username:  param.Username,
		Password:  cmn.PasswordEncryption(param.Password),
		Name:      param.Name,
		HeadImage: param.HeadImage,
		Status:    param.Status,
		Role:      param.Role,
		Mail:      param.Username,
	}

	// 验证账号是否存在
	if _, err := mUser.CheckUsernameExist(param.Username); err != nil {
		apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
		return
	}

	userInfo, err := mUser.CreateOne()

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"userId": userInfo.ID})
}

func (a UsersApi) Deletes(c *gin.Context) {
	type UserIds struct {
		UserIds []uint
	}
	param := UserIds{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	if err := global.Db.Delete(&models.User{}, &param.UserIds).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

func (a UsersApi) Update(c *gin.Context) {
	param := models.User{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}
	param.Password = "-"        // 修改不允许修改密码，为了验证通过
	param.Mail = param.Username // 密码邮箱同时修改
	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	allowField := []string{"Username", "Name", "HeadImage", "Status", "Role", "Mail", "Token"}
	mUser := models.User{}

	// 验证账号是否存在
	if userInfo, err := mUser.CheckUsernameExist(param.Username); err != nil {
		if userInfo.ID != param.UserId {
			apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
			return
		}
	}

	if err := global.Db.Select(allowField).Where("id=?", param.UserId).Updates(&param).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	// 返回token等基本信息
	apiReturn.Success(c)
}

func (a UsersApi) GetList(c *gin.Context) {

	type ParamsStruct struct {
		models.User
		Limit int
		Page  int
	}

	param := ParamsStruct{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	var (
		list  []models.User
		count int64
	)

	if err := global.Db.Limit(param.Limit).Offset((param.Page - 1) * param.Limit).Find(&list).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	resMap := []map[string]interface{}{}
	for _, v := range list {
		resMap = append(resMap, map[string]interface{}{
			"userId":    v.ID,
			"name":      v.Name,
			"headImage": v.HeadImage,
			"status":    v.Status,
			"role":      v.Role,
			"username":  v.Username,
		})
	}

	apiReturn.SuccessListData(c, resMap, count)
}
