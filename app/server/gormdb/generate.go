package gormdb

import (
	"bufio"
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

	// format sql
	var sb str.Builder
	sc := bufio.NewScanner(str.NewReader(sql))
	for sc.Scan() {
		line := sc.Text()

		if !str.StartsWith(line, "CREATE TABLE") {
			sb.WriteString(line)
			sb.WriteByte('\n')
			continue
		}

		i := str.IndexByte(line, '(')
		if i >= 0 {
			sb.WriteString(line[:i+1])

			var keys []string
			line = line[i+1:]
			if i = str.LastIndexByte(line, ')'); i >= 0 {
				line = line[:i]
				if i = str.Index(line, "PRIMARY KEY"); i >= 0 {
					keys = str.FieldsByte(line[i:], ',')
					line = line[:i]
				}
			}

			ss := str.Strips(str.FieldsByte(line, ','))
			for i, s := range ss {
				if '0' <= s[0] && s[0] <= '9' {
					sb.WriteByte(',')
					sb.WriteString(s)
					continue
				}

				if i > 0 {
					sb.WriteByte(',')
				}
				sb.WriteString("\n\t")
				sb.WriteString(s)
			}

			for _, k := range keys {
				if k[0] == '"' || k[0] == '`' {
					sb.WriteByte(',')
					sb.WriteString(k)
					continue
				}
				sb.WriteString(",\n\t")
				sb.WriteString(k)
			}
			sb.WriteString("\n);\n")
		}
	}

	return fsu.WriteString(outfile, sb.String(), 0660)
}
