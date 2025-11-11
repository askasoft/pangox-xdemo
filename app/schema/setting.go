package schema

import (
	"errors"
	"time"

	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pangox-xdemo/app/models"
)

func (sm Schema) SelectSettings(tx sqlx.Sqlx, items ...string) (settings []*models.Setting, err error) {
	sqb := tx.Builder()
	sqb.Select().From(sm.TableSettings())
	if len(items) > 0 {
		sqb.In("name", items)
	}
	sqb.Order("order")
	sqb.Order("name")
	sql, args := sqb.Build()

	err = tx.Select(&settings, sql, args...)
	return
}

func (sm Schema) UpdateSettingValue(tx sqlx.Sqlx, name, value string) (int64, error) {
	sqb := tx.Builder()

	sqb.Update(sm.TableSettings())
	sqb.Setc("value", value)
	sqb.Eq("name", name)
	sql, args := sqb.Build()

	r, err := tx.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func (sm Schema) ListSettingsByRole(tx sqlx.Sqlx, actor, role string) (settings []*models.Setting, err error) {
	sqb := tx.Builder()
	sqb.Select().From(sm.TableSettings())
	sqb.Gte(actor, role)
	sqb.Order("order")
	sqb.Order("name")
	sql, args := sqb.Build()

	err = tx.Select(&settings, sql, args...)
	return
}

type UnsavedSettingItemsError struct {
	Locale string
	Items  []string
}

func (usie *UnsavedSettingItemsError) Error() string {
	nms := make([]string, 0, len(usie.Items))
	for _, it := range usie.Items {
		nms = append(nms, tbs.GetText(usie.Locale, "setting."+it, it))
	}
	return tbs.Format(usie.Locale, "setting.error.unsaved", str.Join(nms, ", "))
}

func AsUnsavedSettingItemsError(err error) (usie *UnsavedSettingItemsError, ok bool) {
	ok = errors.As(err, &usie)
	return
}

func IsUnsavedSettingItemsError(err error) bool {
	_, ok := AsUnsavedSettingItemsError(err)
	return ok
}

func (sm Schema) SaveSettingsByRole(tx sqlx.Sqlx, au *models.User, settings []*models.Setting, locale string) error {
	sqb := tx.Builder()
	sqb.Update(sm.TableSettings())
	sqb.Setc("value", "")
	sqb.Setc("updated_at", "")
	sqb.Eq("name", "")
	sqb.Gte("editor", "")
	sql := tx.Rebind(sqb.SQL())

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var eits []string

	now := time.Now()
	for _, stg := range settings {
		r, err := stmt.Exec(stg.Value, now, stg.Name, au.Role)
		if err != nil {
			return err
		}

		cnt, _ := r.RowsAffected()
		if cnt != 1 {
			eits = append(eits, stg.Name)
		}
	}

	if len(eits) > 0 {
		return &UnsavedSettingItemsError{locale, eits}
	}

	return nil
}
