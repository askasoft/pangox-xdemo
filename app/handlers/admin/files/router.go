package files

import (
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app/handlers/files"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox/xwa/xmwas"
)

func Router(rg *xin.RouterGroup) {
	rg.GET("/", FileIndex)
	rg.POST("/list", FileList)
	rg.POST("/updates", FileUpdates)
	rg.POST("/deletes", FileDeletes)
	rg.POST("/deleteb", FileDeleteBatch)

	rg.GET("/preview/*fid", files.Preview)

	xin.StaticFSFunc(rg, "/dnload/", middles.TenantHFS, xin.DisableAcceptRanges, xmwas.XCC.Handle)
}
