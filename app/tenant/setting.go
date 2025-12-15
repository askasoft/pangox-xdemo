package tenant

import (
	"strings"
	"sync"
	"time"

	"github.com/askasoft/pango/bol"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
	"github.com/askasoft/pangox-xdemo/app/utils/varutil"
)

// TSETS write lock
var muTSETS sync.Mutex

func (tt *Tenant) PurgeSettings() {
	muTSETS.Lock()
	app.TSETS.Remove(string(tt.Schema))
	muTSETS.Unlock()
}

func (tt *Tenant) getSettings() map[string]string {
	if stgs, ok := app.TSETS.Get(string(tt.Schema)); ok {
		return stgs
	}

	muTSETS.Lock()
	defer muTSETS.Unlock()

	// get again to prevent duplicated load
	if stgs, ok := app.TSETS.Get(string(tt.Schema)); ok {
		return stgs
	}

	stgs, err := tt.loadSettings(app.SDB())
	if err != nil {
		panic(err)
	}

	app.TSETS.Set(string(tt.Schema), stgs)
	return stgs
}

func (tt *Tenant) loadSettings(tx sqlx.Sqlx) (map[string]string, error) {
	sqb := tx.Builder()
	sqb.Select().From(tt.TableSettings())
	sql, args := sqb.Build()

	settings := []*models.Setting{}
	if err := tx.Select(&settings, sql, args...); err != nil {
		return nil, err
	}

	stgs := make(map[string]string, len(settings))

	var sr *str.Replacer
	for _, stg := range settings {
		if stg.Name == "tenant_vars" {
			vars, err := varutil.BuildVariables(stg.Value)
			if err != nil {
				tt.Logger("SET").Errorf("Invalid tenant_vars: %s", stg.Value)
			} else {
				sr = varutil.BuildVarReplacer(vars)
			}
			break
		}
	}

	for _, stg := range settings {
		cv := stg.Value
		if sr != nil && stg.Validation == "" && (stg.Style == models.SettingStyleDefault || stg.Style == models.SettingStyleTextarea) {
			cv = sr.Replace(cv)
		}
		stgs[stg.Name] = cv
	}

	return stgs, nil
}

func (tt *Tenant) VariablesReplacer() *strings.Replacer {
	return varutil.BuildVarReplacer(tt.variables)
}

func (tt *Tenant) Variables() map[string]string {
	return tt.variables
}

func (tt *Tenant) Settings() map[string]string {
	return tt.settings
}

func (tt *Tenant) SV(key string, defs ...string) string {
	val := tt.settings[key]
	if val == "" && len(defs) > 0 {
		return defs[0]
	}
	return val
}

func (tt *Tenant) SVs(key string) []string {
	return str.FieldsByte(tt.SV(key), '\t')
}

func (tt *Tenant) SB(key string, defs ...bool) bool {
	return bol.Atob(tt.SV(key), defs...)
}

func (tt *Tenant) SD(key string, defs ...time.Duration) time.Duration {
	return tmu.Atod(tt.SV(key), defs...)
}

func (tt *Tenant) SI(key string, defs ...int) int {
	return num.Atoi(tt.SV(key), defs...)
}

func (tt *Tenant) SL(key string, defs ...int64) int64 {
	return num.Atol(tt.SV(key), defs...)
}

func (tt *Tenant) SZ(key string, defs ...int64) int64 {
	return num.Atoz(key, defs...)
}
