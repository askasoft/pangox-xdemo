package gormdb

import (
	"time"

	"github.com/askasoft/gogormx/log/sqlog/gormlog"
	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/server/gormdb/mymodels"
	"github.com/askasoft/pangox/xfs"
	"github.com/askasoft/pangox/xjm"
	"github.com/askasoft/pangox/xwa/xsqls"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
)

var pgModels = []any{
	&xfs.File{},
	&xjm.Job{},
	&xjm.JobLog{},
	&xjm.JobChain{},
	&models.User{},
	&models.Config{},
	&models.AuditLog{},
	&models.Pet{},
	&xsqls.SchemaChange{},
}

var myModels = []any{
	&xfs.File{},
	&xjm.Job{},
	&xjm.JobLog{},
	&xjm.JobChain{},
	&models.User{},
	&models.Config{},
	&mymodels.AuditLog{},
	&mymodels.Pet{},
	&xsqls.SchemaChange{},
}

func dbmodels(dbt string) []any {
	switch dbt {
	case "mysql":
		return myModels
	default:
		return pgModels
	}
}

func dialector(dbt string) gorm.Dialector {
	dsn := ini.GetString("database", "source")

	log.Debugf("Connect Database (%s): %s", dbt, dsn)

	switch dbt {
	case "mysql":
		return mysql.Open(dsn)
	default:
		return postgres.Open(dsn)
	}
}

func OpenDatabase(schema ...string) (*gorm.DB, error) {
	dbt := ini.GetString("database", "driver")

	gdd := dialector(dbt)

	gdc := &gorm.Config{
		Logger: gormlog.NewGormLogger(
			log.GetLogger("SQL"),
			ini.GetDuration("database", "slowSQL", 2*time.Second),
		),
		SkipDefaultTransaction: true,
	}
	if sm := asg.First(schema); sm != "" {
		gdc.NamingStrategy = gormschema.NamingStrategy{TablePrefix: sm + "."}
	}

	return gorm.Open(gdd, gdc)
}
