package schema

import (
	"errors"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/data"
	"github.com/askasoft/pangox/xwa/xsqls"
)

func (sm Schema) ExecSQL(sqls string) error {
	logger := log.GetLogger("SQL")
	logger.Info(str.PadCenter(" "+string(sm)+" ", 60, "="))

	sqls = str.ReplaceAll(sqls, "SCHEMA", string(sm))

	return xsqls.ExecSQL(app.SDB(), string(sm), sqls, logger)
}

func (sm Schema) CheckSchema(db *sqlx.DB) {
	logger := log.GetLogger("SQL")
	logger.Info(str.PadCenter(" "+string(sm)+" ", 60, "="))

	sqb := db.Builder()
	for it := tables.Iterator(); it.Next(); {
		tb, val := sm.Table(it.Key()), it.Value()

		sqb.Reset()
		sql, args := sqb.Select().From(tb).Limit(1).Build()
		err := db.Get(val, sql, args...)
		if err == nil {
			logger.Infof("%s = OK", tb)
			continue
		}
		if errors.Is(err, sqlx.ErrNoRows) {
			logger.Warnf("%s = %s", tb, err)
			continue
		}
		logger.Errorf("%s = %s", tb, err)
	}
}

func (sm Schema) UpdateSchema(db *sqlx.DB) error {
	logger := log.GetLogger("SQL")
	logger.Info(str.PadCenter(" "+string(sm)+" ", 60, "="))

	return xsqls.ApplySchemaChanges(db, string(sm), data.FS, "sqls/"+app.DBType(), logger)
}

func (sm Schema) VacuumSchema(db *sqlx.DB) error {
	logger := log.GetLogger("SQL")
	logger.Info(str.PadCenter(" "+string(sm)+" ", 60, "="))

	for it := tables.Iterator(); it.Next(); {
		tb := sm.Table(it.Key())

		_, err := db.Exec("VACUUM FULL " + tb)
		if err != nil {
			return err
		}
	}
	return nil
}
