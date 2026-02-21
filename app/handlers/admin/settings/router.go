package settings

import (
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox/xwa/xmwas"
)

func Router(rg *xin.RouterGroup) {
	rg.GET("/", SettingIndex)
	rg.POST("/save", SettingSave)
	rg.POST("/export", SettingExport)
	rg.POST("/import", SettingImport)

	addFilesHandlers(rg.Group("/files"))
}

func addFilesHandlers(rg *xin.RouterGroup) {
	rg.POST("/upload", SettingFileUpload)
	rg.GET("/preview/*fid", SettingFilePreview)

	xin.StaticFSFunc(rg, "/dnload/", middles.TenantHFS, xin.DisableAcceptRanges, xmwas.XCC.Handle)
}
