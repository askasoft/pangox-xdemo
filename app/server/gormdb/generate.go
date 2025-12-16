package gormdb

import (
	"path/filepath"

	"github.com/askasoft/gogormx/gormx"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/str"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
)

// Generate DDL sql
func GenerateDDL(outdir string) error {
	driver := ini.GetString("database", "driver")
	if outdir == "" {
		outdir = "./data/sqls/"
	}

	outfile := filepath.Join(outdir, str.If(driver == "pgx", "postgres", driver)+".sql")

	log.Infof("Generate schema DDL: '%s'", outfile)

	gsp := &gormx.GormSQLPrinter{}

	dbc := &gorm.Config{
		DryRun:         true,
		NamingStrategy: gormschema.NamingStrategy{TablePrefix: "build."},
		Logger:         gsp,
	}

	gdd := dialector(driver)
	dms := dbmodels(driver)

	gdb, err := gorm.Open(gdd, dbc)
	if err != nil {
		return err
	}

	gmi := gdb.Migrator()
	for _, m := range dms {
		gsp.Printf("---------------------------------")
		if err := gmi.CreateTable(m); err != nil {
			return err
		}
	}

	qte := sqx.GetQuoter(driver)

	sql := gsp.SQL()
	sql = str.ReplaceAll(sql, "idx_build_", "idx_")
	sql = str.ReplaceAll(sql, qte.Quote("build"), qte.Quote("SCHEMA"))

	return fsu.WriteString(outfile, sql, 0660)
}
