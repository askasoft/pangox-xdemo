package files

import (
	"net/http"

	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/args"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/tenant"
)

func FileUpdates(c *xin.Context) {
	fua := &args.FileUpdatesArg{}
	if err := fua.Bind(c); err != nil {
		args.AddBindErrors(c, err, "file.")
	}
	if len(c.Errors) > 0 {
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	tt := tenant.FromCtx(c)

	var cnt int64
	err := app.SDB().Transaction(func(tx *sqlx.Tx) (err error) {
		cnt, err = tt.UpdateFiles(tx, fua)
		if err != nil {
			return
		}
		if cnt > 0 {
			err = tt.AddAuditLog(tx, c, models.AL_FILES_UPDATES, cnt, fua.String())
		}
		return
	})
	if err != nil {
		c.AddError(err)
		c.JSON(http.StatusInternalServerError, middles.E(c))
		return
	}

	c.JSON(http.StatusOK, xin.H{
		"success": tbs.Format(c.Locale, "pet.success.updates", cnt),
		"updates": fua,
	})
}

func FileDeletes(c *xin.Context) {
	pka := &args.PKArg{}
	if err := pka.Bind(c); err != nil {
		c.AddError(args.InvalidIDError(c))
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	tt := tenant.FromCtx(c)

	var cnt int64
	err := app.SDB().Transaction(func(tx *sqlx.Tx) (err error) {
		cnt, err = tt.DeleteFiles(tx, pka.PKs()...)
		if err != nil {
			return
		}

		if cnt > 0 {
			err = tt.AddAuditLog(tx, c, models.AL_FILES_DELETES, cnt, pka.String())
		}
		return
	})
	if err != nil {
		c.AddError(err)
		c.JSON(http.StatusInternalServerError, middles.E(c))
		return
	}

	c.JSON(http.StatusOK, xin.H{
		"success": tbs.Format(c.Locale, "file.success.deletes", cnt),
	})
}

func FileDeleteBatch(c *xin.Context) {
	pqa, err := bindFileQueryArg(c)
	if err != nil {
		args.AddBindErrors(c, err, "file.")
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	if !pqa.HasFilters() {
		c.AddError(tbs.Error(c.Locale, "error.param.nofilter"))
		c.JSON(http.StatusBadRequest, middles.E(c))
		return
	}

	tt := tenant.FromCtx(c)

	var cnt int64
	err = app.SDB().Transaction(func(tx *sqlx.Tx) (err error) {
		cnt, err = tt.DeleteFilesQuery(tx, pqa)
		if err != nil {
			return
		}

		if cnt > 0 {
			err = tt.AddAuditLog(tx, c, models.AL_FILES_DELETES, cnt, pqa.String())
		}
		return
	})
	if err != nil {
		c.AddError(err)
		c.JSON(http.StatusInternalServerError, middles.E(c))
		return
	}

	c.JSON(http.StatusOK, xin.H{
		"success": tbs.Format(c.Locale, "file.success.deletes", cnt),
	})
}
