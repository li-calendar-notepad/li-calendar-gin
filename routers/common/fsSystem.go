package common

// ================
// Fs文件系统
// ================

import (
	"calendar-note-gin/assets"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const assets_path = "assets/"

// 文件监听 （所有参数后不能以"/"结尾）
func FSListenFile(r *gin.Engine, path, fs_file_name string) {
	r.GET(path, func(c *gin.Context) {
		FsOutputFileByAssets(c, fs_file_name)
	})
}

// 文件夹监听  （所有参数后不能以"/"结尾）
func FSListenDir(r *gin.Engine, path, fs_dir string) {
	r.GET(path+"/*action", func(c *gin.Context) {
		action := c.Param("action")
		file_name := fs_dir + action
		FsOutputFileByAssets(c, file_name)
	})
}

// 文件夹监听扩展名  （所有参数后不能以"/"结尾）
func FSListenDirExt(r *gin.Engine, path, fs_dir, ext string) {
	fs_dir = assets_path + fs_dir
	r.GET(path+"/*action", func(c *gin.Context) {
		action := c.Param("action")

		file_name := fs_dir + action + ext
		FsOutputFileByAssets(c, file_name)
	})
}

// 获取类型
// ：文件名称，文件字节内容
func getContentType(file_name string, file_content []byte) string {
	ctype := mime.TypeByExtension(filepath.Ext(file_name))
	if ctype == "" {
		ctype = http.DetectContentType(file_content)
	}
	return ctype
}

// 文件系统输出文件包内文件
func FsOutputFileByAssets(c *gin.Context, file_name string) {
	file_name = assets_path + file_name
	indexHtml, err := assets.Asset(file_name)
	if err != nil {
		c.Writer.WriteHeader(404)
		return
	}
	c.Writer.WriteHeader(200)
	c.Writer.Header().Add("Content-Type", getContentType(file_name, indexHtml))
	_, _ = c.Writer.Write(indexHtml)
	c.Writer.Flush()
}
