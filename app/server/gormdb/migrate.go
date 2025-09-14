package gormdb

import (
	"time"

	"github.com/askasoft/gogormx/log/sqlog/gormlog"
	"github.com/askasoft/gogormx/xsm/mysm/mygormsm"
	"github.com/askasoft/gogormx/xsm/pgsm/pggormsm"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
)

// Migrate Schemas
func MigrateSchemas(schemas ...string) (err error) {
	if len(schemas) == 0 {
		if ini.GetBool("tenant", "multiple") {
			schemas, err = listSchemas()
			if err != nil {
				return
			}
		} else {
			schemas = []string{ini.GetString("database", "schema", "public")}
		}
	}

	for _, schema := range schemas {
		if err = migrateSchema(schema); err != nil {
			return
		}
	}
	return
}

func migrateSchema(schema string) error {
	log.Infof("Migrate schema %q", schema)

	dbc := &gorm.Config{
		NamingStrategy: gormschema.NamingStrategy{TablePrefix: schema + "."},
		Logger: gormlog.NewGormLogger(
			log.GetLogger("SQL"),
			time.Second,
		),
	}

	dbt := ini.GetString("database", "driver")
	gdd := dialector(dbt)

	gdb, err := gorm.Open(gdd, dbc)
	if err != nil {
		return err
	}

	err = gdb.AutoMigrate(dbmodels(dbt)...)

	if db, err := gdb.DB(); err == nil {
		db.Close()
	}
	return err
}

func listSchemas() ([]string, error) {
	gdb, err := openDatabase()
	if err != nil {
		return nil, err
	}

	dbt := ini.GetString("database", "driver")
	switch dbt {
	case "mysql":
		return mygormsm.SM(gdb).ListSchemas()
	default:
		return pggormsm.SM(gdb).ListSchemas()
	}
}

func openDatabase() (*gorm.DB, error) {
	dbt := ini.GetString("database", "driver")

	gdd := dialector(dbt)

	gdc := &gorm.Config{
		Logger: gormlog.NewGormLogger(
			log.GetLogger("SQL"),
			time.Second,
		),
		SkipDefaultTransaction: true,
	}

	return gorm.Open(gdd, gdc)
}
