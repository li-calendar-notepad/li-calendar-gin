package system

import (
	v1 "calendar-note-gin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitLoginRouter(router *gin.RouterGroup) {
	loginApi := v1.ApiGroupApp.LoginApi

	router.POST("/login", loginApi.Login)
	router.POST("/logout", loginApi.Logout)
	router.POST("/login/register", loginApi.Register)
	router.POST("/login/register/commit", loginApi.Commit)
	router.POST("/login/sendResetPasswordVCode", loginApi.SendResetPasswordVCode)
	router.POST("/login/resetPasswordByVCode", loginApi.ResetPasswordByVCode)

}
