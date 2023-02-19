package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// var ApiReturn = apiReturn.ApiReturn

type SubjectApi struct {
}

func (a *SubjectApi) DeleteByIdAndItemId(c *gin.Context) {
	param := ParamEventCreateData{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.Get("common.api_error_param_format")+err.Error())
		c.Abort()
		return
	}

	mEvent := models.Event{}
	if err := mEvent.DeleteByIdAndItemId(param.EventId, param.ItemId); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.db_error", "[", err.Error(), "]"))
	} else {
		apiReturn.Success(c)
	}

}

// 创建一条
func (a *SubjectApi) Create(c *gin.Context) {
	param := models.Subject{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}
	allowField := []string{"Title", "ClassName", "ItemId", "content"}
	find := models.Subject{}
	if rows := global.Db.First(&find, "title=? AND item_id=?", param.Title, param.ItemId).RowsAffected; rows != 0 {
		apiReturn.Error(c, global.Lang.GetAndInsert("subject.err_title_already_exists"))
		return
	}

	if err := global.Db.Select(allowField).Create(&param).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{
		"subjectId": param.ID,
	})

}

func (a *SubjectApi) Update(c *gin.Context) {
	param := models.Subject{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}
	allowField := []string{"Title", "ClassName", "ItemId", "SubjectId", "content"}

	find := models.Subject{}
	if rows := global.Db.First(&find, "title=? AND item_id=?", param.Title, param.ItemId).RowsAffected; rows != 0 && param.SubjectId != find.ID {
		apiReturn.Error(c, global.Lang.GetAndInsert("subject.err_title_already_exists"))
		return
	}

	if err := global.Db.Select(allowField).Where("id=?", param.SubjectId).Updates(&param).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)

}

func (a *SubjectApi) Deletes(c *gin.Context) {
	type SubjectIds struct {
		SubjectIds []uint
	}
	param := SubjectIds{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}
	// global.Db.Debug().Delete(&models.Subject{}, &param.SubjectIds)
	if err := global.Db.Delete(&models.Subject{}, &param.SubjectIds).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

func (a *SubjectApi) GetList(c *gin.Context) {
	param := models.Subject{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}
	subjectList := []models.Subject{}
	var count int64
	err := global.Db.Order("id Desc").Select("ClassName", "Content", "Title", "ItemId", "ID").Find(&subjectList, "item_id=?", param.ItemId).Count(&count).Error
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	subjectListMap := []map[string]interface{}{}
	for _, v := range subjectList {
		subjectListMap = append(subjectListMap, map[string]interface{}{
			"subjectId": v.ID,
			"className": v.ClassName,
			"content":   v.Content,
			"title":     v.Title,
			"itemId":    v.ItemId,
		})
	}
	apiReturn.SuccessListData(c, subjectListMap, count)
}
