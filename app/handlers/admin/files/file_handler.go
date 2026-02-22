package files

import (
	"net/http"

	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/args"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox-xdemo/app/utils/tbsutil"
)

var fileListCols = []string{
	"id",
	"name",
	"ext",
	"tag",
	"size",
	"time",
}

func bindFileQueryArg(c *xin.Context) (fqa *args.FileQueryArg, err error) {
	fqa = &args.FileQueryArg{}
	fqa.Order = "-time"

	err = c.Bind(fqa)
	fqa.Orders.Normalize(fileListCols...)
	return
}

func bindFileMaps(c *xin.Context, h xin.H) {
}

func FileIndex(c *xin.Context) {
	h := middles.H(c)

	fqa, _ := bindFileQueryArg(c)

	h["Q"] = fqa
	bindFileMaps(c, h)

	c.HTML(http.StatusOK, "admin/files/files", h)
}

func FileList(c *xin.Context) {
	fqa, err := bindFileQueryArg(c)
	if err != nil {
		args.AddBindErrors(c, err, "file.")
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	tt := tenant.FromCtx(c)

	fqa.Total, err = tt.CountFiles(app.SDB(), fqa)
	if err != nil {
		c.AddError(err)
		c.JSON(http.StatusInternalServerError, middles.E(c))
		return
	}

	h := middles.H(c)

	fqa.Pager.Normalize(tbsutil.GetPagerLimits(c.Locale)...)

	if fqa.Total > 0 {
		results, err := tt.FindFiles(app.SDB(), fqa, fileListCols...)
		if err != nil {
			c.AddError(err)
			c.JSON(http.StatusInternalServerError, middles.E(c))
			return
		}

		h["Files"] = results
		fqa.Count = len(results)
	}

	h["Q"] = fqa
	bindFileMaps(c, h)

	c.HTML(http.StatusOK, "admin/files/files_list", h)
}
