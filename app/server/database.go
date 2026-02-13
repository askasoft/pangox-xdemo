package server

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/doc/csvx"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/schema"
	"github.com/askasoft/pangox-xdemo/app/utils/sqlutil"
	"github.com/askasoft/pangox/xwa/xsqls"
)

func init() {
	xsqls.RegisterGetErrLogLevel("mysql", sqlutil.GetMysqlErrLogLevel)
	xsqls.RegisterGetErrLogLevel("pgx", sqlutil.GetPgxErrLogLevel)
}

func initDatabase() {
	if err := xsqls.OpenDatabase(); err != nil {
		log.Fatal(app.ExitErrDB, err)
	}
}

func reloadDatabase() {
	if err := xsqls.OpenDatabase(); err != nil {
		log.Error(err)
	}
}

func dbIterateSchemas(fn func(sm schema.Schema) error, schemas ...string) error {
	if len(schemas) == 0 {
		return schema.Iterate(fn)
	}

	for _, s := range schemas {
		if err := fn(schema.Schema(s)); err != nil {
			return err
		}
	}
	return nil
}

func dbMigrateSettings(schemas ...string) error {
	settings, err := schema.ReadSettingsFile()
	if err != nil {
		return err
	}

	return dbIterateSchemas(func(sm schema.Schema) error {
		return sm.MigrateSettings(settings)
	}, schemas...)
}

func dbExportSettings(outdir string, schemas ...string) error {
	if err := os.MkdirAll(outdir, 0770); err != nil {
		return err
	}

	return dbIterateSchemas(func(sm schema.Schema) error {
		outfile := filepath.Join(outdir, string(sm)+".csv")

		log.Infof("Export settings %q to '%s'", sm, outfile)

		fw, err := os.Create(outfile)
		if err != nil {
			return err
		}
		defer fw.Close()

		settings, err := sm.SelectSettings(app.SDB())
		if err != nil {
			return err
		}

		return exportSettings(fw, settings)
	}, schemas...)
}

func exportSettings(w io.Writer, settings []*models.Setting) error {
	cw := csv.NewWriter(w)
	cw.UseCRLF = true
	defer cw.Flush()

	if err := cw.Write([]string{"Name", "Value"}); err != nil {
		return err
	}

	for _, s := range settings {
		if err := cw.Write([]string{s.Name, s.Value}); err != nil {
			return err
		}
	}

	return nil
}

func dbImportSettings(srcdir string) error {
	return dbIterateSchemas(func(sm schema.Schema) error {
		srcfile := filepath.Join(srcdir, string(sm)+".csv")

		if err := fsu.FileExists(srcfile); err != nil {
			log.Infof("Skip import settings %q: %v", sm, err)
			return nil
		}

		log.Infof("Import settings %q from '%s'", sm, srcfile)

		var settings []*models.Setting
		if err := csvx.ScanFile(srcfile, &settings); err != nil {
			return err
		}

		err := app.SDB().Transaction(func(tx *sqlx.Tx) error {
			return sm.SaveSettings(tx, settings, asg.First(app.Locales()))
		})
		return err
	})
}

func dbMigrateSupers(schemas ...string) error {
	return dbIterateSchemas(func(sm schema.Schema) error {
		return sm.MigrateSuper()
	}, schemas...)
}

func dbExecSQL(sqlfile string, schemas ...string) error {
	log.Infof("Read SQL file '%s'", sqlfile)

	sql, err := fsu.ReadString(sqlfile)
	if err != nil {
		return err
	}

	return dbIterateSchemas(func(sm schema.Schema) error {
		return sm.ExecSQL(sql)
	}, schemas...)
}

func dbSchemaInit(schemas ...string) error {
	return dbIterateSchemas(func(sm schema.Schema) error {
		return sm.InitSchema()
	}, schemas...)
}

func dbSchemaCheck(schemas ...string) error {
	return dbIterateSchemas(func(sm schema.Schema) error {
		sm.CheckSchema(app.SDB())
		return nil
	}, schemas...)
}

func dbSchemaUpdate(schemas ...string) error {
	return dbIterateSchemas(func(sm schema.Schema) error {
		return sm.UpdateSchema(app.SDB())
	}, schemas...)
}

func dbSchemaVacuum(schemas ...string) error {
	return dbIterateSchemas(func(sm schema.Schema) error {
		return sm.VacuumSchema(app.SDB())
	}, schemas...)
}
