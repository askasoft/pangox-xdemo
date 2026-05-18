package docutil

import (
	"bytes"
	"errors"
	"unicode/utf8"

	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/doc/htmlx"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/wcu"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"
)

const (
	CharsetDetectLength = 4096
)

func ParseHTMLFile(name string, charsets ...string) (*html.Node, error) {
	return htmlx.ParseHTMLFile(name, CharsetDetectLength, charsets...)
}

func ParseHTMLData(data []byte, charsets ...string) (*html.Node, error) {
	r, _, err := wcu.DetectAndTransform(bytes.NewReader(data), CharsetDetectLength, true)
	if err != nil {
		return nil, err
	}
	return html.Parse(r)
}

func DetectAndReadFile(filename string, charsets ...string) ([]byte, string, error) {
	return wcu.DetectAndReadFile(filename, CharsetDetectLength, charsets...)
}

func ReadTextFromTextData(data []byte) string {
	r, _, err := wcu.DetectAndTransform(bytes.NewReader(data), CharsetDetectLength, false)
	if err == nil {
		r, err = iox.SkipBOM(r)
		if err == nil {
			bs, err := iox.ReadAll(r)
			if err == nil {
				return str.UnsafeString(bs)
			}
		}
	}
	return str.UnsafeString(data)
}

func ReadTextFromTextFile(filename string, charset string) (string, error) {
	bs, _, err := DetectAndReadFile(filename, charset)
	if err != nil {
		return "", err
	}
	if len(bs) > 0 {
		r, z := utf8.DecodeRune(bs)
		if r == iox.BOM {
			bs = bs[z:]
		}
	}
	return str.UnsafeString(bs), err
}

func LimitRows(rows [][]string, total, limit int) ([][]string, int) {
	if limit <= 0 {
		return rows, 0
	}

	for y, r := range rows {
		if total > limit {
			rows = rows[:y]
			break
		}

		for x, c := range r {
			if total > limit {
				r[x] = ""
				continue
			}

			count := str.RuneCount(c)
			total += count
			if total > limit {
				r[x] = str.Ellipsis(c, count-(total-limit))
			}
		}
	}

	return rows, total
}

var ErrOverflow = errors.New("overflow")

func ReadXlsxDataToMap(data []byte, limit int) (*linkedhashmap.LinkedHashMap[string, [][]string], error) {
	fx, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	total := 0

	xsm := linkedhashmap.NewLinkedHashMap[string, [][]string]()
	for i := range fx.SheetCount {
		sn := fx.GetSheetName(i)
		if total >= limit {
			xsm.Set(sn, nil)
			continue
		}

		rows, err := fx.GetRows(sn)
		if err != nil {
			return nil, err
		}

		rows, total = LimitRows(rows, total, limit)
		xsm.Set(sn, rows)
	}

	return xsm, gog.If(total > limit, ErrOverflow, nil)
}

// TODO: fix number time
func FixXlsxNumberTime(fx *excelize.File, sn string, rows [][]string) error {
	// xop := excelize.Options{RawCellValue: true}

	// try to convert number to time for unset cell type
	for y, cols := range rows {
		for x, cv := range cols {
			if str.IsDecimal(cv) {
				cn, err := excelize.CoordinatesToCellName(x+1, y+1)
				if err != nil {
					return err
				}

				ct, err := fx.GetCellType(sn, cn)
				if err != nil {
					return err
				}

				if ct == excelize.CellTypeUnset {
					t, err := excelize.ExcelDateToTime(num.Atof(cv), false)
					if err == nil {
						if str.IsNumber(cv) {
							cols[x] = app.FormatDate(t)
						} else {
							cols[x] = app.FormatTime(t)
						}
					}
				}
			}
		}
	}

	return nil
}
