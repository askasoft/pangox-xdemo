package settings

import (
	"github.com/askasoft/pango/xin"
)

func Router(rg *xin.RouterGroup) {
	rg.GET("/", SettingIndex)
	rg.POST("/save", SettingSave)
	rg.POST("/export", SettingExport)
	rg.POST("/import", SettingImport)
}
