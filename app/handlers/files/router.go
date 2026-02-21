package files

import (
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox/xwa/xmwas"
)

func Router(rg *xin.RouterGroup) {
	rg.Use(xmwas.XAC.Handle) // access control
	rg.OPTIONS("/*path", xin.Next)

	rg.POST("/upload", Upload)
	rg.POST("/uploads", Uploads)

	rg.GET("/preview/*fid", Preview)

	xin.StaticFSFunc(rg, "/dnload/", middles.TenantHFS, xin.DisableAcceptRanges, xmwas.XCC.Handle)
}
