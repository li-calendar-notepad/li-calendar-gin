package admin

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/cmn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/lib/jsonConfig"
	"calendar-note-gin/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type SpecialDayApi struct {
}

func (a *SpecialDayApi) GetInfoDays(c *gin.Context) {
	type paramStruct struct {
		SpecialId uint `json:"specialId"`
	}
	param := paramStruct{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	list := []models.SpecialDay{}
	if err := global.Db.Where("special_id=?", param.SpecialId).Order("start_time DESC").Find(&list).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	listMap := []map[string]interface{}{}
	for _, v := range list {
		listMap = append(listMap, map[string]interface{}{
			"note":      v.Note,
			"type":      v.Type,
			"startTime": v.StartTime.Time.Format("2006-01-02"),
			"endTime":   v.EndTime.Time.Format("2006-01-02"),
		})
	}
	apiReturn.SuccessListData(c, listMap, 0)
}

// 获取列表
func (a *SpecialDayApi) GetList(c *gin.Context) {
	list := []models.Special{}
	if err := global.Db.Order("id desc").Find(&list).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	listMap := []map[string]interface{}{}
	for _, v := range list {
		listMap = append(listMap, map[string]interface{}{
			"id":         v.ID,
			"name":       v.Name,
			"onlyId":     v.OnlyId,
			"updateTime": v.UpdatedAt.Format("2006-01-02"),
		})
	}
	apiReturn.SuccessListData(c, listMap, 0)
}

// 导入
func (a *SpecialDayApi) SpecialDayImport(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.Error(c, "上传失败")
		return
	}
	src, err := f.Open()
	defer src.Close()
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	contentByte, err := ioutil.ReadAll(src)
	if err != nil {
		// fmt.Println("aaa")
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	configFile := jsonConfig.SpecialDayModel{}

	if err := json.Unmarshal(contentByte, &configFile); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if !jsonConfig.ConfigModelCheck(&configFile.ConfigModel, jsonConfig.ABILITY_MODE_SPECIAL_DAY, "1") {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format"))
		return
	}

	// checkCode := cmn.Md5(string(contentByte))
	var specialId uint
	dateFormatMode := "2006-01-02"
	// 文件更新时间
	configUpdateTime, err := cmn.StrToTime(cmn.TimeFormatMode1, configFile.Data.UpdateTime)
	if err != nil {
		fmt.Println("时间访问时", configFile.Data.UpdateTime, err)
	}
	// 基础数据
	baseSpecial := models.Special{
		Name:   configFile.Data.Name,
		OnlyId: configFile.Data.OnlyId,
		// Year:   configFile.Data.Year,
		UpdateTime: sql.NullTime{
			Time:  configUpdateTime,
			Valid: true,
		},
		// CheckCode: checkCode,
	}

	findRes := models.Special{}
	findCount := global.Db.Where("only_id=?", configFile.Data.OnlyId).First(&findRes).RowsAffected
	if findCount != 0 {
		// 更新Special表数据
		global.Db.Delete(&models.SpecialDay{}, "special_id=? AND start_time>= ? AND start_time<=?", findRes.ID, configFile.Data.StartDate, configFile.Data.EndDate)
		specialId = findRes.ID
		global.Db.Model(&models.Special{}).Where("id=?", specialId).Updates(&baseSpecial)
		// fmt.Println("更新表", err)
	} else {
		// 插入基础数据
		if err := global.Db.Create(&baseSpecial).Error; err != nil {
			// fmt.Println("插入基础数据错误", err.Error())
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
		specialId = baseSpecial.ID
	}

	// fmt.Println("开始执行 ---")
	days := []models.SpecialDay{}
	for k, v := range configFile.Data.Days {
		holiday := 1
		if v.Holiday {
			holiday = 2
		}
		startTime := sql.NullTime{}
		endTime := sql.NullTime{}

		{
			sTime, err := cmn.StrToTime(dateFormatMode, k)
			if err != nil {
				fmt.Println("时间转换错误", err.Error())
			}
			endTime = sql.NullTime{
				Time:  sTime.Add(24 * time.Hour),
				Valid: true,
			}
			startTime = sql.NullTime{
				Time:  sTime,
				Valid: true,
			}
		}

		days = append(days, models.SpecialDay{
			Note:      v.Name,
			Type:      holiday,
			StartTime: startTime,
			EndTime:   endTime,
			SpecialID: specialId,
		})

	}
	insertCount := global.Db.Create(&days).RowsAffected
	apiReturn.SuccessData(c, gin.H{
		"insertCount": insertCount,
	})
}

// 删除一个
func (a *SpecialDayApi) DeleteOne(c *gin.Context) {
	type SpecialParam struct {
		SpecialID uint `json:"specialId"`
	}
	param := SpecialParam{}
	if err := c.ShouldBindWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	global.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Special{}, param.SpecialID).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return err
		}
		if err := tx.Delete(&models.SpecialDay{}, "special_id=?", param.SpecialID).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return err
		}
		return nil
	})
	apiReturn.Success(c)
}
