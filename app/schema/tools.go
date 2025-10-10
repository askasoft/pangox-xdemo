package schema

import (
	"errors"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox/xwa/xsqls"
)

func (sm Schema) ExecSQL(sqls string) error {
	return xsqls.ExecSQL(app.SDB(), string(sm), sqls)
}

func (sm Schema) CheckSchema(tx sqlx.Sqlx) {
	logger := log.GetLogger("SQL")
	logger.Info(str.Repeat("=", 40))

	sqb := tx.Builder()
	for it := tables.Iterator(); it.Next(); {
		tb, val := sm.Table(it.Key()), it.Value()

		sqb.Reset()
		sql, args := sqb.Select().From(tb).Limit(1).Build()
		err := tx.Get(val, sql, args...)
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

func (sm Schema) VacuumSchema(tx sqlx.Sqlx) error {
	logger := log.GetLogger("SQL")
	logger.Info(str.Repeat("=", 40))

	for it := tables.Iterator(); it.Next(); {
		tb := sm.Table(it.Key())

		_, err := tx.Exec("VACUUM FULL " + tb)
		if err != nil {
			return err
		}
	}
	return nil
}
