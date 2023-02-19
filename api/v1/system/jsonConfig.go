package system

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/global"
	"calendar-note-gin/models"

	"github.com/gin-gonic/gin"
)

type JsonConfig struct {
}

// 样式文件导入
func (l JsonConfig) EventStyleImport(c *gin.Context) {
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

// // 样式文件导出
// func (l JsonConfig) EventStyleExport(c *gin.Context) {
// 	model := jsonConfig.ConfigModel{
// 		Data: jsonConfig.ElementStyleModel{
// 			Title:     "testExport",
// 			ClassName: "test",
// 		},
// 	}
// 	content, _ := jsonConfig.BuildExportFile(&model)
// 	jsonConfig.Write(c, "EventStyle.ey.", content)
// }
