package schema

import (
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pangox-xdemo/app/args"
	"github.com/askasoft/pangox/xfs"
	"github.com/askasoft/pangox/xwa/xsqbs"
)

func (sm Schema) CountFiles(tx sqlx.Sqlx, fqa *args.FileQueryArg) (cnt int, err error) {
	sqb := tx.Builder()

	sqb.Count()
	sqb.From(sm.TableFiles())
	fqa.AddFilters(sqb)
	sql, args := sqb.Build()

	err = tx.Get(&cnt, sql, args...)
	return
}

func (sm Schema) FindFiles(tx sqlx.Sqlx, fqa *args.FileQueryArg, cols ...string) (files []*xfs.File, err error) {
	sqb := tx.Builder()

	sqb.Select(cols...)
	sqb.From(sm.TableFiles())
	fqa.AddFilters(sqb)
	fqa.AddOrders(sqb, "id")
	fqa.AddPager(sqb)
	sql, args := sqb.Build()

	err = tx.Select(&files, sql, args...)
	return
}

func (sm Schema) DeleteFilesQuery(tx sqlx.Sqlx, fqa *args.FileQueryArg) (int64, error) {
	sqb := tx.Builder()

	sqb.Delete(sm.TableFiles())
	fqa.AddFilters(sqb)
	sql, args := sqb.Build()

	return tx.Update(sql, args...)
}

func (sm Schema) DeleteFiles(tx sqlx.Sqlx, ids ...string) (int64, error) {
	sqb := tx.Builder()

	sqb.Delete(sm.TableFiles())
	xsqbs.AddIn(sqb, "id", ids)
	sql, args := sqb.Build()

	return tx.Update(sql, args...)
}

func (sm Schema) UpdateFiles(tx sqlx.Sqlx, fua *args.FileUpdatesArg) (int64, error) {
	sqb := tx.Builder()

	sqb.Update(sm.TableFiles())
	fua.AddUpdates(sqb)
	xsqbs.AddIn(sqb, "id", fua.PKs())
	sql, args := sqb.Build()

	return tx.Update(sql, args...)
}
