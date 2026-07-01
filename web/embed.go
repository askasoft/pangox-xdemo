package web

import (
	"embed"
	"io/fs"
	"os"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-assets/html/bootstrap5"
	"github.com/askasoft/pangox-assets/html/bootswatch5/cosmo"
	"github.com/askasoft/pangox-assets/html/bootswatch5/flatly"
	"github.com/askasoft/pangox-assets/html/bootswatch5/pulse"
	"github.com/askasoft/pangox-assets/html/corejs"
	"github.com/askasoft/pangox-assets/html/docxjs"
	"github.com/askasoft/pangox-assets/html/flagicons"
	"github.com/askasoft/pangox-assets/html/fontawesome7"
	"github.com/askasoft/pangox-assets/html/jquery3"
	"github.com/askasoft/pangox-assets/html/jszip"
	"github.com/askasoft/pangox-assets/html/pdfjs"
	"github.com/askasoft/pangox-assets/html/pdfviewer"
	"github.com/askasoft/pangox-assets/html/plugins"
	"github.com/askasoft/pangox-assets/html/summernote"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox/xwa/xmwas"
)

// Static embed static folder
var Statics = map[string]fs.FS{
	"bootstrap5":         bootstrap5.FS,
	"bootswatch5/cosmo":  cosmo.FS,
	"bootswatch5/flatly": flatly.FS,
	"bootswatch5/pulse":  pulse.FS,
	"corejs":             corejs.FS,
	"flagicons":          flagicons.FS,
	"fontawesome7":       fontawesome7.FS,
	"jquery3":            jquery3.FS,
	"jszip":              jszip.FS,
	"docxjs":             docxjs.FS,
	"pdfjs":              pdfjs.FS,
	"pdfviewer":          pdfviewer.FS,
	"plugins":            plugins.FS,
	"summernote":         summernote.FS,
}

//go:embed assets favicon.ico robots.txt
var FS embed.FS

var (
	WAS = fsu.ProxyFS{FS: FS}
	LDS = fsu.ProxyFS{}
)

func AddWebAssetsHandlers(rg *xin.RouterGroup) {
	xmwas.AddStaticsHandlers(rg.Group("/static/"+app.Revision()), Statics)
	xmwas.AddStaticSubFolderHandlers(rg.Group(""), WAS, app.Revision())
	xmwas.AddStaticSubFileHandlers(rg.Group(""), WAS)
	xmwas.AddStaticSubFolderHandlers(rg.Group(""), LDS)
}

func ConfigWebAssets() {
	was := ini.GetString("app", "webassets")
	if was == "" {
		WAS.FS = FS
		log.Info("Using embedded web assets")
	} else {
		WAS.FS = os.DirFS(was)
		log.Infof("Using external web assets: '%s'", was)
	}

	lds := ini.GetString("app", "localdocs")
	if lds == "" {
		LDS.FS = nil
		log.Info("Using empty local documents")
	} else {
		LDS.FS = os.DirFS(lds)
		log.Infof("Using local documents: '%s'", lds)
	}
}
