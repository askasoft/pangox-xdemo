package server

import (
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pangox-xdemo/app"
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
