package docutil

import (
	"bytes"

	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/doc/htmlx"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/wcu"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"
)

const (
	CharsetDetectLength = 4096
)

func ParseHTMLFile(name string, charsets ...string) (*html.Node, error) {
	return htmlx.ParseHTMLFile(name, CharsetDetectLength, charsets...)
}

func DetectAndReadFile(filename string, charsets ...string) ([]byte, string, error) {
	return wcu.DetectAndReadFile(filename, CharsetDetectLength, charsets...)
}

func ReadTextFromTextData(data []byte) string {
	r, _, err := wcu.DetectAndTransform(bytes.NewReader(data), CharsetDetectLength, false)
	if err == nil {
		bs, err := iox.ReadAll(r)
		if err == nil {
			return str.UnsafeString(bs)
		}
	}
	return str.UnsafeString(data)
}

func ReadTextFromTextFile(filename string, charset string) (string, error) {
	bs, _, err := DetectAndReadFile(filename, charset)
	if err != nil {
		return "", err
	}
	return str.UnsafeString(bs), err
}

func ReadXlsxDataToMap(data []byte) (*linkedhashmap.LinkedHashMap[string, [][]string], error) {
	fx, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// xop := excelize.Options{RawCellValue: true}
	xsm := linkedhashmap.NewLinkedHashMap[string, [][]string]()
	for i := range fx.SheetCount {
		sn := fx.GetSheetName(i)

		rows, err := fx.GetRows(sn)
		if err != nil {
			return nil, err
		}

		// try to convert number to time for unset cell type
		// for y, cols := range rows {
		// 	for x, cv := range cols {
		// 		if str.IsDecimal(cv) {
		// 			cn, err := excelize.CoordinatesToCellName(x+1, y+1)
		// 			if err != nil {
		// 				return nil, err
		// 			}

		// 			ct, err := fx.GetCellType(sn, cn)
		// 			if err != nil {
		// 				return nil, err
		// 			}

		// 			if ct == excelize.CellTypeUnset {
		// 				t, err := excelize.ExcelDateToTime(num.Atof(cv), false)
		// 				if err == nil {
		// 					if str.IsNumber(cv) {
		// 						cols[x] = app.FormatDate(t)
		// 					} else {
		// 						cols[x] = app.FormatTime(t)
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}
		// }
		xsm.Set(sn, rows)
	}

	return xsm, nil
}
