package files

import (
	"encoding/csv"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/xin"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/middles"
	"github.com/askasoft/pangox-xdemo/app/tenant"
	"github.com/askasoft/pangox-xdemo/app/utils/docutil"
)

const (
	PREVIEW_MAX_RUNES = 50000
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
			if !errors.Is(err, docutil.ErrOverflow) {
				c.AddError(err)
				middles.InternalServerError(c)
				return
			}
			h["Warning"] = tbs.GetText(c.Locale, "file.warning.overflow")
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
		case ".csv", ".tsv":
			cr := csv.NewReader(str.NewReader(text))
			cr.Comma = gog.If(ext == ".tsv", '\t', ',')

			data, err := cr.ReadAll()
			if err == nil {
				rows, total := docutil.LimitRows(data, 0, PREVIEW_MAX_RUNES)
				if total > PREVIEW_MAX_RUNES {
					h["Warning"] = tbs.GetText(c.Locale, "file.warning.overflow")
				}
				h["Rows"] = rows
			} else {
				h["Text"] = str.Ellipsis(text, PREVIEW_MAX_RUNES)
			}
			ext = ".csv"
		default:
			h["Text"] = str.Ellipsis(text, PREVIEW_MAX_RUNES)
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

	return docutil.ReadXlsxDataToMap(data, 100)
}
