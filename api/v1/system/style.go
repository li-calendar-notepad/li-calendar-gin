package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/jsonConfig"
	"calendar-note-gin/models"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type StyleApi struct {
}

type ParamStyle struct {
	models.Style
	StyleId uint `json:"styleId"`
}

type ParamStyleIdsArr struct {
	StyleIds []uint `json:"styleIds"`
}

func (l StyleApi) GetStyleList(c *gin.Context) {
	styleList := []models.Style{}
	var count int64
	global.Db.Order("sort Desc").Find(&styleList).Count(&count)
	list := []map[string]interface{}{}
	for _, v := range styleList {
		list = append(list, gin.H{
			"title":           v.Title,
			"sort":            v.Sort,
			"className":       v.ClassName,
			"textColor":       v.TextColor,
			"backgroundColor": v.BackgroundColor,
			"borderColor":     v.BorderColor,
			"styleId":         v.ID,
		})
	}
	apiReturn.SuccessListData(c, list, count)
}

// 编辑
func (l StyleApi) Edit(c *gin.Context) {
	param := ParamStyle{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 查询是否有重复的名字或者类名
	findRes := models.Style{}
	if rowsAffected := global.Db.Where("title=?", param.Title).First(&findRes).RowsAffected; rowsAffected != 0 {
		// 如果是修改则判定ID是否一致
		if param.StyleId != 0 {
			if findRes.ID != param.StyleId {
				apiReturn.Error(c, global.Lang.Get("style.err_title_already_exists"))
				return
			}
		} else {
			apiReturn.Error(c, global.Lang.Get("style.err_title_already_exists"))
			return
		}
	}

	if rowsAffected := global.Db.Where("class_name=?", param.ClassName).First(&findRes).RowsAffected; rowsAffected != 0 {
		// 如果是修改则判定ID是否一致
		if param.StyleId != 0 {
			if findRes.ID != param.StyleId {
				apiReturn.Error(c, global.Lang.Get("style.err_class_name_already_exists"))
				return
			}
		} else {
			apiReturn.Error(c, global.Lang.Get("style.err_class_name_already_exists"))
			return
		}
	}
	var styleId uint
	allowField := []string{"Title", "Sort", "ClassName", "TextColor", "BackgroundColor", "BorderColor"}
	if param.StyleId != 0 {
		// 编辑
		global.Db.Select(allowField).Where("id=?", param.StyleId).Updates(&param.Style)
		styleId = param.StyleId
	} else {
		// 添加
		if err := global.Db.Select(allowField).Create(&param.Style).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
		}
		styleId = param.Style.ID
	}
	apiReturn.SuccessData(c, gin.H{
		"styleId": styleId,
	})
}

func (l StyleApi) Deletes(c *gin.Context) {

	param := ParamStyleIdsArr{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	if err := global.Db.Delete(&models.Style{}, param.StyleIds).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)

}

// 获取css样式的代码
func (l StyleApi) GetStyleCssCode(c *gin.Context) {
	cssCode := ""
	// var cssCode strings.Builder
	styleList := []models.Style{}
	global.Db.Order("sort Desc").Find(&styleList)

	for _, v := range styleList {
		cssCode += "." + v.ClassName + "{background-color:" + v.BackgroundColor + "!important;border:1px solid " + v.BorderColor + "!important;}"
		cssCode += "." + v.ClassName + " .fc-event-title{color:" + v.TextColor + "!important;}"
	}
	// c.String(200,)
	c.String(200, cssCode)
}

func (l StyleApi) EventStyleExport(c *gin.Context) {
	type ExportParam struct {
		ParamStyleIdsArr
		FileName string `json:"fileName"`
	}
	param := ExportParam{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	datas := []models.Style{}
	if err := global.Db.Where("id in ?", param.StyleIds).Find(&datas).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	elementStyleDatas := []jsonConfig.EventStyleDataModel{}
	for _, v := range datas {
		elementStyleDatas = append(elementStyleDatas, jsonConfig.EventStyleDataModel{
			Title:           v.Title,
			ClassName:       v.ClassName,
			TextColor:       v.TextColor,
			BackgroundColor: v.BackgroundColor,
			BorderColor:     v.BorderColor,
		})
	}

	config := jsonConfig.NewConfigModel(jsonConfig.ABILITY_MODE_EVENT_STYLE, "1")
	config.Data = elementStyleDatas
	byte, _ := jsonConfig.BuildExportFile(config)
	jsonConfig.Write(c, param.FileName, byte)
}

// 导入文件
func (l StyleApi) EventStyleImport(c *gin.Context) {

	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.Error(c, "上传失败")
		return
	}
	src, err := f.Open()
	defer src.Close()
	if err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	contentByte, err := ioutil.ReadAll(src)
	if err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	configFile := jsonConfig.EventStyleModel{}
	if err := json.Unmarshal(contentByte, &configFile); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	if !jsonConfig.ConfigModelCheck(&configFile.ConfigModel, jsonConfig.ABILITY_MODE_EVENT_STYLE, "1") {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}
	allowField := []string{"Title", "Sort", "ClassName", "TextColor", "BackgroundColor", "BorderColor"}
	for _, v := range configFile.Data {
		cteateData := models.Style{
			TextColor:       v.TextColor,
			BackgroundColor: v.BackgroundColor,
			ClassName:       v.ClassName,
			Title:           v.Title,
			BorderColor:     v.BorderColor,
		}
		// // 查询是否有重复的名字或者类名
		// {
		// 	findRes := models.Style{}
		// 	var rowsAffected int64 = 1
		// 	for rowsAffected != 0 {
		// 		rowsAffected = global.Db.Debug().Where("title=?", v.Title).First(&findRes).RowsAffected
		// 		if rowsAffected != 0 {
		// 			cteateData.Title += "_1"
		// 			v.Title = cteateData.Title
		// 		}

		// 	}
		// }
		// {
		// 	findRes := models.Style{}
		// 	var rowsAffected int64 = 1
		// 	for rowsAffected != 0 {
		// 		rowsAffected = global.Db.Debug().Where("class_name=?", v.ClassName).First(&findRes).RowsAffected
		// 		if rowsAffected != 0 {
		// 			cteateData.ClassName += "_1"
		// 			v.ClassName = cteateData.ClassName
		// 		}

		// 	}
		// }

		// 添加
		if err := global.Db.Select(allowField).Create(&cteateData).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
		}
	}

	apiReturn.Success(c)

}
