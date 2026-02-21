package files

import (
	"encoding/csv"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox-xdemo/app/utils/docutil"
)

func Preview(c *xin.Context) {
	fid := c.Param("fid")
	dnloadURL := app.Base() + "/files/dnload" + fid

	PreviewFile(c, fid, dnloadURL)
}

func PreviewFile(c *xin.Context, fid, dnloadURL string) {
	if fid == "" {
		middles.NotFound(c)
		return
	}

	tt := tenant.FromCtx(c)

	file, err := tt.FS().FindFile(fid)
	if errors.Is(err, sqlx.ErrNoRows) {
		middles.NotFound(c)
		return
	}
	if err != nil {
		c.AddError(err)
		middles.InternalServerError(c)
		return
	}

	h := middles.H(c)
	h["File"] = file
	h["DnloadURL"] = dnloadURL

	ext := str.ToLower(filepath.Ext(file.Name))
	if ext == ".htm" {
		ext = ".html"
	}

	switch ext {
	case ".docx", ".html", ".pdf", ".pptx", ".xls":
		c.HTML(http.StatusOK, "files/preview"+ext, h)
		return
	case ".xlsx":
		xsm, err := readXlsx(tt, fid)
		if err != nil {
			c.AddError(err)
			middles.InternalServerError(c)
			return
		}
		h["Sheets"] = xsm
		c.HTML(http.StatusOK, "files/preview"+ext, h)
		return
	case ".txt", ".tsv", ".csv":
		data, err := tt.FS().ReadFile(fid)
		if err != nil {
			c.AddError(err)
			middles.InternalServerError(c)
			return
		}

		text := docutil.ReadTextFromTextData(data)

		switch ext {
		case ".tsv":
			cr := csv.NewReader(str.NewReader(text))
			cr.Comma = '\t'
			rows, err := cr.ReadAll()
			if err == nil {
				h["Rows"] = rows
			} else {
				h["Text"] = text
			}
			ext = ".csv"
		case ".csv":
			cr := csv.NewReader(str.NewReader(text))
			rows, err := cr.ReadAll()
			if err == nil {
				h["Rows"] = rows
			} else {
				h["Text"] = text
			}
		default:
			h["Text"] = text
		}

		c.HTML(http.StatusOK, "files/preview"+ext, h)
		return
	default:
		c.Redirect(http.StatusFound, dnloadURL)
		return
	}
}

func readXlsx(tt *tenant.Tenant, fid string) (*linkedhashmap.LinkedHashMap[string, [][]string], error) {
	data, err := tt.FS().ReadFile(fid)
	if err != nil {
		return nil, err
	}

	return docutil.ReadXlsxDataToMap(data)
}
