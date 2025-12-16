package data

import (
	"embed"
)

// FS embed data
//
//go:embed sqls *.csv
var FS embed.FS
