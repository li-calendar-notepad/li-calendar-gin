package apiReturn

import (
	"calendar-note-gin/lib/global"

	"github.com/gin-gonic/gin"
)

func ApiReturn(ctx *gin.Context, code int, msg string, data interface{}) {
	returnData := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	if data != nil {
		returnData["data"] = data
	}
	ctx.JSON(200, returnData)
}

// 返回成功
func SuccessData(ctx *gin.Context, data interface{}) {
	ApiReturn(ctx, 0, "OK", data)
}

// 返回列表
func SuccessListData(ctx *gin.Context, list interface{}, total int64) {

	ApiReturn(ctx, 0, "OK", gin.H{
		"list":  list,
		"total": total,
	})
}

// 返回成功，没有data数据
func Success(ctx *gin.Context) {
	ApiReturn(ctx, 0, "OK", nil)
}

// 返回列表数据
func ListData(ctx *gin.Context, list interface{}, total int64) {
	data := map[string]interface{}{
		"list":  list,
		"total": total,
	}
	ApiReturn(ctx, 0, "OK", data)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorCode(ctx *gin.Context, code int, errMsg string, data interface{}) {
	ApiReturn(ctx, code, errMsg, data)
}

// 返回错误 普通提示错误
func Error(ctx *gin.Context, errMsg string) {
	ErrorCode(ctx, -1, errMsg, nil)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorNoAccess(ctx *gin.Context) {
	ErrorCode(ctx, 1005, global.Lang.Get("common.no_access"), nil)
}

// 返回错误 参数错误
func ErrorParamFomat(ctx *gin.Context, errMsg string) {
	Error(ctx, global.Lang.GetAndInsert("common.api_error_param_format", "[", errMsg, "]"))
}

// 返回错误 数据库
func ErrorDatabase(ctx *gin.Context, errMsg string) {
	Error(ctx, global.Lang.GetAndInsert("common.db_error", "[", errMsg, "]"))
}
