package gormdb

import (
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/server/gormdb/mymodels"
	"github.com/askasoft/pangox/xfs"
	"github.com/askasoft/pangox/xjm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	log.Infof("Connect Database (%s): %s", dbt, dsn)

	switch dbt {
	case "mysql":
		return mysql.Open(dsn)
	default:
		return postgres.Open(dsn)
	}
}
