package admin

import (
	"calendar-note-gin/api/v1/common/apiReturn"
	"calendar-note-gin/lib/cmn"

	"github.com/gin-gonic/gin"
)

type AboutApi struct {
}

func (a *AboutApi) Get(c *gin.Context) {
	version := cmn.GetSysVersionInfo()

	apiReturn.SuccessData(c, gin.H{
		"version": version,
	})
}
