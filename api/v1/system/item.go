package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/api/v1/common/base"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/mail"
	"calendar-note-gin/lib/systemSetting"
	"calendar-note-gin/models"
	"calendar-note-gin/models/datatype"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var StyleDefault = ""

type ParamItemVisit struct {
	ParamItemGetID
	VisitToken string
}

type ParamItemGetID struct {
	ItemId uint
}

type ParamCheckVisitPassword struct {
	ParamItemGetID
	Password string
}

type ParamItemSetting struct {
	ParamItemGetID
	Title       string
	Description string
	Password    string
	datatype.ItemStyleConfig
}
type ItemApi struct {
}

func (a *ItemApi) Create(c *gin.Context) {
	UserInfo, _ := base.GetCurrentUserInfo(c)
	param := ParamItemSetting{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
	}

	// 默认创建项目的设置
	defaultStyleConfig := datatype.ItemStyleConfig{
		WeekStartDay: 7,
		// MonthViewMaxEvent:                "",
		// WeekNumbers:                      0,
		DayTimeDisplayMode: 2,
		// SpecialDayHolidayBackgroundColor: "",
		// SpecialDaySpecialID: 0,
		// SpecialDayTextDisplay: false,
	}

	interfaceData := models.Item{
		Title:       param.Title,
		Description: param.Description,
		UserId:      UserInfo.ID,
		StyleConfig: defaultStyleConfig,
	}
	if err := global.Db.Create(&interfaceData).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
	}

	apiReturn.SuccessData(c, gin.H{
		"itemId": interfaceData.ID,
	})
}

// 获取我的项目
func (a *ItemApi) GetMyItemList(c *gin.Context) {

	param := base.PageLimitVerify{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	mItem := models.Item{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	itemList, count := mItem.GetList(map[string]interface{}{
		"user_id": userInfo.ID,
	})
	list := []map[string]interface{}{}
	for _, v := range itemList {
		requirePassword := false
		if v.Password != "" {
			requirePassword = true
		}

		list = append(list, map[string]interface{}{
			"title":           v.Title,
			"description":     v.Description,
			"createTime":      v.CreatedAt.Format(cmn.TimeFormatMode4),
			"requirePassword": requirePassword,
			"itemId":          v.ID,
		})
	}
	apiReturn.SuccessListData(c, list, count)
}

func (a *ItemApi) GetPublicInfoByItemId(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)

	param := ParamItemGetID{}
	if err := c.BindJSON(&param); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	mItem := models.Item{}
	info, err := mItem.GetInfo(map[string]interface{}{
		"user_id": userInfo.ID,
		"id":      param.ItemId,
	})
	if err != nil {
		apiReturn.ErrorNoAccess(c)
		return
	}

	requirePassword := false
	if info.Password != "" {
		requirePassword = true
	}
	apiReturn.SuccessData(c, gin.H{
		"title":           info.Title,
		"description":     info.Description,
		"createTime":      info.CreatedAt.Format(cmn.TimeFormatMode4),
		"requirePassword": requirePassword,
		"itemId":          info.ID,
	})
}

// 验证访问密码密码
func (a *ItemApi) CheckVisitPassword(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	param := ParamCheckVisitPassword{}

	if err := c.BindJSON(&param); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	mItem := models.Item{}
	info, err := mItem.GetInfo(map[string]interface{}{
		"id":      param.ItemId,
		"user_id": userInfo.ID,
	})
	if err != nil || info.Password != param.Password {
		apiReturn.Error(c, global.Lang.Get("common.api_password_error"))
		return
	}

	// 返回访问token
	apiReturn.SuccessData(c, gin.H{
		"visitToken": a.visitPasswordToVisitTocken(info.Password, info.CreatedAt.Format(cmn.TimeFormatMode2)),
	})
}

// 中间件 验证是否有权限访问当前项目
func (a *ItemApi) MiddlewarePrivate(c *gin.Context) {
	param := ParamItemVisit{}
	currentUserInfo, _ := base.GetCurrentUserInfo(c)

	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		c.Abort()
		return
	}

	mItem := models.Item{}
	info, err := mItem.GetInfo(map[string]interface{}{
		"id": param.ItemId,
	})

	c.Set("itemInfo", info)

	// 项目不存在 或者用户对不上
	if err != nil || info.UserId != currentUserInfo.ID {
		apiReturn.ErrorNoAccess(c)
		// cmn.Debug("项目不存在 或者用户对不上")
		c.Abort()
		return
	}

	// 项目本身存在密码，但是token不正确
	if info.Password != "" {
		if a.visitPasswordToVisitTocken(info.Password, info.CreatedAt.Format(cmn.TimeFormatMode2)) != param.VisitToken {
			apiReturn.ErrorNoAccess(c)
			// cmn.Debug("项目本身存在密码，但是token不正确")
			c.Abort()
			return
		}
	}

	// 无需密码或者密码验证通过

}

// 访问密码转token
func (a *ItemApi) visitPasswordToVisitTocken(password string, key string) string {
	return cmn.Md5(cmn.Md5(password) + cmn.Md5(key))
}

// 获取详情
func (a ItemApi) GetPrivateInfo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)

	param := ParamItemGetID{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	mItem := models.Item{}
	info, err := mItem.GetInfo(map[string]interface{}{
		"user_id": userInfo.ID,
		"id":      param.ItemId,
	})
	if err != nil {
		apiReturn.ErrorNoAccess(c)
		return
	}

	requirePassword := false
	if info.Password != "" {
		requirePassword = true
	}
	// styleConfig := map[string]interface{}{}
	// json.Unmarshal([]byte(info.StyleConfig), &styleConfig)
	apiReturn.SuccessData(c, gin.H{
		"title":           info.Title,
		"description":     info.Description,
		"createTime":      info.CreatedAt.Format(cmn.TimeFormatMode4),
		"requirePassword": requirePassword,
		"itemId":          info.ID,
		"styleConfig":     info.StyleConfig,
	})
}

