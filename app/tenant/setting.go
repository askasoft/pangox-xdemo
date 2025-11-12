package tenant

import (
	"strings"
	"sync"
	"time"

	"github.com/askasoft/pango/bol"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/models"
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
			var err error
			sr, err = buildSettingVarsReplacer(stg.Value)
			if err != nil {
				tt.Logger("SET").Errorf("Invalid tenant_vars: %s", stg.Value)
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

func buildSettingVarsReplacer(vars string) (*strings.Replacer, error) {
	i := ini.NewIni()

	err := i.LoadData(str.NewReader(vars))
	if err != nil {
		return nil, err
	}

	var kvs []string
	sec := i.Section("")
	for _, key := range sec.Keys() {
		kvs = append(kvs, "{{"+key+"}}", sec.GetString(key))
	}
	return str.NewReplacer(kvs...), nil
}

func (tt *Tenant) SettingVarsReplacer() (*strings.Replacer, error) {
	return buildSettingVarsReplacer(tt.SV("tenant_vars"))
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