// 保存配置
func (a ItemApi) SaveConfig(c *gin.Context) {
	param := ParamItemSetting{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	styleConfig, err := json.Marshal(param.ItemStyleConfig)
	if err != nil {
		apiReturn.Error(c, "json 数据格式化错误")
		return
	}
	mItem := models.Item{}
	errUpdate := mItem.UpdateByCondition(map[string]interface{}{
		"id": param.ItemId,
	}, map[string]interface{}{
		"title":       param.Title,
		"description": param.Description,
		"password":    param.Password,

		"style_config": string(styleConfig),
	})

	if errUpdate != nil {
		// fmt.Println("修改记录出错", errUpdate)
		apiReturn.Error(c, "数据修改失败")
		return
	}
	apiReturn.Success(c)
}

// 获取配置项
func (a ItemApi) GetConfig(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	param := ParamItemGetID{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	mItem := models.Item{}
	info, err := mItem.GetInfo(map[string]interface{}{
		"user_id": userInfo.ID,
		"id":      param.ItemId,
	})
	if err != nil {
		apiReturn.ErrorNoAccess(c)
		return
	}

	// styleConfig := ParamItemStyleConfig{}
	// json.Unmarshal([]byte(info.StyleConfig), &styleConfig)
	apiReturn.SuccessData(c, gin.H{
		"title":       info.Title,
		"description": info.Description,
		"createTime":  info.CreatedAt.Format(cmn.TimeFormatMode4),
		"password":    info.Password,
		"itemId":      info.ID,
		"styleConfig": info.StyleConfig,
	})
}

// 我非常喜欢编程，当做爱好，平时在家没日没夜的敲着代码。
// 一天有人对我说，你为了什么这么拼，我想了想说，想让更多人看到自己吧，争取更多赚钱的机会
// 他说：那你就不是爱好，就是为了赚钱。
// 我强调说：才不是，就是爱好
// 他说：如果你非常有钱你还会继续坚持敲代码吗？我沉思了...
func (a *ItemApi) DeleteOne(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	param := ParamItemGetID{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	mItem := models.Item{}
	if err := mItem.DeleteByIdAndUserId(param.ItemId, userInfo.ID); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.db_error", "[", err.Error(), "]"))
	} else {
		apiReturn.Success(c)
	}
}

func (a *ItemApi) GetSummaryData(c *gin.Context) {
	// 第一个事件：
	// 最新事件
	// 事件总数
	// 主题数

	param := ParamItemGetID{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}
	firstEvent := models.Event{}
	lastEvent := models.Event{}
	var (
		eventCount          int64
		subjectCount        int64
		firstEventStartTime string
		lastEventStartTime  string
	)

	if rows := global.Db.Where("item_id=?", param.ItemId).Order("start_time asc").First(&firstEvent).RowsAffected; rows == 0 {
		firstEventStartTime = ""
	} else {
		firstEventStartTime = firstEvent.StartTime.Time.Format(cmn.TimeFormatMode4)
	}
	if rows := global.Db.Where("item_id=?", param.ItemId).Order("start_time desc").First(&lastEvent).RowsAffected; rows == 0 {
		lastEventStartTime = ""
	} else {
		lastEventStartTime = lastEvent.StartTime.Time.Format(cmn.TimeFormatMode4)
	}

	global.Db.Model(&models.Event{}).Where("item_id=?", param.ItemId).Order("start_time desc").Count(&eventCount)
	global.Db.Model(&models.Subject{}).Where("item_id=?", param.ItemId).Order("start_time desc").Count(&subjectCount)

	apiReturn.SuccessData(c, gin.H{
		"firstEventStartTime": firstEventStartTime,
		"lastEventStartTime":  lastEventStartTime,
		"eventCount":          eventCount,
		"subjectCount":        subjectCount,
	})

}

// 找回密码，将密码发送至邮箱
func (a *ItemApi) ForgotVisitPassword(c *gin.Context) {
	param := ParamItemGetID{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mItem := models.Item{}

	visitPassword := ""
	itemName := ""
	itemInfo, err := mItem.GetInfo(map[string]interface{}{"id": param.ItemId})
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	if itemInfo.UserId != userInfo.ID {
		apiReturn.ErrorNoAccess(c)
		return
	}
	visitPassword = itemInfo.Password
	itemName = itemInfo.Title

	emailInfoConfig := systemSetting.Email{}
	systemSetting.GetValueByInterface("system_email", &emailInfoConfig)
	emailInfo := mail.EmailInfo{
		Username: emailInfoConfig.Mail,
		Password: emailInfoConfig.Password,
		Host:     emailInfoConfig.Host,
		Port:     emailInfoConfig.Port,
	}
	emailer := mail.NewEmailer(emailInfo)
	title := global.Lang.Get("mail.item_retrieval_password_title")
	content := global.Lang.GetWithFields("mail.item_retrieval_password_content", map[string]string{
		"ItemName": itemName,
	})
	if err = emailer.SendMailOfVCode(userInfo.Mail, title, content, visitPassword); err != nil {
		global.Logger.Errorf("[item retrieval password] failed to send email to %s , err:%+v\n", userInfo.Mail, err)
		apiReturn.Error(c, global.Lang.Get("common.contact_admin"))
		return
	}
	global.Logger.Infof("[item retrieval password] send email to %s ", userInfo.Mail)
	apiReturn.Success(c)
}
